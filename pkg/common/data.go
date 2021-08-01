package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(DatagramSize)

	d.Block++

	err := binary.Write(buf, binary.BigEndian, OpData)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, d.Block)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(buf, d.Payload, DatagramSize-4)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Data) UnmarshalBinary(b []byte) error {
	if l := len(b); l < 4 || l > DatagramSize {
		return errors.New("invalid data length")
	}
	var code OpCode

	err := binary.Read(bytes.NewReader(b[:2]), binary.BigEndian, &code)
	if err != nil {
		return errors.New("unable to read code")
	}

	if code != OpData {
		return fmt.Errorf("invalid code for OpData, received %d", code)
	}

	err = binary.Read(bytes.NewReader(b[2:4]), binary.BigEndian, &d.Block)
	if err != nil {
		return errors.New("unable to block code")
	}

	d.Payload = bytes.NewBuffer(b[4:])

	return nil
}
