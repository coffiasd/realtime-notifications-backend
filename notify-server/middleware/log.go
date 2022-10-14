package middleware

import (
	"notify-server/config"
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func GetLogger() *logrus.Logger {
	return logger
}

func InitLog() {
	logger.Formatter = &logrus.JSONFormatter{}
	logFile, _ := os.OpenFile(config.ParseConfig.Log.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger.SetOutput(logFile)
}
