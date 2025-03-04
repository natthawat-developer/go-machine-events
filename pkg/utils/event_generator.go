package utils

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Event โครงสร้างของ Event
type Event struct {
	Type      string `json:"type"`
	MachineID string `json:"machine_id"`
	Quantity  int    `json:"quantity"`
}

// randomMachine เลือกเครื่องจักรแบบสุ่ม
func randomMachine() string {
	machines := []string{"001", "002", "003"}
	return machines[rand.Intn(len(machines))]
}

// GenerateEvent สร้าง Event แบบสุ่ม
func GenerateEvent() ([]byte, error) {
	rand.Seed(time.Now().UnixNano())

	var event Event
	if rand.Float64() < 0.5 {
		event = Event{
			Type:      "sale",
			MachineID: randomMachine(),
			Quantity:  []int{1, 2}[rand.Intn(2)], // ขาย 1 หรือ 2 ชิ้น
		}
	} else {
		event = Event{
			Type:      "refill",
			MachineID: randomMachine(),
			Quantity:  []int{3, 5}[rand.Intn(2)], // เติม 3 หรือ 5 ชิ้น
		}
	}

	return json.Marshal(event)
}
