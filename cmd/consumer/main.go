package main

import (
	"go-machine-events/config"
	"go-machine-events/pkg/logger"

	"github.com/IBM/sarama"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, nil)
	if err != nil {
		log.Error("Error creating Kafka consumer group: %v", err)
		return
	}
	defer consumerGroup.Close()

	log.Info("Kafka Consumer started with group: %s", cfg.Kafka.GroupID)
}
