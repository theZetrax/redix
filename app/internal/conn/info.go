// Responsible for handling the INFO command
// Author: Zablon Dawit
// Date: Mar-30-2024
package conn

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

func (h *HttpHandler) handleInfo(conn net.Conn, _ internal.Request) {
	resp_raw := []string{"#Replication", "role:master"}
	resp := strings.Join(resp_raw, "\n")

	_, err := conn.Write([]byte(encoder.NewBulkString(resp)))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
