package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Cache interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(conn RedisConfig) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     conn.Host + ":" + conn.Port,
		Password: conn.Password, // no password set
		DB:       conn.DB,       // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to redis")
	}

	return &redisClient{
		client: client,
	}
}

func (r redisClient) Set(ctx context.Context, key string, value []byte) error {
	err := r.client.Set(ctx, key, string(value), 0).Err()
	if err != nil {
		logrus.WithError(err).Error("failed to set value in redis")
	}

	return err
}

func (r redisClient) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		logrus.WithError(err).Error("failed to get value from redis")
		return nil, err
	}

	return []byte(result), nil
}
