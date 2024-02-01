package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	KafkaHosts []string `mapstructure:"kafka_hosts"`
}

func LoadConfigs(isLocal bool) Config {
	if isLocal {
		loadLocalConfigs()
	}

	viper.AutomaticEnv()

	viper.SetDefault("kafka_hosts", "localhost:9092")

	var data Config

	err := viper.Unmarshal(&data)
	if err != nil {
		logrus.WithError(err).Fatal("error in unmarshalling config")
	}

	return data
}

func loadLocalConfigs() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("error in reading config")
	}
}
