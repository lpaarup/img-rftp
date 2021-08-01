package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/lpaarup/img-rftp/pkg/common"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Retries int
	Timeout time.Duration
}

func (s *Server) ListenAndServe(port int) error {
	log.Infof("started server on addr: 0.0.0.0:%d", port)

	// Create a UDP connection
	conn, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not listen on port %d: %v", port, err)
	}
	defer conn.Close()

	return s.serve(conn)
}

func (s *Server) serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}

	if s.Retries == 0 {
		s.Retries = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}

	// Loop for new client connections, check if they are
	// sending a valid Read Request and handle it.
	var rrq common.ReadReq
	for {
		buf := make([]byte, common.DatagramSize)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Infof("[%s] bad request: %v", addr, err)
			continue
		}

		go s.handle(addr.String(), rrq)
	}
}

func (s Server) handle(addr string, rrq common.ReadReq) {
	log.Infof("[%s] requested file: %s", addr, rrq.Filename)

	// Use a new UDP connection to send the data, not to
	// block the principal connection
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Printf("[%s] could not dial: %v", addr, err)
	}
	defer conn.Close()

	// Download the desired image
	res, err := http.Get(rrq.Filename)
	if err != nil {
		err := s.sendError(conn, common.Err{Error: common.ErrNotFound, Message: err.Error()})
		if err != nil {
			log.Printf("[%s] failed to send error: %v", addr, err)
		}
		return
	}

	// Check if is exists
	if res.StatusCode != http.StatusOK {
		log.Errorf("received invalid status for %s, %s", rrq.Filename, res.Status)
		err := s.sendError(conn, common.Err{Error: common.ErrNotFound, Message: fmt.Sprintf("could not find %s", rrq.Filename)})
		if err != nil {
			log.Printf("[%s] failed to send error: %v", addr, err)
		}
		return
	}

	defer res.Body.Close()

	var (
		ackP  common.Ack
		errP  common.Err
		dataP = common.Data{Payload: res.Body}
		buf   = make([]byte, common.DatagramSize)
	)

NEXTPACKET:
	// While the sent data size if equal to the determined size,
	// keep sending data
	for n := common.DatagramSize; n == common.DatagramSize; {
		// Create a new data packet
		b, err := dataP.MarshalBinary()
		if err != nil {
			return
		}

	RETRY:
		for i := 0; i < s.Retries; i++ {
			n, err = conn.Write(b)
			if err != nil {
				return
			}

			// Wait for the specified timeout to see if there
			// is an ACK comming from the client.
			// If not, retry
			conn.SetReadDeadline(time.Now().Add(s.Timeout))
			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue RETRY
				}
				return
			}

			switch {
			case ackP.UnmarshalBinary(buf) == nil:
				if uint16(ackP) == dataP.Block {
					continue NEXTPACKET
				}
			case errP.UnmarshalBinary(buf) == nil:
				log.Printf("[%s] received error: %s", conn.RemoteAddr(), errP.Message)
				return
			}
		}
		log.Printf("[%s] exhausted retries", conn.RemoteAddr())
		return
	}

	log.Printf("[%s] send %d blocks", conn.RemoteAddr(), dataP.Block)
}

func (Server) sendError(conn net.Conn, e common.Err) error {
	b, err := e.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.Write(b)
	if err != nil {
		return err
	}

	return nil
}
