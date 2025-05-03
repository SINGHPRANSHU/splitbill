package handler

import (
	"log"
)

type StandardLogger struct {
	logger *log.Logger
}

// NewStandardLogger initializes a new StandardLogger
func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		logger: log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile),
	}
}

// Debug logs debug-level messages
func (l *StandardLogger) Debug(v ...interface{}) {
	l.logger.SetPrefix("DEBUG: ")
	l.logger.Println(v...)
}

// Fatal logs fatal-level messages and exits the application
func (l *StandardLogger) Fatal(v ...interface{}) {
	l.logger.SetPrefix("FATAL: ")
	l.logger.Fatalln(v...)
}

// Println logs general messages
func (l *StandardLogger) Println(v ...interface{}) {
	l.logger.SetPrefix("")
	l.logger.Println(v...)
}
