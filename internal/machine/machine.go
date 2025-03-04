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
		machine.StockLevel += amount
		r.log.Info("Machine %s updated. New stock: %d", machine.ID, machine.StockLevel) 
	} else {
		r.log.Error("Machine %s not found", id)
	}
}


func (r *MachineRepository) GetMachine(id string) (*Machine, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machine, exists := r.machines[id]
	return machine, exists
}


func (r *MachineRepository) ListMachines() []*Machine {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machines := []*Machine{}
	for _, machine := range r.machines {
		machines = append(machines, machine)
	}
	return machines
}
