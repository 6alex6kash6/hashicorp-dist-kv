package raft

import (
	"distributed-kv/internal/protocol"
	"distributed-kv/internal/storage"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/raft"
)

type KVFSM struct {
	store storage.KVStorage
}

func NewKVFSM(store storage.KVStorage) *KVFSM {
	return &KVFSM{
		store: store,
	}
}

func (f *KVFSM) Apply(log *raft.Log) interface{} {

	switch log.Type {
	case raft.LogCommand:
		var op protocol.Operation
		if err := op.Decode(log.Data); err != nil {
			fmt.Print(err)
			return nil
		}
		switch op.Type {
		case protocol.PUT:
			return f.store.Set(op.Key, op.Value)
		case protocol.DELETE:
			return f.store.Delete(op.Key)
		}
	}
	return nil
}

func (f *KVFSM) Snapshot() (raft.FSMSnapshot, error) {
	snapshot := f.store.Snapshot()
	// Create a snapshot of the key-value store
	// This is used for compact log entries and efficient transfers
	return &kvSnapshot{snapshot}, nil
}

func (f *KVFSM) Restore(sink io.ReadCloser) error {
	snapshot, err := newKVSnapshot(sink)
	if err != nil {
		log.Printf("[ERR] Failed to create snapshot: %v", err)
		return err
	}
	defer snapshot.Release()

	if err := f.store.Restore(snapshot.Data); err != nil {
		log.Printf("[ERR] Failed to restore snapshot: %v", err)
		return err
	}

	return nil
}
