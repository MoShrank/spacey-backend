package logger

import (
	"github.com/moshrank/spacey-backend/config"
	log "github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	WithFields(fields log.Fields) *log.Entry
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var logLevelMapping = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
	"FATAL": log.FatalLevel,
	"PANIC": log.PanicLevel,
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
