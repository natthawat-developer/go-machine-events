package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Producer ใช้สำหรับจัดการ Kafka Producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer สร้าง Producer ใหม่
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_8_0_0

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
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
		log.Printf("Error publishing message to Kafka: %v", err)
		return err
	}

	log.Println("Message published successfully")
	return nil
}
