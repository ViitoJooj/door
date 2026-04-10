package logger

import (
	"os"
	"sync"
	"time"
)

type Logger struct {
	out   *os.File
	mutex sync.Mutex
}

func NewLogger(out *os.File) *Logger {
	return &Logger{out: out}
}

func (l *Logger) write(level, msg string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	row := `{"time":"` + now + `","level":"` + level + `","msg":"` + msg + `"}` + "\n"

	l.out.WriteString(row)
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg)
}

func (l *Logger) Warn(msg string) {
	l.write("WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
}
