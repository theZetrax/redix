package manager

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func NewReplica(host, port string) *resp.NodeInfo {
	return &resp.NodeInfo{
		Host: host,
		Port: port,
	}
}

// handshake with master node.
func Handshake(port string) (conn net.Conn, err error) {
	fmt.Println("Connecting to master: ", "localhost:"+port)
	conn, err = net.Dial("tcp", "localhost:"+port)
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
			resp.EncodeBulkString(port),
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
		if _, err = io.Copy(conn, bytes.NewReader(message)); err != nil {
			log.Printf(
				"Error writing to master[%s]: %s\n",
				"localhost:"+port,
				err.Error(),
			)
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		read_bytes, err := conn.Read(buf) // read the response
		if err != nil {
			log.Panicln("Failed to read:", err)
		}

		raw := string(buf[:read_bytes])
		log.Printf(
			"Master[%s] raw response: %s\n",
			"localhost:"+port,
			strings.ReplaceAll(raw, resp.CRLF, "\\r\\n"),
		)
	}

	log.Printf("Connected to master: %s\n", "localhost:"+port)
	return conn, nil
}
