package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	os.Args = []string{"app-name", "-a", "localhost:1111", "-b", "localhost:2222"}
	cfg := Load()

	assert.Equal(t, cfg.ServerAddress, "localhost:1111")
	assert.Equal(t, cfg.ResultHost, "localhost:2222")
}
