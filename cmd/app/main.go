package main

import (
	"go-machine-events/config"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
	"go-machine-events/pkg/logger"
	"go-machine-events/pkg/utils"
	"time"
)

func main() {
	log := logger.NewLogger()

	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
		return
	}

	machineRepo := machine.NewMachineRepository()
	machineRepo.AddMachine("001", 10)
	machineRepo.AddMachine("002", 10)
	machineRepo.AddMachine("003", 10)

	pubsubService, err := pubsub.NewPubSub(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		log.Fatal("Failed to initialize PubSub: %v", err)
	}
	defer pubsubService.Close()

	go func() {
		saleSub := machine.NewSaleSubscriber(machineRepo)
		refillSub := machine.NewRefillSubscriber(machineRepo)

		pubsubService.Subscribe("sale", saleSub.HandleSaleEvent)
		pubsubService.Subscribe("refill", refillSub.HandleRefillEvent)
	}()

	go func() {
		for range 50 {
			event, err := utils.GenerateEvent()
			if err != nil {
				log.Error("Error generating event: %v", err)
				continue
			}
			pubsubService.PublishEvent(event)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {}
}
