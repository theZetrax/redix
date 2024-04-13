package cmd

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/logger"
	"github.com/codecrafters-io/redis-starter-go/app/repository"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CMD_TYPE string
type CMD_OPTS struct {
	Store       *repository.Store
	ReplicaInfo *resp.NodeInfo
}
type CMD_HANDLER func(opts CMD_OPTS, args []any) []byte
type CMD_MULTI_HANDLER func(opts CMD_OPTS, args []any) [][]byte
type CMD struct {
	Name           CMD_TYPE
	Args           []any
	handler        CMD_HANDLER       // handle request with single response
	handleMultiple CMD_MULTI_HANDLER // handle request with multiple responses
	Store          *repository.Store
	ReplicaInfo    *resp.NodeInfo
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
)

func NewCMD(raw []any, opts CMD_OPTS) *CMD {
	cmd := &CMD{
		Name:        CMD_INVALID,
		Args:        raw[1:],
		Store:       opts.Store,
		ReplicaInfo: opts.ReplicaInfo,
	}

	switch strings.ToUpper(raw[0].(string)) {
	case string(CMD_PING):
		cmd.Name = CMD_PING
		cmd.handler = handlePing
	case string(CMD_SET):
		cmd.Name = CMD_SET
		cmd.handler = handleSet
	case string(CMD_ECHO):
		cmd.Name = CMD_ECHO
		cmd.handler = handleEcho
	case string(CMD_GET):
		cmd.Name = CMD_GET
		cmd.handler = handleGet
	case string(CMD_REPLCONF):
		cmd.Name = CMD_REPLCONF
		cmd.handler = handleReplConf
	case string(CMD_INFO):
		cmd.Name = CMD_INFO
		cmd.handler = handleInfo
	case string(CMD_PSYNC):
		cmd.Name = CMD_PSYNC
		cmd.handleMultiple = handlePsync
	}

	return cmd
}

func (c *CMD) Process(conn *net.Conn) {
	log.Println("Executing Command: ", c.Name, c.Args)
	if c.handler != nil {
		response := c.handler(CMD_OPTS{Store: c.Store, ReplicaInfo: c.ReplicaInfo}, c.Args)
		logger.LogResp("Sending Response: ", response)
		_, err := (*conn).Write(response)
		if err != nil {
			log.Println("Failed to write to master: ", err)
			os.Exit(1)
		}
	}

	if c.handleMultiple != nil {
		responses := c.handleMultiple(CMD_OPTS{Store: c.Store, ReplicaInfo: c.ReplicaInfo}, c.Args)

		for _, response := range responses {
			_, err := (*conn).Write(response)
			if err != nil {
				log.Println("Failed to write to master: ", err)
				os.Exit(1)
			}

			time.Sleep(4 * time.Millisecond) // sleep for 4ms
		}
	}
}
