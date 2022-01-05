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
	logLevel, graylogConnection string,
	curConfig config.ConfigInterface,
) LoggerInterface {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(logLevelMapping[logLevel])
	log.SetOutput(os.Stderr)

	logger := log.WithFields(
		log.Fields{
			"service":                   "spacey-backend",
			"port":                      curConfig.GetPort(),
			"user-service-db-name":      curConfig.GetUserSeviceDBNAME(),
			"flashcard-service-db-name": curConfig.GetFlashcardServiceDBName(),
		},
	)

	gelfWriter, err := gelf.NewUDPWriter(graylogConnection)

	if err != nil {
		log.Warn("Failed to connect to graylog: ", err)
	} else {
		log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
	}

	return logger
}
