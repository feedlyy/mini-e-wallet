package consumer

import "github.com/Shopify/sarama"

// KafkaConsumer hold sarama consumer
type KafkaConsumer struct {
	Consumer sarama.Consumer
}
