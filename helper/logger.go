package helper

import (
	"log"
	"os"
	"runtime"
	"strings"
)

type Logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
}

func NewLogger() *Logger {
	flags := log.LstdFlags
	return &Logger{
		infoLog:  log.New(os.Stdout, "\033[34m ✓ [INFO] ", flags),
		warnLog:  log.New(os.Stdout, "\033[33m ! [WARN] ", flags),
		errorLog: log.New(os.Stderr, "\033[31m ⨯ [ERROR] ", flags),
	}
}

func (l *Logger) Info(v ...any) {
	l.infoLog.Println(v...)
}

func (l *Logger) Warn(v ...any) {
	l.warnLog.Println(v...)
}

func (l *Logger) Error(v ...any) {
	_, filename, line, _ := runtime.Caller(1)
	lastSlashIdx := strings.LastIndex(filename, "/")
	filename = filename[lastSlashIdx+1:]
	l.errorLog.Printf("%s, line #%d, %+v", filename, line, v)
}
