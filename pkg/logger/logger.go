package logger

import "log"

type Logger struct {
	LogLevel string
}

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func NewLogger(logLevel string) LoggerInterface {
	return &Logger{
		LogLevel: logLevel,
	}
}

func (l *Logger) Debug(args ...interface{}) {
	log.Println(args)
}

func (l *Logger) Info(args ...interface{}) {
	log.Println(args)
}

func (l *Logger) Warn(args ...interface{}) {
	log.Println(args)
}

func (l *Logger) Error(args ...interface{}) {
	log.Println(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	log.Fatal(args)
}
