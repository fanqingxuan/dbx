package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
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
	default:
		return "UNKNOWN"
	}
}

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	SetLevel(level Level)
}

type defaultLogger struct {
	level  Level
	output io.Writer
}

func New() Logger {
	return &defaultLogger{level: LevelInfo, output: os.Stdout}
}

func NewWithOutput(w io.Writer) Logger {
	return &defaultLogger{level: LevelInfo, output: w}
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) log(level Level, format string, args ...any) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, msg)
}

func (l *defaultLogger) Debug(format string, args ...any) { l.log(LevelDebug, format, args...) }
func (l *defaultLogger) Info(format string, args ...any)  { l.log(LevelInfo, format, args...) }
func (l *defaultLogger) Warn(format string, args ...any)  { l.log(LevelWarn, format, args...) }
func (l *defaultLogger) Error(format string, args ...any) { l.log(LevelError, format, args...) }

var std = New()

func SetLevel(level Level)                    { std.SetLevel(level) }
func Debug(format string, args ...any)        { std.Debug(format, args...) }
func Info(format string, args ...any)         { std.Info(format, args...) }
func Warn(format string, args ...any)         { std.Warn(format, args...) }
func Error(format string, args ...any)        { std.Error(format, args...) }
