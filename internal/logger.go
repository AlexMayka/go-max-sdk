package internal

import (
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// ConsoleLogger is the default logger implementation
type ConsoleLogger struct {
	mu     sync.RWMutex
	level  core.LogLevel
	logger *log.Logger
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger() core.Logger {
	return &ConsoleLogger{
		level:  core.INFO,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *ConsoleLogger) SetLevel(level core.LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *ConsoleLogger) shouldLog(level core.LogLevel) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return level >= l.level
}

func (l *ConsoleLogger) log(level core.LogLevel, component, event, details string) {
	if !l.shouldLog(level) {
		return
	}

	levelStr := l.levelToString(level)
	timestamp := time.Now().Format("15:04:05")

	// Format: [LEVEL] timestamp component: event - details
	message := fmt.Sprintf("[%s] %s %s: %s", levelStr, timestamp, component, event)
	if details != "" {
		message += fmt.Sprintf(" - %s", details)
	}

	l.logger.Println(message)
}

func (l *ConsoleLogger) levelToString(level core.LogLevel) string {
	switch level {
	case core.DEBUG:
		return "DEBUG"
	case core.INFO:
		return "INFO"
	case core.WARN:
		return "WARN"
	case core.ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *ConsoleLogger) Debug(component, event, details string) {
	l.log(core.DEBUG, component, event, details)
}

func (l *ConsoleLogger) Info(component, event, details string) {
	l.log(core.INFO, component, event, details)
}

func (l *ConsoleLogger) Warn(component, event, details string) {
	l.log(core.WARN, component, event, details)
}

func (l *ConsoleLogger) Error(component, event, details string) {
	l.log(core.ERROR, component, event, details)
}

// NoopLogger is a logger that does nothing
type NoopLogger struct{}

func NewNoopLogger() core.Logger {
	return &NoopLogger{}
}

func (n *NoopLogger) Debug(component, event, details string) {}
func (n *NoopLogger) Info(component, event, details string)  {}
func (n *NoopLogger) Warn(component, event, details string)  {}
func (n *NoopLogger) Error(component, event, details string) {}
func (n *NoopLogger) SetLevel(level core.LogLevel)           {}

// ParseLogLevel parses string log level
func ParseLogLevel(s string) core.LogLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return core.DEBUG
	case "INFO":
		return core.INFO
	case "WARN":
		return core.WARN
	case "ERROR":
		return core.ERROR
	default:
		return core.INFO
	}
}
