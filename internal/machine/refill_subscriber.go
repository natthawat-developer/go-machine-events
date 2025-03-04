package machine

import (
	"encoding/json"
	"go-machine-events/pkg/logger"
)

type EventRefill struct {
	MachineID string `json:"machine_id"`
	Refill    int    `json:"refill"`
}

type RefillSubscriber struct {
	repo *MachineRepository
	log  *logger.Logger
}

func NewRefillSubscriber(repo *MachineRepository) *RefillSubscriber {
	return &RefillSubscriber{
		repo: repo,
		log:  logger.NewLogger(),
	}
}

func (r *RefillSubscriber) HandleRefillEvent(data []byte) {
    r.log.Info("Received Refill Event: %s", string(data))

    var event EventRefill
    err := json.Unmarshal(data, &event)
    if err != nil {
        r.log.Error("Failed to parse Refill Event: %v", err)
        return
    }

    r.repo.UpdateStock(event.MachineID, event.Refill)
    r.log.Info("Refill event processed for machine %s, refill %d", event.MachineID, event.Refill)
}
