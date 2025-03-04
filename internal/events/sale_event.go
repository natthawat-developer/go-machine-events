package events

type MachineSaleEvent struct {
	Sold      int    `json:"sold"`
	MachineID string `json:"machine_id"`
	Quantity  int    `json:"quantity"`
}

type LowStockWarningEvent struct {
	MachineID string `json:"machine_id"`
}

type StockLevelLowEvent struct {
	MachineID string `json:"machine_id"`
}

func (e MachineSaleEvent) Type() string         { return "sale" }
func (e MachineSaleEvent) GetMachineID() string { return e.MachineID }

func (e LowStockWarningEvent) Type() string         { return "low_stock_warning" }
func (e LowStockWarningEvent) GetMachineID() string { return e.MachineID }
