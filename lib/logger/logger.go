package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger provides structured logging to a file
type Logger struct {
	file      *os.File
	startTime time.Time
	mu        sync.Mutex
}

// GetInstance returns the singleton logger instance
func GetInstance() *Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

// newLogger creates a new logger instance
func newLogger() *Logger {
	startTime := time.Now()

	// Create logs directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get home directory: %v\n", err)
		return &Logger{startTime: startTime}
	}

	logsDir := filepath.Join(homeDir, ".vibecast", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logs directory: %v\n", err)
		return &Logger{startTime: startTime}
	}

	// Create log file with timestamp
	timestamp := startTime.Format("2006-01-02_15-04-05")
	logFile := filepath.Join(logsDir, fmt.Sprintf("%s.log", timestamp))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
		return &Logger{startTime: startTime}
	}

	logger := &Logger{
		file:      file,
		startTime: startTime,
	}

	logger.Info("Logger initialized", "log_file", logFile)
	return logger
}

// log writes a log entry with the given level and message
func (l *Logger) log(level, msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Build log line
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level, msg)

	// Add fields if present
	if len(fields) > 0 {
		logLine += " {"
		first := true
		for key, value := range fields {
			if !first {
				logLine += ", "
			}
			logLine += fmt.Sprintf("%s=%v", key, value)
			first = false
		}
		logLine += "}"
	}

	logLine += "\n"

	l.file.WriteString(logLine)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logWithFields("INFO", msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logWithFields("ERROR", msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logWithFields("DEBUG", msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logWithFields("WARN", msg, fields...)
}

// logWithFields handles alternating string keys and values
func (l *Logger) logWithFields(level, msg string, fields ...interface{}) {
	fieldMap := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			fieldMap[key] = fields[i+1]
		}
	}
	l.log(level, msg, fieldMap)
}

// LogError is a convenience function to log an error with context
func (l *Logger) LogError(context string, err error) {
	if err != nil {
		l.Error("error_occurred", "context", context, "error", err.Error())
	}
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GetLogFilePath returns the path to the current log file
func (l *Logger) GetLogFilePath() string {
	if l.file != nil {
		return l.file.Name()
	}
	return ""
}
