package config

import (
	"aggregation/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pkg/kafka"
	"pkg/sqlx"
	"strings"
)

type Config struct {
	Kafka      kafka.KafkaConfig     `mapstructure:"kafka"`
	Postgresql sqlx.PostgresqlConfig `mapstructure:"postgresql"`
	Redis      redis.RedisConfig     `mapstructure:"redis"`
}

func LoadConfigs(isLocal bool) Config {
	if isLocal {
		loadLocalConfigs()
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	viper.AutomaticEnv()

	viper.SetDefault("kafka.hosts", "localhost:9092")
	viper.SetDefault("kafka.consumer_group", "analyze")
	viper.SetDefault("kafka.timeout", "2s")
	viper.SetDefault("kafka.topic", "stock")
	viper.SetDefault("kafka.partition", 0)

	viper.SetDefault("postgresql.host", "localhost")
	viper.SetDefault("postgresql.port", 5432)
	viper.SetDefault("postgresql.username", "postgres")
	viper.SetDefault("postgresql.password", "postgres")
	viper.SetDefault("postgresql.dbname", "postgres")
	viper.SetDefault("postgresql.sslmode", "disable")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

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
