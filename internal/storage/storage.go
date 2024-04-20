package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
)

type KVStorage interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte) bool
	Delete(key string) bool
	Snapshot() []byte
	Restore(data []byte) error
}

type Storage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string][]byte),
	}
}

func (s *Storage) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	return value, ok
}

func (s *Storage) Set(key string, value []byte) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return true
}

func (s *Storage) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)

	return true
}

func (s *Storage) Snapshot() []byte {
	encoded, err := encode(s.data)
	if err != nil {
		fmt.Println("err while making snapshot", err)
	}
	return encoded
}

func (s *Storage) Restore(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	snapshot, err := decode(data)
	if err != nil {
		return err
	}
	s.data = snapshot
	return nil
}

func encode(data map[string][]byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decode(data []byte) (map[string][]byte, error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)

	var snapshot map[string][]byte
	err := dec.Decode(&snapshot)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}
