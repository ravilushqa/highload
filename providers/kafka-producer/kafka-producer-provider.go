package kafkaproducerprovider

import (
	"github.com/Shopify/sarama"
)

// KafkaProducer kafka producer
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// New creates kafka producer
func New(addrs []string, topic string, kfkConfig *sarama.Config) (*KafkaProducer, error) {
	if kfkConfig == nil {
		kfkConfig = sarama.NewConfig()
		kfkConfig.Producer.Return.Successes = true
	}

	producer, err := sarama.NewSyncProducer(addrs, kfkConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Close closes producer
func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}

// SendMessage sends message
func (kp *KafkaProducer) SendMessage(payload, key []byte) error {
	kfkMessage := sarama.ProducerMessage{Topic: kp.topic, Value: sarama.ByteEncoder(payload)}

	if len(key) > 0 {
		kfkMessage.Key = sarama.ByteEncoder(key)
	}

	_, _, err := kp.producer.SendMessage(&kfkMessage)
	return err
}
