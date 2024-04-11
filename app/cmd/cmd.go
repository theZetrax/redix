package cmd

import (
	"log"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/repository"
)

type CMD_TYPE string
type CMD_OPTS struct {
	Store *repository.Store
}
type CMD_HANDLER func(opts CMD_OPTS, args []any) []byte
type CMD struct {
	Name    CMD_TYPE
	Args    []any
	handler CMD_HANDLER
	Store   *repository.Store
}

const (
	CMD_PING    CMD_TYPE = "PING"
	CMD_SET     CMD_TYPE = "SET"
	CMD_ECHO    CMD_TYPE = "ECHO"
	CMD_INVALID CMD_TYPE = "INVALID"
	CMD_GET     CMD_TYPE = "GET"
)

func NewCMD(raw []any, opts CMD_OPTS) *CMD {
	cmd := &CMD{
		Name:  CMD_INVALID,
		Args:  raw[1:],
		Store: opts.Store,
	}

	switch strings.ToUpper(raw[0].(string)) {
	case "PING":
		cmd.Name = CMD_PING
		cmd.handler = handlePing
	case "SET":
		cmd.Name = CMD_SET
		cmd.handler = handleSet
	case "ECHO":
		cmd.Name = CMD_ECHO
		cmd.handler = handleEcho
	case "GET":
		cmd.Name = CMD_GET
		cmd.handler = handleGet
	}

	return cmd
}

func (c *CMD) Process() []byte {
	log.Println("Executing Command: ", c.Name, c.Args)
	return c.handler(CMD_OPTS{Store: c.Store}, c.Args)
}
