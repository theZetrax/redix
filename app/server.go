package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/manager"
)

func main() {
	port := "6379"

	cm := manager.NewClientManager()
	server := &manager.ConnManager{
		ClientManager: cm,
	}

	server.Serve(port)
	server.Start()
}
