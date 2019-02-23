package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
	}
}

//noinspection GoUnusedExportedFunction
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

//noinspection GoUnusedExportedFunction
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

//noinspection GoUnusedExportedFunction
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

//noinspection GoUnusedExportedFunction
func Info(args ...interface{}) {
	logger.Info(args...)
}

//noinspection GoUnusedExportedFunction
func Error(args ...interface{}) {
	logger.Error(args...)
}
