package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var logLevelMapping = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
}

func NewLogger(logLevel string) LoggerInterface {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevelMapping[logLevel])

	return log.New()
}
