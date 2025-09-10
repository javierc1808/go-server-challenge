package logger

import (
	"log"
	"os"
)

// Logger define la interfaz para logging
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

// SimpleLogger implementa Logger con el logger estándar de Go
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewSimpleLogger crea una nueva instancia de SimpleLogger
func NewSimpleLogger() Logger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.LstdFlags),
	}
}

// Info registra un mensaje de información
func (l *SimpleLogger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Error registra un mensaje de error
func (l *SimpleLogger) Error(msg string, err error) {
	if err != nil {
		l.errorLogger.Printf("%s: %v", msg, err)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Debug registra un mensaje de debug
func (l *SimpleLogger) Debug(msg string) {
	l.debugLogger.Println(msg)
}
