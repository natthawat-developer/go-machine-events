package kafka

import (
	"go-machine-events/pkg/logger"

	"github.com/IBM/sarama"
)

// Producer ใช้สำหรับจัดการ Kafka Producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	log      *logger.Logger // ✅ ใช้ logger
}

// NewProducer สร้าง Producer ใหม่
func NewProducer(brokers []string, topic string) (*Producer, error) {
	log := logger.NewLogger() // ✅ ใช้ logger

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_8_0_0

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error("Error creating Kafka producer: %v", err) // ✅ ใช้ log.Error()
		return nil, err
	}

	log.Info("Kafka producer initialized successfully")

	return &Producer{
		producer: producer,
		topic:    topic,
		log:      log, // ✅ ใช้ logger
	}, nil
}

// Publish ส่ง Event ไปยัง Kafka
func (p *Producer) Publish(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Error("Error publishing message to Kafka: %v", err) // ✅ ใช้ log.Error()
		return err
	}

	p.log.Info("Message published successfully to Kafka")
	return nil
}
