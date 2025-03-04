package services

import (
	"encoding/json"
	"log"

	"go-machine-events/internal/events"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
)

// **RefillService ใช้สำหรับจัดการ Refill Event**
type RefillService struct {
	machineRepo *machine.MachineRepository
	pubsub      *pubsub.PubSub
}

// **NewRefillService สร้าง RefillService ใหม่**
func NewRefillService(machineRepo *machine.MachineRepository, ps *pubsub.PubSub) *RefillService {
	return &RefillService{
		machineRepo: machineRepo,
		pubsub:      ps,
	}
}

// **HandleRefill จัดการเติมสต็อก**
func (s *RefillService) HandleRefill(event events.MachineRefillEvent) {
	// ✅ เรียกใช้ event.MachineID() แทนการเข้าถึง field ตรงๆ
	machine, exists := s.machineRepo.GetMachine(event.GetMachineID())
	if !exists {
		log.Printf("Machine %s not found", event.GetMachineID())
		return
	}

	// เติมสต็อก
	machine.StockLevel += event.Refill
	log.Printf("Machine %s refilled by %d, new stock: %d", event.GetMachineID(), event.Refill, machine.StockLevel)

	// ถ้าสต็อกกลับมา ≥ 3 ส่ง StockLevelOkEvent
	if machine.StockLevel >= 3 && machine.LowStockWarning {
		machine.LowStockWarning = false
		stockOkEvent := events.StockLevelOkEvent{MachineID: event.GetMachineID()}

		eventData, err := json.Marshal(stockOkEvent)
		if err != nil {
			log.Printf("Error marshaling StockLevelOkEvent: %v", err)
			return
		}

		s.pubsub.PublishEvent(eventData)
		log.Printf("StockLevelOkEvent sent for Machine %s", event.GetMachineID())
	}
}
