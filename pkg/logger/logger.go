package logger

import "log"

type Logger struct {
}

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func NewLogger() LoggerInterface {
	return &Logger{}
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
