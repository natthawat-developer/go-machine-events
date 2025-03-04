package main

import (
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
	"go-machine-events/pkg/logger"
	"time"
)

func main() {
	log := logger.NewLogger() // ✅ ใช้ logger

	// สร้างเครื่องจักร
	machineRepo := machine.NewMachineRepository()
	machineRepo.AddMachine("001", 10)
	machineRepo.AddMachine("002", 10)
	machineRepo.AddMachine("003", 10)

	// ✅ ใช้ NewPubSub และเช็ค error
	pubsubService, err := pubsub.NewPubSub([]string{"localhost:9092"}, "machine-events")
	if err != nil {
		log.Fatal("Failed to initialize PubSub: %v", err) // ✅ ใช้ log.Fatal()
	}
	defer pubsubService.Close()

	// ✅ ใช้ machine.NewSaleSubscriber
	saleSub := machine.NewSaleSubscriber(machineRepo)
	refillSub := machine.NewRefillSubscriber(machineRepo)

	// ✅ สมัคร subscriber โดยใช้ handler function
	pubsubService.Subscribe("sale", saleSub.HandleSaleEvent)
	pubsubService.Subscribe("refill", refillSub.HandleRefillEvent)

	// ✅ ส่ง Events จำลอง
	pubsubService.PublishEvent([]byte(`{"type": "sale", "sold": 2, "machine_id": "001"}`))
	pubsubService.PublishEvent([]byte(`{"type": "refill", "refill": 5, "machine_id": "002"}`))

	// รอให้ Kafka ประมวลผล
	time.Sleep(5 * time.Second)
	log.Info("Done!") // ✅ ใช้ log.Info()
}
