package client

import (
	"bytes"
	"io"
	"net"

	"github.com/pkg/errors"
)

type RFTPClient struct {
	remoteAddr string
}

func (c *RFTPClient) Read(filename string) (io.Reader, error) {
	conn, err := net.Dial("tcp", c.remoteAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not dial")
	}
	defer conn.Close()

	_, err = conn.Write([]byte("START"))
	if err != nil {
		return nil, errors.Wrap(err, "could not send start")
	}

	buf := make([]byte, 10000)
	_, err = conn.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "could not read")
	}

	return bytes.NewReader(buf), nil
}

func New(remote string) *RFTPClient {
	return &RFTPClient{remoteAddr: remote}
}
