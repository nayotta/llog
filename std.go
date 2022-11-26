package llog

import (
	"fmt"
	"io"
	"os"
)

var std *Logger

func init() {
	std = &Logger{
		out: &mutexWriter{
			Writer: os.Stderr,
		},
	}
}

func Default() *Logger {
	return std
}

func SetTag(tag string) {
	std.tag = tag
}

func SetOutput(out io.Writer) {
	std.out = &mutexWriter{
		Writer: out,
	}
}

func SetLevelString(s string) {
	std.setLevelString(s)
}

func SetLevel(level Level) {
	std.setLevel(level)
}

func SetFileAndLine(included bool) {
	std.fileAndLine = included
}

func WithTag(tag string) *Logger {
	return std.WithTag(tag)
}

func WithLevel(level Level) *Logger {
	return std.WithLevel(level)
}

func WithOutput(out io.Writer) *Logger {
	return std.WithOutput(out)
}

func WithFileAndLine(included bool) *Logger {
	return std.WithFileAndLine(included)
}

func Error(v ...any) {
	std.output(LevelError, fmt.Sprint(v...))
}

func Errorf(format string, v ...any) {
	std.output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	std.output(LevelWarning, fmt.Sprint(v...))
}

func Warnf(format string, v ...any) {
	std.output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	std.output(LevelInfo, fmt.Sprint(v...))
}

func Infof(format string, v ...any) {
	std.output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(v ...any) {
	std.output(LevelDebug, fmt.Sprint(v...))
}

func Debugf(format string, v ...any) {
	std.output(LevelDebug, fmt.Sprintf(format, v...))
}

func Fatal(v ...any) {
	std.output(LevelError, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	std.output(LevelError, fmt.Sprintf(format, v...))
	os.Exit(1)
}
