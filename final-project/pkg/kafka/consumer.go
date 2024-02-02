package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"pkg/broker"
	"time"
)

type kafkaConsumer struct {
	reader  *kafka.Reader
	timeout time.Duration
}

func NewKafkaConsumer(config KafkaConfig) broker.Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   config.Hosts,
		Topic:     config.Topic,
		Partition: config.Partition,
	})

	return &kafkaConsumer{
		reader:  reader,
		timeout: config.Timeout,
	}
}

func (c *kafkaConsumer) Consume(messagesChan chan<- *broker.Message) {
	for {
		message, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			logrus.WithError(err).Error("error in consume")
			continue
		}

		messagesChan <- &broker.Message{
			Key:   string(message.Key),
			Value: message.Value,
		}
	}
}
