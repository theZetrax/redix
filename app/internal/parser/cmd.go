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
)

const (
	CMD_ECHO = "echo"
	CMD_PING = "ping"
	CMD_SET  = "set"
	CMD_GET  = "get"
)

// do not modify this
var COMMANDS = []string{
	CMD_ECHO,
	CMD_PING,
	CMD_GET,
	CMD_SET,
}

type CMD struct {
	CMD  string
	Args []any
}

func NewCMD(raw_cmd []any) (CMD, error) {
	cmd := reflect.ValueOf(raw_cmd).Index(0).Interface().(string)

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
