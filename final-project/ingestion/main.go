package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"ingestion/internal/api"
	"ingestion/internal/config"
	"ingestion/internal/kafka"
	"ingestion/internal/service"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	user := os.Getenv("USER")
	configs := config.LoadConfigs(user == "divar")
	fmt.Println(configs)

	kafkaProducer := kafka.NewKafkaProducer("stock", configs.KafkaHosts, 2)

	ingestionService := service.NewIngestionService(kafkaProducer)

	ingestionAPI := api.NewIngestionAPI(ingestionService)
	ingestionAPI.Start()
}
