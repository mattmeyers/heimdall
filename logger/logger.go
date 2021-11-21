package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Logger represents a standard level logging interface. Every method logs the provided
// message using fmt.Printf with a timestamp and level prefix. Fatal logs the message
// like the other methods, but calls os.Exit(1) afterwards.
type Logger interface {
	// Logs at the LevelDebug level.
	Debug(format string, args ...interface{})
	// Logs at the LevelInfo level.
	Info(format string, args ...interface{})
	// Logs at the LevelWarn level.
	Warn(format string, args ...interface{})
	// Logs at the LevelError level.
	Error(format string, args ...interface{})
	// Logs at the LevelFatal level then calls os.Exit(1).
	Fatal(format string, args ...interface{})
}

// Level represents a logging level. This restricts the logger to print only messages with
// at least this level.
type Level int

// The available logging levels.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}

	return fmt.Sprintf("%%!(Level=%d)", l)
}

// ParseLevel converts a string to the corresponding Level. Comparisons are case insensitive.
// If an unknown level is provided, then an error will be returned.
func ParseLevel(l string) (Level, error) {
	switch strings.ToLower(l) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	}

	return Level(-1), errors.New("invalid log level")
}

func (l Level) validate() error {
	if strings.HasPrefix("%!", l.String()) {
		return errors.New("invalid Level")
	}

	return nil
}

// LevelLogger implements the Logger interface using the defined Level constants. The provided
// level is treated as the minimum. Any messages passed to a level that is at least the defined
// level will be printed.
//
// Every log message is treated as a single line. If there is no newline at the end of the
// message, then one will be added.
type LevelLogger struct {
	w     io.Writer
	level Level
}

// NewLevelLogger constructs a new logger. An error will be returned if an invalid level is
// provided. If no output writer is provided, then os.Stdout will be used.
func NewLevelLogger(level Level, out io.Writer) (*LevelLogger, error) {
	if err := level.validate(); err != nil {
		return nil, err
	}

	if out == nil {
		out = os.Stdout
	}

	return &LevelLogger{
		w:     out,
		level: level,
	}, nil
}

// Logs at the LevelDebug level.
func (l *LevelLogger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.printPrefixTag(LevelDebug)
		l.printMessage([]byte(fmt.Sprintf(format, args...)))
	}
}

// Logs at the LevelInfo level.
func (l *LevelLogger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.printPrefixTag(LevelInfo)
		l.printMessage([]byte(fmt.Sprintf(format, args...)))
	}
}

// Logs at the LevelWarn level.
func (l *LevelLogger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.printPrefixTag(LevelWarn)
		l.printMessage([]byte(fmt.Sprintf(format, args...)))
	}
}

// Logs at the LevelError level.
func (l *LevelLogger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.printPrefixTag(LevelError)
		l.printMessage([]byte(fmt.Sprintf(format, args...)))
	}
}

// Logs at the LevelFatal level then calls os.Exit(1).
func (l *LevelLogger) Fatal(format string, args ...interface{}) {
	if l.level <= LevelFatal {
		l.printPrefixTag(LevelFatal)
		l.printMessage([]byte(fmt.Sprintf(format, args...)))
		os.Exit(1)
	}
}

func (l *LevelLogger) printPrefixTag(level Level) {
	l.w.Write([]byte(fmt.Sprintf("%s [%s]: ", time.Now().Format(time.RFC3339), level)))
}

var newline = []byte{'\n'}

func (l *LevelLogger) printMessage(message []byte) {
	l.w.Write(message)
	if len(message) == 0 || message[len(message)-1] != '\n' {
		l.w.Write(newline)
	}
}
