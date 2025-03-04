package pubsub

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// **Publisher สำหรับ Kafka**
type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

// **NewPublisher สร้าง Publisher ใหม่**
func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		producer: producer,
		topic:    topic,
	}, nil
}

// **PublishEvent ส่ง Event ไป Kafka**
func (p *Publisher) PublishEvent(event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(data),
	}

	_, _, err = p.producer.SendMessage(message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	log.Printf("Published Event: %s\n", data)
	return nil
}

// **Close ปิด Kafka Producer**
func (p *Publisher) Close() error {
	return p.producer.Close()
}
