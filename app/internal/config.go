package internal

type Config struct {
	Port       string
	ReplicaOf  HOST
	IsMaster   bool
	ReplId     string
	ReplOffset string
}

func NewConfig(cli_args CLIArgs) *Config {
	config := &Config{}
	config.Port = cli_args.GetPort()
	config.ReplicaOf, config.IsMaster = cli_args.GetReplicaOf()
	config.ReplId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	config.ReplOffset = "0"

	return config
}
