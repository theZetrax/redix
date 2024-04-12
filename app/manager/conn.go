package manager

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ConnManager struct {
	Role          resp.NodeRole
	server        *net.Listener
	ClientManager *ClientManager
	ReplicaInfo   *resp.NodeInfo
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

// ConnectToMaster connects the replica to the master node.
func (n *ConnManager) ConnectToMaster(replicaInfo *resp.NodeInfo) {
	if replicaInfo == nil {
		return
	}

	Handshake(replicaInfo.Port)
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
