package machine

import (
	"encoding/json"
	"go-machine-events/pkg/logger"
)

type EventSale struct {
	MachineID string `json:"machine_id"`
	Sold      int    `json:"sold"`
}

type SaleSubscriber struct {
	repo *MachineRepository
	log  *logger.Logger
}

func NewSaleSubscriber(repo *MachineRepository) *SaleSubscriber {
	return &SaleSubscriber{
		repo: repo,
		log:  logger.NewLogger(),
	}
}

func (s *SaleSubscriber) HandleSaleEvent(data []byte) {
    s.log.Info("Received Sale Event: %s", string(data))

    var event EventSale
    err := json.Unmarshal(data, &event)
    if err != nil {
        s.log.Error("Failed to parse Sale Event: %v", err)
        return
    }

    s.repo.UpdateStock(event.MachineID, -event.Sold)
    s.log.Info("Sale event processed for machine %s, sold %d", event.MachineID, event.Sold)
}
