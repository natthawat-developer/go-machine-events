package services

import (
	"encoding/json"
	"go-machine-events/internal/events"
	"go-machine-events/internal/machine"
	"go-machine-events/internal/pubsub"
	"go-machine-events/pkg/logger"
)

type RefillService struct {
	machineRepo *machine.MachineRepository
	pubsub      *pubsub.PubSub
	log         *logger.Logger
}

func NewRefillService(machineRepo *machine.MachineRepository, ps *pubsub.PubSub) *RefillService {
	return &RefillService{
		machineRepo: machineRepo,
		pubsub:      ps,
		log:         logger.NewLogger(),
	}
}

func (s *RefillService) HandleRefill(event events.MachineRefillEvent) {
    if event.Refill <= 0 {
        s.log.Error("Invalid refill amount %d for machine %s", event.Refill, event.GetMachineID())
        return
    }

    machine, exists := s.machineRepo.GetMachine(event.GetMachineID())
    if !exists {
        s.log.Error("Machine %s not found", event.GetMachineID())
        return
    }

    machine.StockLevel += event.Refill
    s.log.Info("Machine %s refilled by %d, new stock: %d", event.GetMachineID(), event.Refill, machine.StockLevel)

    // Handle low stock warning
    if machine.StockLevel >= 3 && machine.LowStockWarning {
        machine.LowStockWarning = false
        stockOkEvent := events.StockLevelOkEvent{MachineID: event.GetMachineID()}

        eventData, err := json.Marshal(stockOkEvent)
        if err != nil {
            s.log.Error("Error marshaling StockLevelOkEvent: %v", err)
            return
        }

        s.pubsub.PublishEvent(eventData)
        s.log.Info("StockLevelOkEvent sent for Machine %s", event.GetMachineID())
    }
}
