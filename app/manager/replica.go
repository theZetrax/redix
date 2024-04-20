package manager

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func NewReplica(host, port string) *resp.NodeInfo {
	return &resp.NodeInfo{
		Host: host,
		Port: port,
	}
}

// Perform a handshake with master node.
func Handshake(master_port string, node_port string) (conn net.Conn, err error) {
	fmt.Println("Connecting to master: ", "localhost:"+master_port)
	conn, err = net.Dial("tcp", "localhost:"+master_port)
	if err != nil {
		return nil, err
	}

	// messages to send to the master
	messages := [][]byte{
		resp.EncodeArray(resp.EncodeBulkString("PING")),
		// replication configuration
		resp.EncodeArray(
			resp.EncodeBulkString("REPLCONF"),
			resp.EncodeBulkString("listening-port"),
			resp.EncodeBulkString(node_port),
		),
		// partial sync PSYNC
		resp.EncodeArray(
			resp.EncodeBulkString("REPLCONF"),
			resp.EncodeBulkString("capa"),
			resp.EncodeBulkString("psync2"),
		),
		resp.EncodeArray(
			resp.EncodeBulkString("PSYNC"),
			resp.EncodeBulkString("?"),
			resp.EncodeBulkString("-1"),
		),
	}

	for _, message := range messages {
		// send ping to master
		buf := make([]byte, 1024)
		if _, err := conn.Write(message); err != nil {
			log.Println("Failed to write to master: ", err)
			return nil, err
		}

		// read response from master
		read_bytes, err := conn.Read(buf)
		if err != nil {
			log.Println("Failed to read from master: ", err)
			return nil, err
		}

		raw := string(buf[:read_bytes])
		log.Printf(
			"Recieved From Master[%s]: %q\n",
			"localhost:"+master_port,
			strings.ReplaceAll(raw, resp.CRLF, "\\r\\n"),
		)
	}

	log.Printf("Connected to master: %s\n", "localhost:"+master_port)
	return conn, nil
}
