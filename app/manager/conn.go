package manager

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/logger"
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
func (n *ConnManager) ConnectToMaster() {
	conn, err := Handshake(n.NodeInfo.MasterPort, n.NodeInfo.Port)
	if err != nil {
		fmt.Println("Failed to connect to master: ", err)
		os.Exit(1)
	}

	psync_cmd := resp.EncodeArray(
		resp.EncodeBulkString("PSYNC"),
		resp.EncodeBulkString("?"),
		resp.EncodeBulkString("-1"),
	)

	// send PSYNC command to master
	if _, err := conn.Write(psync_cmd); err != nil {
		log.Println("Failed to write to master: ", err)
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

		log.Println("OffsetCounter: ", n.NodeInfo.OffsetCount, "BufferCount: ", read_bytes)
		if n.NodeInfo.OffsetCount != -1 {
			n.NodeInfo.OffsetCount += read_bytes
		}
		logger.LogResp("From Master: ", buf[:read_bytes])
		message := buf[:read_bytes]

		if len(message) == 0 {
			continue
		}

		type requestData struct {
			segment      any
			segment_type resp.RESP_TYPE
		}
		req_data := make([]requestData, 0)
		for segment, t, rest, err := resp.Parse(message); ; segment, t, rest, err = resp.Parse(rest) {
			if err != nil {
				break
			}
			req_data = append(req_data, requestData{
				segment:      segment,
				segment_type: t,
			})
		}

		for _, d := range req_data {
			switch d.segment_type {
			case resp.TYPE_ARRAY:
				arr := d.segment
				cmd_handler := cmd.NewCMD(arr.([]any), cmd.CMD_OPTS{
					Store:       n.ClientManager.store,
					ReplicaInfo: n.ClientManager.node_info,
				})

				switch {
				case cmd_handler.Name == cmd.CMD_SET:
					cmd_handler.Process(nil, nil)
				case cmd_handler.Name == cmd.CMD_REPLCONF:
					cmd_handler.Process(&conn, nil)
				}
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
