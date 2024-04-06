package main

import (
	"fmt"
	"io"
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
	var connInstance net.Conn
	var err error

	if !config.IsMaster {
		connInstance, err = service.Handshake(config.ReplicaOf.Raw, config.Port)
		if err != nil {
			log.Printf("Error connecting to master[%s]: %s\n", config.ReplicaOf.Raw, err.Error())
			os.Exit(1)
		}
		resp_handler := &service.ResponseHandler{
			StorageEngine: storageEngine,
		}

		defer connInstance.Close()

		// handle incoming responses from master
		// and update the storage engine
		for {
			var read int
			buf := make([]byte, 1024)
			read, err = connInstance.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading from connection: ", err.Error())
				os.Exit(1)
			}

			resp_handler.HandleResponse(buf[:read])
		}

	} else {
		req_handler := &service.ReqHandler{
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
		defer req_handler.Close()

		// handle incoming connections
		for {
			// accept connection if master node
			// else use the replication connection
			connInstance, err = l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			var read int // read bytes length
			buf := make([]byte, 1024)
			read, err = connInstance.Read(buf)
			if err != nil {
				fmt.Println("Error reading from connection: ", err.Error())
				connInstance.Close()
				os.Exit(1)
			}

			shouldClose := service.IsLongLived(buf, read)

			go req_handler.HandleRequest(
				connInstance,
				&buf,
				read,
				service.RequestHandlerOptions{
					IsMaster:    config.IsMaster,
					ShouldClose: shouldClose,
				},
			)
		}
	} // end of else

}
