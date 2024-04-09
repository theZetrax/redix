package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/service"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

func main() {
	cli_args := internal.InitFlags()
	config := internal.NewConfig(cli_args)

	storageEngine := repository.NewStorageEngine()

	// connection instance
	var err error

	fmt.Println("Starting redis-server: ", config.IsMaster)

	if !config.IsMaster {
		replicaConn, err := service.Handshake(config.ReplicaOf.Raw, config.Port)
		if err != nil {
			log.Printf("Error connecting to master[%s]: %s\n", config.ReplicaOf.Raw, err.Error())
			os.Exit(1)
		}
		resp_handler := &service.ReplicaNode{
			StorageEngine: storageEngine,
		}

		// defer replicaConn.Close()

		// handle incoming responses from master
		// and update the storage engine
		go func() {
			for {
				buf := make([]byte, 1024)
				read, err := replicaConn.Read(buf)
				if err != nil {
					log.Println("Error reading from connection: ", err.Error())
					replicaConn.Close()
					return
				}

				log.Println("Data: ", string(buf[:read]))
				resp_handler.Handle(buf[:read])
			}
		}()
	}

	req_handler := &service.MainNode{
		StorageEngine: storageEngine,
		Config:        config,
		ConnPool:      make(map[string]net.Conn),
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+config.Port)
	if err != nil {
		fmt.Println("Failed to bind to port:", config.Port)
		os.Exit(1)
	}

	fmt.Printf("Redis-server listening on: %s\n", config.Port)

	defer l.Close()
	// defer req_handler.Close()

	// handle incoming connections
	for {
		connInstance, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		read, err := connInstance.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			continue
		}

		go req_handler.Handle(
			connInstance,
			buf,
			read,
			service.MainNodeOptions{
				IsMaster:    config.IsMaster,
				ShouldClose: !service.IsLongLived(buf, read),
			},
		)
	}
}
