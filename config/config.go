package config

import (
	"os"

	"go-machine-events/pkg/logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`
}


func LoadConfig() *Config {
	log := logger.NewLogger()

	config := &Config{}
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Error("Error reading config file: %v", err)
		return nil
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Error("Error parsing config file: %v", err)
		return nil
	}

	log.Info("Config loaded successfully: %+v", config)
	return config
}
