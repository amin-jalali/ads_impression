package server

import (
	"errors"
	"os"
)

func GetEnv(key string) (error, string) {
	if value := os.Getenv(key); value != "" {
		return nil, value
	}
	return errors.New("env is not set"), ""
}
