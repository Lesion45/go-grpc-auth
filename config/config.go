package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token-ttl" env-default:"7h"`
	GRPC     GRPC
	Storage  Storage
}

type GRPC struct {
	Port int `yaml:"port"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"dbname"`
}

// MustLoad loads configuration from config.yaml
// Shut down the application if the config doesn't exist or if there is an error reading the config.
func MustLoad() *Config {
	configPath := "C:\\Users\\maus1\\GolandProjects\\go-grpc-auth\\config\\config.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
