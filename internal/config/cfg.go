package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func MustLoad() Cfg {
	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "result base url")
	flag.Parse()

	cfg := Cfg{
		ServerAddress: *serverAddress,
		BaseURL:       *baseURL,
	}

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
