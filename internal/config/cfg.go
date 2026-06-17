package config

import "flag"

type Cfg struct {
	ServerAddress string
	ResultHost    string
}

func Load() Cfg {
	serverAddress := flag.String("a", ":8080", "server address")
	resultAddress := flag.String("b", "http://localhost:8080", "result host")
	flag.Parse()

	return Cfg{
		ServerAddress: *serverAddress,
		ResultHost:    *resultAddress,
	}
}
