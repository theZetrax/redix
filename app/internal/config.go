package internal

type Config struct {
	Port      string
	ReplicaOf HOST
	IsMaster  bool
}

func NewConfig(cli_args CLIArgs) *Config {
	config := &Config{}
	config.Port = cli_args.GetPort()
	config.ReplicaOf, config.IsMaster = cli_args.GetReplicaOf()

	return config
}
