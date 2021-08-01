package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type Err struct {
	Error   ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {
	capacity := 2 + 2 + len(e.Message) + 1

	buf := new(bytes.Buffer)
	buf.Grow(capacity)

	err := binary.Write(buf, binary.BigEndian, OpErr)
	if err != nil {
		return nil, errors.New("unable to write OpErr")
	}

	err = binary.Write(buf, binary.BigEndian, e.Error)
	if err != nil {
		return nil, errors.New("unable to write ErrCode")
	}

	_, err = buf.WriteString(e.Message)
	if err != nil {
		return nil, errors.New("unable to write message")
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, errors.New("unable to write ending 0 byte")
	}

	return buf.Bytes(), nil
}

func (e *Err) UnmarshalBinary(b []byte) error {
	r := bytes.NewBuffer(b)

	var code OpCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return errors.New("unable to read code")
	}

	if code != OpErr {
		return fmt.Errorf("invalid code for OpErr, received %d", code)
	}

	err = binary.Read(r, binary.BigEndian, &e.Error)
	if err != nil {
		return errors.New("unable to read ErrCode")
	}

	e.Message, err = r.ReadString(0)
	if err != nil {
		return errors.New("unable to read message")
	}

	e.Message = strings.TrimRight(e.Message, "\x00")
	if len(e.Message) == 0 {
		return errors.New("unable to trim message")
	}

	return nil
}
