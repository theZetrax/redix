// DESCRIPTION: This file contains the implementation of the replication protocol.
// AUTHOR: Zablon Dawit
// Date: APR-01-2024
package conn

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

// handshake with master node.
func Handshake(master_host string) {
	conn, err := net.Dial("tcp", master_host)
	if err != nil {
		log.Printf("Error connecting to master[%s]: %s\n", master_host, err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	message := encoder.NewArray(encoder.NewBulkString("PING"))
	if _, err = io.Copy(conn, strings.NewReader(message)); err != nil {
		log.Printf("Error writing to master[%s]: %s\n", master_host, err.Error())
		os.Exit(1)
	}

	log.Printf("Connected to master: [%s]\n", master_host)
}
