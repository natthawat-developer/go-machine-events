package events

// **MachineRefillEvent ใช้เมื่อมีการเติมสต็อก**
type MachineRefillEvent struct {
	Refill    int    `json:"refill"`
	MachineID string `json:"machine_id"`
}

// **StockLevelOkEvent ใช้เมื่อสต็อกกลับมาปกติ**
type StockLevelOkEvent struct {
	MachineID string `json:"machine_id"`
}

// Implement IEvent สำหรับ MachineRefillEvent
func (e MachineRefillEvent) Type() string      { return "refill" }
func (e MachineRefillEvent) GetMachineID() string { return e.MachineID } // ✅ เปลี่ยนชื่อ method

// Implement IEvent สำหรับ StockLevelOkEvent
func (e StockLevelOkEvent) Type() string      { return "stock_level_ok" }
func (e StockLevelOkEvent) GetMachineID() string { return e.MachineID } // ✅ เปลี่ยนชื่อ method
