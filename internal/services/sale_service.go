package services

import (
	"encoding/json"
	"log"

	"go-machine-events/internal/events"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
)

// **SaleService ใช้สำหรับจัดการ Sale Event**
type SaleService struct {
	machineRepo *machine.MachineRepository
	pubsub      *pubsub.PubSub
}

// **NewSaleService สร้าง SaleService ใหม่**
func NewSaleService(machineRepo *machine.MachineRepository, ps *pubsub.PubSub) *SaleService {
	return &SaleService{
		machineRepo: machineRepo,
		pubsub:      ps,
	}
}

// **HandleSale จัดการการขายสินค้า**
func (s *SaleService) HandleSale(event events.MachineSaleEvent) {
	machineID := event.GetMachineID() // ✅ เรียกใช้เมธอด MachineID() 
	machine, exists := s.machineRepo.GetMachine(machineID) // ✅ รับค่าทั้ง 2 ตัวแปร

	if !exists { // ✅ ตรวจสอบว่าพบเครื่องจักรหรือไม่
		log.Printf("Machine %s not found", machineID)
		return
	}

	// ✅ ตรวจสอบว่าสต็อกพอหรือไม่
	if machine.StockLevel < event.Sold {
		log.Printf("Machine %s has insufficient stock", machineID)
		return
	}

	// ✅ อัปเดตสต็อก
	machine.StockLevel -= event.Sold
	log.Printf("Machine %s sold %d, new stock: %d", machineID, event.Sold, machine.StockLevel)

	// ✅ ถ้าสต็อกต่ำกว่า 3 และยังไม่มี LowStockWarning => ส่ง LowStockWarningEvent
	if machine.StockLevel < 3 && !machine.LowStockWarning {
		machine.LowStockWarning = true
		lowStockEvent := events.LowStockWarningEvent{MachineID: machineID}

		eventData, err := json.Marshal(lowStockEvent)
		if err != nil {
			log.Printf("Error marshaling LowStockWarningEvent: %v", err)
			return
		}

		s.pubsub.PublishEvent(eventData)
		log.Printf("LowStockWarningEvent sent for Machine %s", machineID)
	}
}
