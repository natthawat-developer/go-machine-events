package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// KafkaClient ใช้จัดการ Kafka Producer และ Consumer
type KafkaClient struct {
	brokers      []string
	topic        string
	producer     sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
}

// NewKafkaClient สร้าง KafkaClient ใหม่
func NewKafkaClient(brokers []string, topic string, groupID string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	// ตั้งค่า Producer
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// ตั้งค่า Consumer
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		brokers:      brokers,
		topic:        topic,
		producer:     producer,
		consumerGroup: consumerGroup,
	}, nil
}

// PublishEvent ส่ง Event ไปยัง Kafka
func (kc *KafkaClient) PublishEvent(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: kc.topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := kc.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error publishing message to Kafka: %v", err)
		return err
	}

	log.Println("Message published successfully")
	return nil
}

// Subscribe รับ Event จาก Kafka และเรียก callback
func (kc *KafkaClient) Subscribe(callback func([]byte)) error {
	handler := &ConsumerHandler{callback: callback}
	ctx := context.Background()

	go func() {
		for {
			err := kc.consumerGroup.Consume(ctx, []string{kc.topic}, handler)
			if err != nil {
				log.Printf("Error consuming messages: %v", err)
			}
		}
	}()

	log.Println("Kafka Consumer started...")
	return nil
}
