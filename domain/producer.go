package domain

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// KafkaProducer hold kafka producer session
type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// SendMessage function to send message into kafka
func (p *KafkaProducer) SendMessage(topic string, t Transaction) error {
	var (
		kafkaMsgs     *sarama.ProducerMessage
		err           error
		serializedMsg []byte
	)

	serializedMsg, err = json.Marshal(t)
	if err != nil {
		logrus.Errorf("marhshal-ing err when send message kafka: %v", err)
		return err
	}

	kafkaMsgs = &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(serializedMsg),
	}

	_, _, err = p.Producer.SendMessage(kafkaMsgs)
	if err != nil {
		logrus.Errorf("Send message error: %v", err)
		return err
	}

	return nil
}
