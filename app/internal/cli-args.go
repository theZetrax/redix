package internal

import (
	"flag"
	"log"
	"os"
	"strings"
)

type CLIArgs struct {
	port       string
	replica_of string
}

type HOST struct {
	port string
	host string
	raw  string
}

func InitFlags() (cli_args CLIArgs) {
	//#region define the flags
	// port flag
	flag.StringVar(&cli_args.port, "port", "6379", "Port to listen on, default: 6379")
	flag.StringVar(&cli_args.replica_of, "p", "6379", "Port to listen on, default: 6379")

	// replicaof
	flag.StringVar(&cli_args.replica_of, "replicaof", "", "Replicate another Redis instance")
	//#endregion

	flag.Parse() // parse the flags

	// check for replicaof flag
	args_raw := os.Args[1:]
	for i, arg := range args_raw {
		switch arg {
		case "--replicaof":
			idx_host := i + 1
			idx_port := i + 2

			if idx_host >= len(args_raw) || idx_port >= len(args_raw) {
				log.Println(
					"Error: replicaof flag requires a host and port argument",
				)
				break
			}

			replica_of_val := args_raw[idx_host] + ":" + args_raw[idx_port]
			// set the replicaof flag
			flag.Set("replicaof", replica_of_val)

			break
		}
	}

	return cli_args
}

func (ca *CLIArgs) GetPort() string {
	return ca.port
}

func (ca *CLIArgs) GetReplicaOf() (_ HOST, is_master bool) {
	if ca.replica_of == "" {
		return HOST{}, true
	}

	host, port, found := strings.Cut(ca.replica_of, ":")
	if !found {
		return HOST{}, true
	}

	// if host is localhost
	if host == "localhost" {
		host = "0.0.0.0"
	}

	return HOST{
		port: port,
		host: host,
		raw:  ca.replica_of,
	}, false
}
