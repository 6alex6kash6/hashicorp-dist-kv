package main

import (
	"distributed-kv/internal/config"
	"distributed-kv/internal/raft"
	"distributed-kv/internal/server"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	node := raft.NewRaftNode(config)

	if config.JoinAddr != "" {
		go func() error {
			url := url.URL{
				Scheme: "http",
				Host:   config.JoinAddr,
				Path:   "join",
			}

			req, err := http.NewRequest(http.MethodPost, url.String(), nil)
			if err != nil {
				return err
			}

			req.Header.Add("X-peer-addr", config.NodeAddr)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("non 200 status code: %d", resp.StatusCode)
			}

			return nil
		}()

	}

	server := server.NewServer(node, config.HTTPPort)

	server.Run()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
}
