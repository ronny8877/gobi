package main

import (
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

// LoggerStruct represents the structured logger
type LoggerStruct struct {
	info  func(msg string, args ...interface{})
	warn  func(msg string, args ...interface{})
	err   func(msg string, args ...interface{})
	debug func(msg string, args ...interface{})
}

// Logger initializes and returns a LoggerStruct
func (app *App) Logger() *LoggerStruct {
	logger := logrus.New()
	logger.SetOutput(colorable.NewColorableStdout())
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
	})

	return &LoggerStruct{
		info: func(msg string, args ...interface{}) {
			if app.Config.Logging != nil && *app.Config.Logging {
				logger.Infof(msg, args...)
			}
		},
		warn: func(msg string, args ...interface{}) {
			if app.Config.Logging != nil && *app.Config.Logging {
				logger.Warnf(msg, args...)
			}
		},
		err: func(msg string, args ...interface{}) {
			if app.Config.Logging != nil && *app.Config.Logging {
				logger.Errorf(msg, args...)
			}
		},
		debug: func(msg string, args ...interface{}) {
			if app.Config.Logging != nil && *app.Config.Logging {
				logger.Debugf(msg, args...)
			}
		},
	}
}
