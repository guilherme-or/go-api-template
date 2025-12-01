package config

import (
	"github.com/guilherme-or/go-api-template/config/env"
	"github.com/guilherme-or/go-api-template/config/logger"
)

func MustSetup() {
	// environment config load
	if err := env.Load(); err != nil {
		panic(err)
	}

	// logger config load
	logger.Load()
}
