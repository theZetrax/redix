package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/conn"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

func main() {
	internal.InitFlags()
	flag.Parse()

	storageEngine := repository.NewStorageEngine()
	handler := &conn.HttpHandler{
		StorageEngine: storageEngine,
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+internal.PORT)
	if err != nil {
		fmt.Printf("Failed to bind to port: %s", internal.PORT)
		os.Exit(1)
	}
	fmt.Printf("Redis-server listening on: %s", internal.PORT)
	defer l.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			go handler.HandleConnection(conn)
		}
	}()

	<-sigChan
}
