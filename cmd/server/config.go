package server

import "github.com/ilyakaznacheev/cleanenv"

type ConfigEnv struct {
	Value string `env:"KEY"`
}

func GetEnv(key string) (string, error) {
	var cfg ConfigEnv

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return "", err
	}

	return cfg.Value, nil
}
