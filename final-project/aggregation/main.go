package main

import (
	"aggregation/internal/cache"
	"aggregation/internal/config"
	"aggregation/internal/repository"
	"aggregation/internal/service"
	"aggregation/pkg/redis"
	"fmt"
	_ "github.com/lib/pq"
	"pkg"
	"pkg/broker"
	"pkg/kafka"
	"pkg/logger"
	"pkg/sqlx"
)

func main() {
	logger.InitLogger()

	configs := config.LoadConfigs(pkg.IsLocalEnv())
	fmt.Println(configs)

	db := sqlx.NewDB(configs.Postgresql)
	aggregatedMetricsRepo := repository.NewAggregatedMetricsRepo(db)

	consumers := make([]broker.Consumer, 0)
	for _, partition := range kafka.StockDataSymbolToPartition {
		kafkaConfigs := configs.Kafka
		kafkaConfigs.Partition = partition
		consumers = append(consumers, kafka.NewKafkaConsumer(kafkaConfigs))
	}

	redisClient := redis.NewRedisClient(configs.Redis)

	metricsCache := cache.NewMetricsCache(redisClient)

	aggregationService := service.NewAggregationService(consumers, aggregatedMetricsRepo, metricsCache)
	aggregationService.Start()
}
