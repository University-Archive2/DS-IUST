package broker

import "context"

type Producer interface {
	Produce(ctx context.Context, key string, value []byte) error
}
