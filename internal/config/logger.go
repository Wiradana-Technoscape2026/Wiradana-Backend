package config

import "github.com/sirupsen/logrus"

func NewLogger(_ *Config) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)
	return log
}
