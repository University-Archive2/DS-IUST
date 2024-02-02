package main

import (
	"analyze/internal/config"
	"analyze/internal/repository"
	"analyze/internal/service"
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
	stockAnalysisRepo := repository.NewStockAnalysisRepo(db)

	consumers := make([]broker.Consumer, 0)
	for _, partition := range kafka.StockDataSymbolToPartition {
		kafkaConfigs := configs.Kafka
		kafkaConfigs.Partition = partition
		consumers = append(consumers, kafka.NewKafkaConsumer(kafkaConfigs))
	}

	analyzeService := service.NewAnalyzeService(consumers, stockAnalysisRepo)
	analyzeService.Start()
}
