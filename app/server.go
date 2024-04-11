package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/manager"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

func main() {
	port := "6379"

	store := repository.NewStore()

	cm := manager.NewClientManager(store)
	server := &manager.ConnManager{
		ClientManager: cm,
	}

	server.Serve(port)
	server.Start()
}
