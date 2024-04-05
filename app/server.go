package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/service"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

func main() {
	cli_args := internal.InitFlags()
	config := internal.NewConfig(cli_args)

	storageEngine := repository.NewStorageEngine()
	handler := &service.HttpHandler{
		StorageEngine: storageEngine,
		Config:        config,
	}

	// connection instance
	var connInstance net.Conn
	var err error

	if !config.IsMaster {
		// replication connection
		connInstance, err = service.Handshake(config.ReplicaOf.Raw, config.Port)
		if err != nil {
			log.Printf("Error connecting to master[%s]: %s\n", config.ReplicaOf.Raw, err.Error())
			os.Exit(1)
		}
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+config.Port)
	if err != nil {
		fmt.Println("Failed to bind to port:", config.Port)
		os.Exit(1)
	}

	fmt.Printf("Redis-server listening on: %s\n", config.Port)

	defer l.Close()
	defer handler.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// handle incoming connections
	go func() {
		for {
			// accept connection if master node
			// else use the replication connection
			connInstance, err = l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			go handler.HandleConnection(connInstance)
		}
	}()

	<-sigChan
}
