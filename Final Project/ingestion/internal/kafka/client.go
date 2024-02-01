package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"ingestion/pkg/broker"
	"time"
)

const retries = 3

type kafkaProducer struct {
	writer  *kafka.Writer
	timeout time.Duration
}

func NewKafkaProducer(topic string, hosts []string, timeout int) broker.Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(hosts...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Balancer:               &StockTypePartitionBalancer{},
	}

	return &kafkaProducer{
		writer:  writer,
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (p *kafkaProducer) Produce(ctx context.Context, key string, value []byte) error {
	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	return p.produceWithRetry(ctx, message)
}

func (p *kafkaProducer) produceWithRetry(ctx context.Context, message kafka.Message) error {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(ctx, p.timeout)
		defer cancel()

		err = p.writer.WriteMessages(ctx, message)
		if err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
				time.Sleep(time.Millisecond * 250)
				continue
			}
			logrus.WithError(err).Error("error in produce")
			break
		}
	}

	return err
}
