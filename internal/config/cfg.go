package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"LOG_LEVEL"`
}

func MustLoad() Cfg {
	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "result base url")
	logLevel := flag.String("l", "debug", "log level")
	flag.Parse()

	cfg := Cfg{
		ServerAddress: *serverAddress,
		BaseURL:       *baseURL,
		LogLevel:      *logLevel,
	}

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
