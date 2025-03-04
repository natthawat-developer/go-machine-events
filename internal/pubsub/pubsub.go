package pubsub

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// **PubSub เป็นตัวกลางสำหรับ Publisher และ Subscriber**
type PubSub struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	topic    string
	handlers map[string]func([]byte) // ✅ เพิ่ม map เก็บ event handlers ตาม topic
}

// **NewPubSub สร้าง Publisher และ Subscriber**
func NewPubSub(brokers []string, topic string) (*PubSub, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &PubSub{
		producer: producer,
		consumer: consumer,
		topic:    topic,
		handlers: make(map[string]func([]byte)), // ✅ กำหนด map ให้ใช้งาน
	}, nil
}

// **Subscribe เพิ่ม handler ให้ topic ที่ต้องการ**
func (p *PubSub) Subscribe(topic string, handler func([]byte)) {
	p.handlers[topic] = handler
	log.Printf("Subscribed to topic: %s\n", topic)
}


func (p *PubSub) StartListening() {
	partitionConsumer, err := p.consumer.ConsumePartition(p.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming Kafka topic: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("Kafka Subscriber started...")

	for msg := range partitionConsumer.Messages() {
		// ✅ แปลง JSON หา type ก่อน
		var event map[string]interface{}
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error parsing event: %v\n", err)
			continue
		}

		// ✅ หา "type" เพื่อเรียก Handler ที่ถูกต้อง
		eventType, ok := event["type"].(string)
		if !ok {
			log.Println("Invalid event type")
			continue
		}

		if handler, found := p.handlers[eventType]; found {
			handler(msg.Value)
		} else {
			log.Printf("No handler found for event type: %s\n", eventType)
		}
	}
}


// **PublishEvent ส่ง Event ไป Kafka**
func (p *PubSub) PublishEvent(event []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(event),
	}

	_, _, err := p.producer.SendMessage(message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	log.Printf("Published Event: %s\n", event)
	return nil
}

// **Close ปิด Producer และ Consumer**
func (p *PubSub) Close() {
	p.producer.Close()
	p.consumer.Close()
}
