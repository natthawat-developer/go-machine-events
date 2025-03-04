package machine

import (
	"fmt"
	"sync"
)

// **Machine โครงสร้างข้อมูลของเครื่องจักร**
type Machine struct {
	ID              string
	StockLevel      int
	LowStockWarning bool // ✅ เพิ่มฟิลด์นี้
}

// **MachineRepository เป็นตัวจัดการข้อมูลเครื่องจักร**
type MachineRepository struct {
	machines map[string]*Machine
	mu       sync.RWMutex // ใช้ Mutex ป้องกันการแก้ไขพร้อมกัน
}

// **NewMachineRepository สร้าง Repository ใหม่**
func NewMachineRepository() *MachineRepository {
	return &MachineRepository{
		machines: make(map[string]*Machine),
	}
}

// **AddMachine เพิ่มเครื่องจักรใหม่**
func (r *MachineRepository) AddMachine(id string, stock int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.machines[id] = &Machine{ID: id, StockLevel: stock}
	fmt.Printf("Machine %s added with stock %d\n", id, stock)
}

// **UpdateStock อัปเดตสต็อกของเครื่องจักร**
func (r *MachineRepository) UpdateStock(id string, amount int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if machine, exists := r.machines[id]; exists {
		machine.StockLevel += amount
		fmt.Printf("Machine %s updated. New stock: %d\n", machine.ID, machine.StockLevel)
	} else {
		fmt.Printf("Machine %s not found\n", id)
	}
}

// **GetMachine ดึงข้อมูลเครื่องจักรตาม ID**
func (r *MachineRepository) GetMachine(id string) (*Machine, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machine, exists := r.machines[id]
	return machine, exists
}

// **ListMachines คืนค่าเครื่องจักรทั้งหมด**
func (r *MachineRepository) ListMachines() []*Machine {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machines := []*Machine{}
	for _, machine := range r.machines {
		machines = append(machines, machine)
	}
	return machines
}
