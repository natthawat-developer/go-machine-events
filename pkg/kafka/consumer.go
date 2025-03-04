package kafka

import (
	"github.com/IBM/sarama"
)

// ConsumerHandler ใช้สำหรับดึงข้อความจาก Kafka และส่งไปให้ callback
type ConsumerHandler struct {
	callback func([]byte)
}

// Setup ถูกเรียกเมื่อ Consumer เริ่มทำงาน
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup ถูกเรียกเมื่อ Consumer หยุดทำงาน
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim รับข้อความจาก Kafka และเรียก callback
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.callback(message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}
