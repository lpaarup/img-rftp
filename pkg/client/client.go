package client

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/lpaarup/img-rftp/pkg/common"
	"github.com/pkg/errors"
)

type Client struct {
	serverAddr string
}

func (c *Client) Read(filename string) (io.Reader, error) {
	// Create a UDP connection
	conn, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		return nil, errors.Wrap(err, "could not start connection")
	}
	defer conn.Close()

	// Create a Read Request for the given network image
	rrq := common.ReadReq{Filename: filename}
	b, err := rrq.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Resolve Server's address in order to use it in conn.WriteTo
	addr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve address %s: %v", c.serverAddr, err)
	}

	// Write the read request to the server
	_, err = conn.WriteTo(b, addr)
	if err != nil {
		return nil, errors.Wrap(err, "could not send start")
	}

	var (
		dataP common.Data
		ackP  common.Ack
		errP  common.Err
		buf   = make([]byte, common.DatagramSize)
		image = new(bytes.Buffer)
	)

	// As long as the received data is equal to the determined size, there are
	// still data packets to receive
	for n := common.DatagramSize; n == common.DatagramSize; {
		// Read data packets (or errors) from the connnection
		var addr net.Addr
		n, addr, err = conn.ReadFrom(buf)
		if err != nil {
			return nil, errors.Wrap(err, "could not read")
		}

		switch {
		// Check if the packet is Data
		case dataP.UnmarshalBinary(buf) == nil:
			// Copy the received data payload into the image buffer
			_, err = io.Copy(image, dataP.Payload)
			if err != nil {
				return nil, errors.Wrap(err, "could not copy payload data")
			}

			// Acknowledge the data block
			ackP = common.Ack(dataP.Block)
			b, err = ackP.MarshalBinary()
			if err != nil {
				return nil, errors.Wrap(err, "could not marshal ack")
			}

			_, err = conn.WriteTo(b, addr)
			if err != nil {
				return nil, errors.Wrap(err, "could not write data")
			}
			continue

		// Check if there was an error
		case errP.UnmarshalBinary(buf) == nil:
			return nil, fmt.Errorf(errP.Message)
		}
	}

	return image, nil
}

func New(remote string) *Client {
	return &Client{serverAddr: remote}
}
