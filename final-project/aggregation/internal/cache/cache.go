package cache

import (
	"aggregation/pkg/redis"
	"context"
)

type MetricsCache struct {
	cache redis.Cache
}

func NewMetricsCache(cache redis.Cache) *MetricsCache {
	return &MetricsCache{
		cache: cache,
	}
}

func (m *MetricsCache) Set(key string, value []byte) error {
	return m.cache.Set(context.Background(), key, value)
}

func (m *MetricsCache) Get(key string) ([]byte, error) {
	return m.cache.Get(context.Background(), key)
}
