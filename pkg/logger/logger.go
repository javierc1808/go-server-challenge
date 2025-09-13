package logger

import (
	"log"
	"os"
)

// Logger defines the interface for logging
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

// SimpleLogger implements Logger using Go's standard logger
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewSimpleLogger creates a new instance of SimpleLogger
func NewSimpleLogger() Logger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.LstdFlags),
	}
}

// Info logs an informational message
func (l *SimpleLogger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, err error) {
	if err != nil {
		l.errorLogger.Printf("%s: %v", msg, err)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string) {
	l.debugLogger.Println(msg)
}
