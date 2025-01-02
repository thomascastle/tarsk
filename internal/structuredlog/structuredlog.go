package structuredlog

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level uint8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	minLevel Level
	mu       sync.Mutex
	out      io.Writer
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		minLevel: minLevel,
		out:      out,
	}
}

func (l *Logger) Error(e error, properties map[string]string) {
	l.print(LevelError, e.Error(), properties)
}

func (l *Logger) Fatal(e error, properties map[string]string) {
	l.print(LevelFatal, e.Error(), properties)
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Write(message []byte) (int, error) {
	return l.print(LevelError, string(message), nil)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	entry := struct {
		Level      string            `json:"level"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Time       string            `json:"time"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Message:    message,
		Properties: properties,
		Time:       time.Now().UTC().Format(time.RFC3339),
	}

	if level >= LevelError {
		entry.Trace = string(debug.Stack())
	}

	var line []byte

	line, e := json.Marshal(entry)
	if e != nil {
		line = []byte(LevelError.String() + ": unable to marshal log messaeg: " + e.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
