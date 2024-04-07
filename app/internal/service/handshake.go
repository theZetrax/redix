// DESCRIPTION: This file contains the implementation of the replication protocol.
// AUTHOR: Zablon Dawit
// Date: APR-01-2024
package service

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

// handshake with master node.
func Handshake(master_host string, port string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", master_host)
	if err != nil {
		return nil, err
	}

	// messages to send to the master
	messages := []string{
		encoder.NewArray(encoder.NewBulkString(decoder.CMD_PING)), // ping
		// replication configuration
		encoder.NewArray(
			encoder.NewBulkString(decoder.CMD_REPLCONF),
			encoder.NewBulkString("listening-port"),
			encoder.NewBulkString(port),
		),
		// partial sync PSYNC
		encoder.NewArray(
			encoder.NewBulkString("REPLCONF"),
			encoder.NewBulkString("capa"),
			encoder.NewBulkString("psync2"),
		),
		encoder.NewArray(
			encoder.NewBulkString(decoder.CMD_PSYNC),
			encoder.NewBulkString("?"),
			encoder.NewBulkString("-1"),
		),
	}

	for _, message := range messages {
		// write the message to the master
		if _, err = io.Copy(conn, strings.NewReader(message)); err != nil {
			log.Printf(
				"Error writing to master[%s]: %s\n",
				master_host,
				err.Error(),
			)
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		read_bytes, err := conn.Read(buf) // read the response
		if err != nil {
			return nil, err
		}

		raw := string(buf[:read_bytes])
		log.Printf(
			"Master[%s] raw response: %s\n",
			master_host,
			strings.ReplaceAll(raw, decoder.CRLF, "\\r\\n"),
		)
	}

	log.Printf("Connected to master: %s\n", master_host)
	return conn, nil
}