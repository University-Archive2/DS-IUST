package logger

import "github.com/sirupsen/logrus"

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
