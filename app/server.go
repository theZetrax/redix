package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

const (
	PORT = "6379"
)

func main() {
	handler := &internal.HttpHandler{}

	l, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		fmt.Println("Failed to bind to port " + PORT)
		os.Exit(1)
	}
	fmt.Println("Redis-server listening on: ", PORT)
	defer l.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler.HandleConnection(conn)
	}()

	<-sigChan
}
