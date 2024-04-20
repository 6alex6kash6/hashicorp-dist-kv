package protocol

import (
	"bytes"
	"encoding/gob"
)

type OperationType int

const (
	PUT OperationType = iota
	GET
	DELETE
)

type Operation struct {
	Type  OperationType
	Key   string
	Value []byte
}

func (op *Operation) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(op.Type); err != nil {
		return nil, err
	}

	if err := enc.Encode(op.Key); err != nil {
		return nil, err
	}

	if op.Type == PUT {
		if err := enc.Encode(op.Value); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (op *Operation) Decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&op.Type); err != nil {
		return err
	}

	if err := dec.Decode(&op.Key); err != nil {
		return err
	}

	if op.Type == PUT {
		if err := dec.Decode(&op.Value); err != nil {
			return err
		}
	}

	return nil
}
