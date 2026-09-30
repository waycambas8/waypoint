package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// Init initializes the loggers.
func Init() {
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}

// Info logs an informational message.
func Info(msg string) {
	if infoLogger == nil {
		Init()
	}
	infoLogger.Println(msg)
}

// Error logs an error message.
func Error(msg string) {
	if errorLogger == nil {
		Init()
	}
	errorLogger.Println(msg)
}

// Success logs a success message.
func Success(msg string) {
	fmt.Printf("[OK] %s\n", msg)
}
