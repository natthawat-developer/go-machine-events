package kafka

import (
	"context"
	"go-machine-events/pkg/logger"

	"github.com/IBM/sarama"
)

// KafkaClient ใช้จัดการ Kafka Producer และ Consumer
type KafkaClient struct {
	brokers       []string
	topic         string
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	log           *logger.Logger 
}

// NewKafkaClient สร้าง KafkaClient ใหม่
func NewKafkaClient(brokers []string, topic string, groupID string) (*KafkaClient, error) {
	log := logger.NewLogger() 

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	// ตั้งค่า Producer
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error("Error creating Kafka producer: %v", err) // ✅ ใช้ log.Error()
		return nil, err
	}

	// ตั้งค่า Consumer
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Error("Error creating Kafka consumer group: %v", err) // ✅ ใช้ log.Error()
		return nil, err
	}

	log.Info("Kafka client initialized successfully")

	return &KafkaClient{
		brokers:       brokers,
		topic:         topic,
		producer:      producer,
		consumerGroup: consumerGroup,
		log:           log, 
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
		kc.log.Error("Error publishing message to Kafka: %v", err) // ✅ ใช้ log.Error()
		return err
	}

	kc.log.Info("Message published successfully to Kafka")
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
				kc.log.Error("Error consuming messages: %v", err) // ✅ ใช้ log.Error()
			}
		}
	}()

	kc.log.Info("Kafka Consumer started...")
	return nil
}
