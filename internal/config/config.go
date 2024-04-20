package config

import (
	"flag"
	"fmt"
)

type Config struct {
	NodeID    string
	NodeAddr  string
	ClusterID string
	Bootstrap bool
	JoinAddr  string
	HTTPPort  int
}

func ParseConfig() (*Config, error) {
	nodeID := flag.String("node-id", "", "Unique node ID for this node")
	nodeAddr := flag.String("node-addr", "", "Address on which this node listens for Raft communication")
	clusterID := flag.String("cluster-id", "", "Unique ID for the Raft cluster")
	bootstrap := flag.Bool("bootstrap", false, "Whether to bootstrap the Raft cluster")
	joinAddr := flag.String("join-addr", "", "Addres of leader node to join")
	httpPort := flag.Int("http-port", 0, "Port on which the HTTP server should listen")
	flag.Parse()

	if *nodeID == "" {
		return nil, fmt.Errorf("node-id is required")
	}

	if *nodeAddr == "" {
		return nil, fmt.Errorf("node-addr is required")
	}

	if *clusterID == "" {
		return nil, fmt.Errorf("cluster-id is required")
	}

	config := &Config{
		NodeID:    *nodeID,
		NodeAddr:  *nodeAddr,
		ClusterID: *clusterID,
		Bootstrap: *bootstrap,
		JoinAddr:  *joinAddr,
		HTTPPort:  *httpPort,
	}

	return config, nil
}
