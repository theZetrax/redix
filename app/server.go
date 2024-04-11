package main

import (
	"flag"
	"log"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/manager"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
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

	cm := manager.NewClientManager(store)
	server := &manager.ConnManager{
		ClientManager: cm,
	}

	server.Serve(port)
	server.Start()
}
