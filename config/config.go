package config

import (
	"github.com/Shopify/sarama"
	"time"
)

func GetKafkaConfig(username, password string, timeout, retry int) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = time.Duration(timeout) * time.Second
	kafkaConfig.Producer.Retry.Max = retry

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}
