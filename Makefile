build:
	go build -o bin/dist-kv

run: build
	./bin/dist-kv

bootstrap: build
	./bin/dist-kv -node-id=node1 -node-addr=localhost:8001 -cluster-id=my-cluster -bootstrap=true -http-port=8080

run_node_2:
	./bin/dist-kv -node-id=node2 -node-addr=localhost:8002 -cluster-id=my-cluster -http-port=8091 -join-addr=localhost:8080

run_node_3:
	./bin/dist-kv -node-id=node3 -node-addr=localhost:8003 -cluster-id=my-cluster -http-port=8093 -join-addr=localhost:8080
