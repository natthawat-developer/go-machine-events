package pubsub

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// **Subscriber สำหรับ Kafka**
type Subscriber struct {
	consumer sarama.Consumer
	topic    string
}

// **NewSubscriber สร้าง Subscriber ใหม่**
func NewSubscriber(brokers []string, topic string) (*Subscriber, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		consumer: consumer,
		topic:    topic,
	}, nil
}

// **Subscribe เริ่มรับ Event จาก Kafka**
func (s *Subscriber) Subscribe(handleFunc func(event map[string]interface{})) error {
	partitionConsumer, err := s.consumer.ConsumePartition(s.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	log.Println("Kafka Subscriber started...")

	for msg := range partitionConsumer.Messages() {
		var event map[string]interface{}
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		handleFunc(event) // ส่ง Event ให้ Handler
	}

	return nil
}

// **Close ปิด Kafka Consumer**
func (s *Subscriber) Close() error {
	return s.consumer.Close()
}
