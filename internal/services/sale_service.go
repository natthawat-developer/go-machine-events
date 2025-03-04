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
	machineID := event.GetMachineID()
	machine, exists := s.machineRepo.GetMachine(machineID)

	if !exists {
		s.log.Error("Machine %s not found", machineID)
		return
	}

	if machine.StockLevel < event.Sold {
		s.log.Error("Machine %s has insufficient stock", machineID)
		return
	}

	machine.StockLevel -= event.Sold
	s.log.Info("Machine %s sold %d, new stock: %d", machineID, event.Sold, machine.StockLevel)

	if machine.StockLevel < 3 && !machine.LowStockWarning {
		machine.LowStockWarning = true
		lowStockEvent := events.LowStockWarningEvent{MachineID: machineID}

		eventData, err := json.Marshal(lowStockEvent)
		if err != nil {
			s.log.Error("Error marshaling LowStockWarningEvent: %v", err)
			return
		}

		s.pubsub.PublishEvent(eventData)
		s.log.Info("LowStockWarningEvent sent for Machine %s", machineID)
	}
}
