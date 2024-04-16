package manager

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ConnManager struct {
	Role          resp.NodeRole
	server        *net.Listener
	ClientManager *ClientManager
	NodeInfo      *resp.NodeInfo
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

		message := buf[:read_bytes]

		if len(message) == 0 {
			continue
		}

		handler, _ := resp.HandleResp(message)
		switch handler.(type) {
		case *resp.Array:
			arr := handler.(*resp.Array)
			cmd_handler := cmd.NewCMD(arr.Parsed, cmd.CMD_OPTS{
				Store:       n.ClientManager.store,
				ReplicaInfo: n.NodeInfo,
			})

			switch {
			case cmd_handler.Name == cmd.CMD_SET:
				cmd_handler.Process(&conn, nil)
			default:
				cmd_handler.Process(&conn, nil)
			}
		}
	}
}

func (n *ConnManager) Start() {
	fmt.Println("Server started")

	go n.ClientManager.setup() // setup the clients
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
