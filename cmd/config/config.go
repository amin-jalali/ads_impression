package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
)

// Config structure to hold configuration values
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	App struct {
		TTL int `yaml:"ttl"`
	} `yaml:"app"`
}

// LoadConfig reads the configuration from file
func LoadConfig() (*Config, error) {
	var cfg Config

	execPath, err := os.Getwd()
	configPath := filepath.Join(execPath, "config.yml")

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
		return nil, err
	}
	return &cfg, nil
}
