package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	capacity := 2 + 2
	b := new(bytes.Buffer)
	b.Grow(capacity)

	err := binary.Write(b, binary.BigEndian, OpAck)
	if err != nil {
		return nil, errors.New("could not write OpAck")
	}

	err = binary.Write(b, binary.BigEndian, a)
	if err != nil {
		return nil, errors.New("could not write block")
	}

	return b.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(b []byte) error {

	var code OpCode
	r := bytes.NewBuffer(b)
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return errors.New("unable to read code")
	}

	if code != OpAck {
		return fmt.Errorf("invalid code for OpAck, received %d", code)
	}

	err = binary.Read(r, binary.BigEndian, a)
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read block id")
	}

	return nil
}
