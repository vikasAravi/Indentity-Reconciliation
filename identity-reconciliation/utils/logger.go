package util

import (
	"bitespeed/identity-reconciliation/config"
	"github.com/google/uuid"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

var Log *Logger

type LoggerError struct {
	Error error
}

func panicIfError(err error) {
	if err != nil {
		panic(LoggerError{err})
	}
}

func SetupLogger() {
	level, err := logrus.ParseLevel(config.LogLevel())
	panicIfError(err)

	logrusVar := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: joonix.NewFormatter(),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	Log = &Logger{logrusVar}

}

func BuildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context":    context,
		"request_id": uuid.NewString(),
	}
}
