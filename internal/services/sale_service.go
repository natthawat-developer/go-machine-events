package services

import (
	"encoding/json"

	"go-machine-events/internal/events"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
	"go-machine-events/pkg/logger"
)

type SaleService struct {
	machineRepo *machine.MachineRepository
	pubsub      *pubsub.PubSub
	log         *logger.Logger
}

func NewSaleService(machineRepo *machine.MachineRepository, ps *pubsub.PubSub) *SaleService {
	return &SaleService{
		machineRepo: machineRepo,
		pubsub:      ps,
		log:         logger.NewLogger(),
	}
}

func (s *SaleService) HandleSale(event events.MachineSaleEvent) {
    machine, exists := s.machineRepo.GetMachine(event.MachineID) // ใช้ event.MachineID แทน event.GetMachineID() ถ้าใช้โครงสร้างนี้
    if !exists {
        s.log.Error("Machine %s not found", event.MachineID)
        return
    }

    if machine.StockLevel < event.Quantity {
        s.log.Info("Not enough stock for machine %s to sell %d items", event.MachineID, event.Quantity)
        return
    }

    machine.StockLevel -= event.Quantity
    s.log.Info("Machine %s sold %d items, new stock: %d", event.MachineID, event.Quantity, machine.StockLevel)

    // Handle low stock warning
    if machine.StockLevel < 3 && !machine.LowStockWarning {
        machine.LowStockWarning = true
        lowStockEvent := events.StockLevelLowEvent{MachineID: event.MachineID}

        eventData, err := json.Marshal(lowStockEvent)
        if err != nil {
            s.log.Error("Error marshaling StockLevelLowEvent: %v", err)
            return
        }

        s.pubsub.PublishEvent(eventData)
        s.log.Info("StockLevelLowEvent sent for Machine %s", event.MachineID)
    }
}
