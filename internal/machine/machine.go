package machine

import (
	"go-machine-events/pkg/logger"
	"sync"
)

type Machine struct {
	ID              string
	StockLevel      int
	LowStockWarning bool
}

type MachineRepository struct {
	machines map[string]*Machine
	mu       sync.RWMutex
	log      *logger.Logger
}

// **NewMachineRepository สร้าง Repository ใหม่**
func NewMachineRepository() *MachineRepository {
	return &MachineRepository{
		machines: make(map[string]*Machine),
		log:      logger.NewLogger(),
	}
}

// **AddMachine เพิ่มเครื่องจักรใหม่**
func (r *MachineRepository) AddMachine(id string, stock int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.machines[id] = &Machine{ID: id, StockLevel: stock}
	r.log.Info("Machine %s added with stock %d", id, stock)
}

// **UpdateStock อัปเดตสต็อกของเครื่องจักร**
func (r *MachineRepository) UpdateStock(id string, amount int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if machine, exists := r.machines[id]; exists {
		oldStock := machine.StockLevel
		machine.StockLevel += amount // ✅ อัปเดตสต็อก

		change := machine.StockLevel - oldStock
		r.log.Info("Machine %s updated. New stock: %d (change: %d)", machine.ID, machine.StockLevel, change)

		// ตรวจสอบว่าการขายทำให้สต็อกติดลบหรือไม่
		if machine.StockLevel < 0 {
			machine.StockLevel = 0
			r.log.Warn("Machine %s stock went negative! Reset to 0.", machine.ID)
		}

		// แจ้งเตือนถ้า Stock ต่ำ
		if machine.StockLevel < 3 {
			r.log.Warn("Machine %s is running low on stock! Current stock: %d", machine.ID, machine.StockLevel)
		}
	} else {
		r.log.Error("Machine %s not found", id)
	}
}



// **GetMachine ดึงข้อมูลเครื่องจักรตาม ID**
func (r *MachineRepository) GetMachine(id string) (*Machine, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machine, exists := r.machines[id]
	return machine, exists
}

// **ListMachines คืนค่ารายชื่อเครื่องจักรทั้งหมด**
func (r *MachineRepository) ListMachines() []*Machine {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machines := []*Machine{}
	for _, machine := range r.machines {
		machines = append(machines, machine)
	}
	return machines
}
