package pubsub

import (
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

// **StartListening เริ่มฟัง event จาก Kafka**
func (p *PubSub) StartListening() {
	partitionConsumer, err := p.consumer.ConsumePartition(p.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming Kafka topic: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("Kafka Subscriber started...")

	for msg := range partitionConsumer.Messages() {
		topic := p.topic
		if handler, found := p.handlers[topic]; found {
			handler(msg.Value)
		} else {
			log.Printf("No handler found for topic: %s\n", topic)
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
