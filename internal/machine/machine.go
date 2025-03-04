package machine

import (
	"fmt"
)

// **SaleSubscriber รับ event เมื่อมีการขายสินค้า**
type SaleSubscriber struct {
	repo *MachineRepository
}

// **NewSaleSubscriber สร้าง Subscriber สำหรับ Sale Events**
func NewSaleSubscriber(repo *MachineRepository) *SaleSubscriber {
	return &SaleSubscriber{repo: repo}
}

// **HandleSaleEvent อัปเดตสต็อกเมื่อมีการขายสินค้า**
func (s *SaleSubscriber) HandleSaleEvent(data []byte) {
	// จำลองการอ่านข้อมูลจาก JSON
	machineID := "001" // ควรใช้ JSON Unmarshal จริง ๆ
	sold := 2

	s.repo.UpdateStock(machineID, -sold)
	fmt.Printf("Sale event processed for machine %s, sold %d\n", machineID, sold)
}

// **RefillSubscriber รับ event เมื่อมีการเติมสินค้า**
type RefillSubscriber struct {
	repo *MachineRepository
}

// **NewRefillSubscriber สร้าง Subscriber สำหรับ Refill Events**
func NewRefillSubscriber(repo *MachineRepository) *RefillSubscriber {
	return &RefillSubscriber{repo: repo}
}

// **HandleRefillEvent อัปเดตสต็อกเมื่อมีการเติมสินค้า**
func (r *RefillSubscriber) HandleRefillEvent(data []byte) {
	// จำลองการอ่านข้อมูลจาก JSON
	machineID := "002"
	refill := 5

	r.repo.UpdateStock(machineID, refill)
	fmt.Printf("Refill event processed for machine %s, refill %d\n", machineID, refill)
}
