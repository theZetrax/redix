package cmd

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/logger"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CMD_TYPE string
type CMD_OPTS struct {
	Store              *repository.Store
	ReplicaInfo        *resp.NodeInfo
	ConnectedNodeCount int
}
type CMD_HANDLER func(opts CMD_OPTS, args []any) []byte
type CMD_MULTI_HANDLER func(opts CMD_OPTS, args []any) [][]byte
type CMD struct {
	Name           CMD_TYPE
	Args           []any
	handler        CMD_HANDLER       // handle request with single response
	handleMultiple CMD_MULTI_HANDLER // handle request with multiple responses
	CMD_OPTS       CMD_OPTS
}

const (
	CMD_PING     CMD_TYPE = "PING"
	CMD_SET      CMD_TYPE = "SET"
	CMD_ECHO     CMD_TYPE = "ECHO"
	CMD_INVALID  CMD_TYPE = "INVALID"
	CMD_GET      CMD_TYPE = "GET"
	CMD_REPLCONF CMD_TYPE = "REPLCONF"
	CMD_PSYNC    CMD_TYPE = "PSYNC"
	CMD_INFO     CMD_TYPE = "INFO"
	CMD_WAIT     CMD_TYPE = "WAIT"
)

func NewCMD(raw []any, opts CMD_OPTS) *CMD {
	cmd := &CMD{
		Name:     CMD_INVALID,
		Args:     raw[1:],
		CMD_OPTS: opts,
	}

	switch strings.ToUpper(raw[0].(string)) {
	case string(CMD_PING):
		cmd.Name = CMD_PING
		cmd.handler = handlePing
	case string(CMD_SET):
		cmd.Name = CMD_SET
		cmd.handler = handleSet
	case string(CMD_GET):
		cmd.Name = CMD_GET
		cmd.handler = handleGet
	case string(CMD_ECHO):
		cmd.Name = CMD_ECHO
		cmd.handler = handleEcho
	case string(CMD_REPLCONF):
		cmd.Name = CMD_REPLCONF
		cmd.handler = handleReplConf
	case string(CMD_INFO):
		cmd.Name = CMD_INFO
		cmd.handler = handleInfo
	case string(CMD_WAIT):
		cmd.Name = CMD_WAIT
		cmd.handler = handleWait
	case string(CMD_PSYNC):
		cmd.Name = CMD_PSYNC
		cmd.handleMultiple = handlePsync
	}

	return cmd
}

// Process the command, and write the response to the client
// if successful, execute the post function
//
// 1. Execute the command handler
// 2. if exists, Execute the handler function
// 4. if exists, Execute the command handler for multiple responses
// 3. execute the post function
// 2. write the response to the client
func (c *CMD) Process(conn *net.Conn, post func()) {
	log.Println("Executing Command: ", c.Name, c.Args)
	if c.handler != nil {
		response := c.handler(c.CMD_OPTS, c.Args)

		if conn != nil {
			logger.LogResp("Sending Response: ", response)
			_, err := (*conn).Write(response)
			if err != nil {
				log.Println("Failed to write to master: ", err)
				os.Exit(1)
			}
		}
	}

	if c.handleMultiple != nil {
		responses := c.handleMultiple(c.CMD_OPTS, c.Args)

		if conn != nil {
			for _, response := range responses {
				_, err := (*conn).Write(response)
				if err != nil {
					log.Println("Failed to write to master: ", err)
					os.Exit(1)
				}
			}
		}
	}

	if post != nil {
		post()
	}
}
