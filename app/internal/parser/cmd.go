// Description: This file contains the parser for the command line arguments.
// The NewCMD function should return a CMD object with the command and its arguments.
// The command should be the first element of the raw_cmd slice and the arguments should be the rest of the elements.
// AUTHOR: Zablon Dawit
// DATE: Mar-30-2024
package parser

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// commands
const (
	CMD_ECHO = "ECHO"
	CMD_PING = "PING"
	CMD_SET  = "SET"
	CMD_GET  = "GET"
	CMD_INFO = "INFO"
)

// subcommands
const (
	SUB_PX = "PX"
)

// do not modify this
var COMMANDS = []string{
	CMD_ECHO,
	CMD_PING,
	CMD_GET,
	CMD_SET,
	CMD_INFO,
}

var SUB_COMMANDS = map[string][]string{
	CMD_SET: {SUB_PX},
}

type CMD struct {
	CMD  string
	Args []any
}

func NewCMD(raw_cmd []any) (CMD, error) {
	cmd := reflect.ValueOf(raw_cmd).Index(0).Interface().(string)
	cmd = strings.ToUpper(cmd) // convert to lowercase

	// check if the command is valid
	if !slices.Contains(COMMANDS, cmd) {
		return CMD{}, fmt.Errorf("Invalid command: %s", cmd)
	}

	args := raw_cmd[1:]
	return CMD{
		CMD:  cmd,
		Args: args,
	}, nil
}
