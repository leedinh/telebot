package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"  env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer
}

type HTTPServer struct {
	Address     string        `yaml:"host" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	// Load config from environment variables
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %s", configPath)
	}

	var config Config
	// Load the config from the file
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return &config
}
