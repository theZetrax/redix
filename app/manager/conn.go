package manager

import (
	"fmt"
	"net"
	"os"
)

type ConnManager struct {
	server        *net.Listener
	ClientManager *ClientManager
}

func (n *ConnManager) Serve(port string) {
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}

	n.server = &l

	fmt.Println("Server listening on localhost:" + port)
}

func (n *ConnManager) Start() {
	fmt.Println("Server started")
	for {
		conn, err := (*n.server).Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			(*n.server).Close()
			os.Exit(1)
		}

		fmt.Println("Accepted connection from: ", conn.RemoteAddr().String())

		client := NewClient(n.ClientManager, conn)
		go client.Setup()
		go client.Read()
	}
}
