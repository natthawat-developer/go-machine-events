package events

// **MachineSaleEvent ใช้เมื่อมีการขายสินค้า**
type MachineSaleEvent struct {
	Sold      int    `json:"sold"`
	MachineID string `json:"machine_id"`
}

// **LowStockWarningEvent ใช้เมื่อสต็อกต่ำ**
type LowStockWarningEvent struct {
	MachineID string `json:"machine_id"`
}

// Implement IEvent สำหรับ MachineSaleEvent
func (e MachineSaleEvent) Type() string       { return "sale" }
func (e MachineSaleEvent) GetMachineID() string { return e.MachineID } // ✅ เปลี่ยนชื่อ method

// Implement IEvent สำหรับ LowStockWarningEvent
func (e LowStockWarningEvent) Type() string       { return "low_stock_warning" }
func (e LowStockWarningEvent) GetMachineID() string { return e.MachineID } // ✅ เปลี่ยนชื่อ method
