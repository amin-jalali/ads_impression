package server

import "github.com/ilyakaznacheev/cleanenv"

// ConfigEnv holds environment variables
type ConfigEnv struct {
	Value string `env:"KEY"`
}

// GetEnv retrieves an environment variable using Cleanenv
func GetEnv(key string) (string, error) {
	var cfg ConfigEnv

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return "", err
	}

	return cfg.Value, nil
}
