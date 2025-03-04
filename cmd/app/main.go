package main

import (
	"encoding/json"
	"go-machine-events/config"
	"go-machine-events/internal/events"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
	"go-machine-events/internal/services" // ✅ Import services
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

	saleService := services.NewSaleService(machineRepo, pubsubService)
	refillService := services.NewRefillService(machineRepo, pubsubService)

	pubsubService.Subscribe("sale", func(data []byte) {
		var event events.MachineSaleEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Error("Failed to unmarshal sale event: %v", err)
			return
		}
		saleService.HandleSale(event)
	})

	pubsubService.Subscribe("refill", func(data []byte) {
		var event events.MachineRefillEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Error("Failed to unmarshal refill event: %v", err)
			return
		}
		refillService.HandleRefill(event)
	})

	go pubsubService.StartListening()

	go func() {
		for range 5 {
			event, err := utils.GenerateEvent()
			if err != nil {
				log.Error("Error generating event: %v", err)
				continue
			}
			pubsubService.PublishEvent(event)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {}
}
