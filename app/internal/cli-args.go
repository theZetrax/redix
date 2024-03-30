package internal

import "flag"

var PORT string

func InitFlags() {
	// port flag
	flag.StringVar(&PORT, "port", "6379", "Port to listen on, default: 6379")
	flag.StringVar(&PORT, "p", "6379", "Port to listen on, default: 6379")
}
