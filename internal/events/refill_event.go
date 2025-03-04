package events

type MachineRefillEvent struct {
	Refill    int    `json:"refill"`
	MachineID string `json:"machine_id"`
}

type StockLevelOkEvent struct {
	MachineID string `json:"machine_id"`
}

func (e MachineRefillEvent) Type() string         { return "refill" }
func (e MachineRefillEvent) GetMachineID() string { return e.MachineID }

func (e StockLevelOkEvent) Type() string         { return "stock_level_ok" }
func (e StockLevelOkEvent) GetMachineID() string { return e.MachineID }
