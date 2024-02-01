package main

import (
	"fmt"
	"ingestion/internal/api"
	"ingestion/internal/config"
	"ingestion/internal/service"
	"pkg"
	"pkg/kafka"
	"pkg/logger"
)

func main() {
	logger.InitLogger()

	configs := config.LoadConfigs(pkg.IsLocalEnv())
	fmt.Println(configs)

	kafkaProducer := kafka.NewKafkaProducer(configs.KafkaHosts, 2)

	ingestionService := service.NewIngestionService(kafkaProducer)

	ingestionAPI := api.NewIngestionAPI(ingestionService)
	ingestionAPI.Start()
}
