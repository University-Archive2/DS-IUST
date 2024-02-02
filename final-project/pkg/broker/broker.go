package broker

import "context"

type Message struct {
	Key   string
	Value []byte
}

type Producer interface {
	Produce(ctx context.Context, message *Message) error
}

type Consumer interface {
	Consume(messagesChan chan<- *Message)
}
