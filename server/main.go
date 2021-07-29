package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "port to listen to")
)

func main() {
	flag.Parse()
	s := RFTPServer{}
	if err := s.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port)); err != nil {
		log.Fatalf("error on ListenAndServe: %v", err)
	}
}

type RFTPServer struct {
	Payload []byte
}

func (s *RFTPServer) ListenAndServe(addr string) error {
	log.Infof("started server on addr: %s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not listen on port 8080: %v", err)
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()

		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

func (s *RFTPServer) handle(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1000)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	switch s := string(buf[:n]); s {
	case "START":
		log.Infof("Received: %s", string(buf))

		b, err := ioutil.ReadFile("logo.png")
		if err != nil {
			return errors.Wrap(err, "could not open file")
		}

		_, err = conn.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}
