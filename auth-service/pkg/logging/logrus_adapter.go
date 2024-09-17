package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

type LogrusAdapter struct {
	logrus *logrus.Logger
}

func NewLogrusAdapter() Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))

	return &LogrusAdapter{
		logrus: logger,
	}
}

func (l *LogrusAdapter) LogInfo(message string) {
	l.logrus.Info(message)
}

func (l *LogrusAdapter) LogError(message string) {
	l.logrus.Error(message)
}

func (l *LogrusAdapter) LogWarn(message string) {
	l.logrus.Warn(message)
}

func (l *LogrusAdapter) LogDebug(message string) {
	l.logrus.Debug(message)
}

func (l *LogrusAdapter) LogPanic(message string) {
	l.logrus.Panic(message)
}
