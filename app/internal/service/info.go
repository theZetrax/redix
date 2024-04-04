// Responsible for handling the INFO command
// Author: Zablon Dawit
// Date: Mar-30-2024
package service

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
		fmt.Sprintf("master_replid:%s", h.Config.ReplId),
		fmt.Sprintf("master_repl_offset:%s", h.Config.ReplOffset),
	}
	resp := strings.Join(resp_raw, "\n")

	_, err := conn.Write([]byte(encoder.NewBulkString(resp)))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
