package common

const (
	DatagramSize = 516
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	_
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnknown = iota
	ErrNotFound
)
