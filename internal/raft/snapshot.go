package raft

import (
	"io"

	"github.com/hashicorp/raft"
)

type kvSnapshot struct {
	Data []byte
}

func newKVSnapshot(r io.ReadCloser) (*kvSnapshot, error) {
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &kvSnapshot{Data: data}, nil
}

func (s *kvSnapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write(s.Data)
	if err != nil {
		sink.Cancel()
		return err
	}

	return sink.Close()
}

func (s *kvSnapshot) Release() {
	// No-op, we don't need to release any resources
}
