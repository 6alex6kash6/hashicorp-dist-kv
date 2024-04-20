package raft

import (
	"distributed-kv/internal/config"
	"distributed-kv/internal/storage"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

const defaultAddr = "localhost:8080"

type RaftNode struct {
	raft       *raft.Raft
	Store      storage.KVStorage
	nodeConfig *config.Config
}

func NewRaftNode(nodeConfig *config.Config) *RaftNode {
	store := storage.NewStorage()
	fsm := NewKVFSM(store)

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeConfig.NodeID)
	config.ElectionTimeout = 10 * time.Second

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()
	snap := raft.NewInmemSnapshotStore()

	addr, err := net.ResolveTCPAddr("tcp", nodeConfig.NodeAddr)
	if err != nil {
		fmt.Printf("Unable to resolve TCP addr %d", err)
	}

	transport, err := raft.NewTCPTransport(nodeConfig.NodeAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		fmt.Printf("Unable to create TCP transport %d", err)
	}

	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snap, transport)
	if err != nil {
		fmt.Printf("Unable to create Raft %d", err)
	}

	if nodeConfig.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		raftNode.BootstrapCluster(configuration)
	}

	return &RaftNode{raft: raftNode, Store: store, nodeConfig: nodeConfig}
}

func (rn *RaftNode) ApplyOperation(cmd []byte) {
	rn.raft.Apply(cmd, 5*time.Second)
}

func (rn *RaftNode) AddNode(peer string) error {
	future := rn.raft.AddVoter(raft.ServerID(peer), raft.ServerAddress(peer), 0, 10000)

	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

func (rn *RaftNode) join() {

}
