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

	// ✅ Subscribe ก่อนเริ่มฟัง Kafka
	saleSub := machine.NewSaleSubscriber(machineRepo)
	refillSub := machine.NewRefillSubscriber(machineRepo)

	pubsubService.Subscribe("sale", saleSub.HandleSaleEvent)
	pubsubService.Subscribe("refill", refillSub.HandleRefillEvent)

	// ✅ เริ่มฟัง Kafka Event
	go pubsubService.StartListening()

	// ✅ ส่ง Event จำลองไป Kafka
	go func() {
		for i := 0; i < 10; i++ {
			event, err := utils.GenerateEvent()
			if err != nil {
				log.Error("Error generating event: %v", err)
				continue
			}
			pubsubService.PublishEvent(event)
			time.Sleep(500 * time.Millisecond) // ✅ ให้เวลาระหว่าง Event ชัดขึ้น
		}
	}()

	// ✅ ป้องกันโปรแกรมจบ
	select {} 
}
