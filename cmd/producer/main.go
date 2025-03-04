package main

import (
	"go-machine-events/config"
	"go-machine-events/pkg/logger"

	"github.com/IBM/sarama"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger()

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, nil)
	if err != nil {
		log.Error("Error creating Kafka producer: %v", err)
		return
	}
	defer producer.Close()

	log.Info("Kafka Producer started")
}
