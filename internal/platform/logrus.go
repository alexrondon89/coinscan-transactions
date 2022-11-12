package platform

import (
	"github.com/sirupsen/logrus"
)

func NewLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.Info("logrus created")
	return logger
}
