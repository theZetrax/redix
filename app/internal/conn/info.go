// Responsible for handling the INFO command
// Author: Zablon Dawit
// Date: Mar-30-2024
package conn

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
)

func (h *HttpHandler) handleInfo(conn net.Conn, _ internal.Request) {
	var role string
	if h.Config.IsMaster {
		role = "master"
	} else {
		role = "slave"
	}

	resp_raw := []string{
		"#Replication",
		fmt.Sprintf("role:%s", role),
	}
	resp := strings.Join(resp_raw, "\n")

	_, err := conn.Write([]byte(encoder.NewBulkString(resp)))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
