package sqlx

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PostgresqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func NewDB(config PostgresqlConfig) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user="+config.Username+" dbname="+config.DBName+" password="+config.Password+" host="+config.Host+" sslmode="+config.SSLMode)
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to db")
	}

	return db
}
