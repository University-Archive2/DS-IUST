package kafka

import "time"

type KafkaConfig struct {
	Hosts         []string      `mapstructure:"hosts"`
	ConsumerGroup string        `mapstructure:"consumer_group"`
	Timeout       time.Duration `mapstructure:"timeout"`
	Topic         string        `mapstructure:"topic"`
	Partition     int           `mapstructure:"partition"`
}
