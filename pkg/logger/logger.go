package logger

import (
	"github.com/moshrank/spacey-backend/config"
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

func NewLogger(
	cfg config.ConfigInterface,
) LoggerInterface {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(logLevelMapping[cfg.GetLogLevel()])

	logger := log.WithFields(
		log.Fields{
			"port": cfg.GetPort(),
		},
	)

	return logger
}
