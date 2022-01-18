package logger

import (
	"io"
	"os"

	"github.com/moshrank/spacey-backend/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
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
	log.SetOutput(os.Stderr)

	logger := log.WithFields(
		log.Fields{
			"service": "spacey-backend",
			"port":    cfg.GetPort(),
		},
	)

	gelfWriter, err := gelf.NewUDPWriter(cfg.GetGrayLogConnection())

	if err != nil {
		log.Warn("Failed to connect to graylog: ", err)
	} else {
		log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
	}

	return logger
}
