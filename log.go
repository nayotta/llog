package llog

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var bufPool sync.Pool

func init() {
	bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, 1024)
			return &buf
		},
	}
}

type mutexWriter struct {
	mu sync.Mutex
	io.Writer
}

func (w *mutexWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.Writer.Write(b)
}

// Level defines what logs should be printed
type Level int

// LogLevels
const (
	LevelError   Level = -2
	LevelWarning Level = -1
	LevelInfo    Level = 0 // default log level
	LevelDebug   Level = 1
)

var levelString = map[Level]string{
	LevelError:   "[E]",
	LevelWarning: "[W]",
	LevelInfo:    "[I]",
	LevelDebug:   "[D]",
}

func (level Level) String() string {
	return levelString[level]
}

// Logger is a simple custom logger support log levels
type Logger struct {
	out *mutexWriter

	level       Level
	tag         string
	fileAndLine bool
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) setLevel(level Level) {
	l.level = level
}

func (l *Logger) setLevelString(s string) {
	switch strings.ToLower(s) {
	case "error", "e":
		l.setLevel(LevelError)
	case "warning", "w":
		l.setLevel(LevelWarning)
	case "info", "i":
		l.setLevel(LevelInfo)
	case "debug", "d":
		l.setLevel(LevelDebug)
	}
}

func (l *Logger) formatHeader(buf *[]byte, level Level) {
	ts := time.Now().Format("2006/01/02 15:04:05.000 ")
	*buf = append(*buf, ts...)

	ls := level.String()
	*buf = append(*buf, ls...)

	if l.tag != "" {
		*buf = append(*buf, '[')
		*buf = append(*buf, l.tag...)
		*buf = append(*buf, ']', ' ')
	}

	if l.fileAndLine {
		var ok bool
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			file = "???"
			line = 0
		}

		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		nu := strconv.Itoa(line)
		*buf = append(*buf, nu...)
		*buf = append(*buf, ' ')
	}
}

func (l *Logger) output(level Level, s string) {
	if level > l.level {
		return
	}

	buf := bufPool.New().(*[]byte)
	defer bufPool.Put(buf)

	*buf = (*buf)[:0]
	l.formatHeader(buf, level)
	*buf = append(*buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	_, err := l.out.Write(*buf)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) clone() *Logger {
	return &Logger{
		out:         l.out,
		tag:         l.tag,
		level:       l.level,
		fileAndLine: l.fileAndLine,
	}
}

func (l *Logger) WithTag(tag string) *Logger {
	clone := l.clone()
	clone.tag = tag
	return clone
}

func (l *Logger) WithLevel(level Level) *Logger {
	clone := l.clone()
	clone.level = level
	return clone
}

func (l *Logger) WithOutput(out io.Writer) *Logger {
	clone := l.clone()
	clone.out = &mutexWriter{
		Writer: out,
	}
	return clone
}

func (l *Logger) WithFileAndLine(included bool) *Logger {
	clone := l.clone()
	clone.fileAndLine = included
	return clone
}

func (l *Logger) Error(v ...any) {
	l.output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...any) {
	l.output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...any) {
	l.output(LevelWarning, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...any) {
	l.output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...any) {
	l.output(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...any) {
	l.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...any) {
	l.output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...any) {
	l.output(LevelDebug, fmt.Sprintf(format, v...))
}
