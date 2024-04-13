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
func (n *ConnManager) ConnectToMaster(node_info *resp.NodeInfo) {
	if node_info == nil {
		return
	}

	conn, err := Handshake(node_info.MasterPort, node_info.Port)
	if err != nil {
		fmt.Println("Failed to connect to master: ", err)
		os.Exit(1)
	}

	for {
		buf := make([]byte, 1024)
		read_bytes, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by master")
				break
			}
		}

		fmt.Println("Received from master: ", string(buf[:read_bytes]))
	}
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
		go n.ClientManager.setup()
		go client.Setup()
		go client.Read()
	}
}
