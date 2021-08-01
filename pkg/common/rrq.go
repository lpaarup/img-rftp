package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type ReadReq struct {
	Filename, Mode string
}

func (rrq ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if rrq.Mode != "" {
		mode = rrq.Mode
	}

	capacity := 2 + 2 + len(rrq.Filename) + 1 + len(mode) + 1

	b := new(bytes.Buffer)
	b.Grow(capacity)

	err := binary.Write(b, binary.BigEndian, OpRRQ)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(rrq.Filename)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (rrq *ReadReq) UnmarshalBinary(b []byte) error {
	r := bytes.NewBuffer(b)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return fmt.Errorf("invalid code for OpRRQ, received %d", code)
	}

	rrq.Filename, err = r.ReadString(0)
	if err != nil {
		return err
	}

	rrq.Filename = strings.TrimRight(rrq.Filename, "\x00")
	if len(rrq.Filename) == 0 {
		return errors.New("filename in OpRRQ is empty")
	}

	rrq.Mode, err = r.ReadString(0)
	if err != nil {
		return err
	}

	rrq.Mode = strings.TrimRight(rrq.Mode, "\x00")
	if len(rrq.Mode) == 0 {
		return errors.New("mode in OpRRQ is empty")
	}

	if strings.ToLower(rrq.Mode) != "octet" {
		return errors.New("only binary supported")
	}

	return nil
}
