package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/manager"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var port string
var replica_of string

func main() {
	flag.StringVar(&port, "port", "6379", "Port to listen on, default: 6379")
	flag.StringVar(&replica_of, "replicaof", "", "Replicate another Redis instance")

	flag.Parse()

	// check for replicaof flag
	args_raw := os.Args[1:]
	for i, arg := range args_raw {
	MATCH_ARG:
		switch arg {
		case "--replicaof":
			idx_host := i + 1
			idx_port := i + 2

			if idx_host >= len(args_raw) || idx_port >= len(args_raw) {
				log.Println(
					"Error: replicaof flag requires a host and port argument",
				)
				break
			}

			replica_of_val := args_raw[idx_host] + ":" + args_raw[idx_port]
			// set the replicaof flag
			flag.Set("replicaof", replica_of_val)
			break MATCH_ARG
		}
	}

	store := repository.NewStore()
	var replica_info *resp.NodeInfo = &resp.NodeInfo{
		Host:                "localhost",
		Port:                port,
		Role:                resp.RoleMaster,
		MasterReplicaId:     "1",
		MasterReplicaOffset: "0",
	}

	if replica_of != "" {
		replica_host, replica_port, found := strings.Cut(replica_of, ":")
		if !found {
			log.Printf("Error: Invalid master node address: %v\n", replica_of)
			os.Exit(1)
		}
		replica_info = manager.NewReplica(replica_host, replica_port)
	}

	cm := manager.NewClientManager(store, replica_info)
	server := &manager.ConnManager{
		ClientManager: cm,
		Role:          resp.RoleMaster,
		ReplicaInfo:   replica_info,
	}
	server.ConnectToMaster(replica_info)

	server.Serve(port)
	server.Start()
}
