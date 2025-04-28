package util

import (
	"fmt"
	"os"
	"strings"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	// LogLevelError only logs errors
	LogLevelError LogLevel = iota
	// LogLevelInfo logs errors and important info
	LogLevelInfo
	// LogLevelDebug logs everything including debug messages
	LogLevelDebug
)

var (
	// CurrentLogLevel controls the verbosity of logging
	CurrentLogLevel LogLevel = LogLevelInfo
)

// SetLogLevel sets the current logging level
func SetLogLevel(level LogLevel) {
	CurrentLogLevel = level
}

// GetLogLevelFromEnv reads log level from environment variable
func GetLogLevelFromEnv() {
	levelStr := os.Getenv("STARLINK_LOG_LEVEL")
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		SetLogLevel(LogLevelDebug)
	case "ERROR":
		SetLogLevel(LogLevelError)
	default:
		SetLogLevel(LogLevelInfo)
	}
}

// LogDebug prints a message if the current log level is Debug or higher
func LogDebug(format string, args ...interface{}) {
	if CurrentLogLevel >= LogLevelDebug {
		fmt.Printf(format, args...)
	}
}

// LogInfo prints a message if the current log level is Info or higher
func LogInfo(format string, args ...interface{}) {
	if CurrentLogLevel >= LogLevelInfo {
		fmt.Printf(format, args...)
	}
}

// LogError prints a message if the current log level is Error or higher (always prints)
func LogError(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}