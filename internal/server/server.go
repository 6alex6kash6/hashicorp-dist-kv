package server

import (
	"distributed-kv/internal/protocol"
	"distributed-kv/internal/raft"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	raft *raft.RaftNode
	port int
}

func NewServer(raft *raft.RaftNode, port int) *Server {
	return &Server{
		raft: raft,
		port: port,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /kv/{key}", s.handleGet)

	mux.HandleFunc("POST /kv/{key}", s.handlePut)

	mux.HandleFunc("POST /join", s.handleJoin)

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := protocol.Operation{
		Type:  protocol.PUT,
		Key:   key,
		Value: value,
	}

	encodedOp, err := op.Encode()
	if err != nil {
		fmt.Println("Err while encoding operation")
	}
	s.raft.ApplyOperation(encodedOp)

	fmt.Println("Put called")
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	val, _ := s.raft.Store.Get(key)
	w.Write(val)
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	peer := r.Header.Get("X-peer-addr")

	err := s.raft.AddNode(peer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Printf("Peer %s successfully joined cluster", peer)

	w.WriteHeader(http.StatusOK)
}
