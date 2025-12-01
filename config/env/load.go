package env

import (
	"errors"

	"github.com/joho/godotenv"
)

// Global environment variable accessor
var G *environment

func Load() error {
	if G != nil {
		return errors.New("environment already loaded")
	}

	e := environment{}

	godotenv.Load()

	if err := bindEnv(&e); err != nil {
		return err
	}

	G = &e

	return nil
}
