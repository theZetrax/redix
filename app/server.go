package main

import (
	"fmt"
	"net"
	"os"

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
	defer l.Close()
	fmt.Println("Redis-server listening on: ", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler.HandleConnection(conn)
	}
}
