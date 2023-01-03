package utils

import (
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	reset  = "\033[0m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
)

type Logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
}

func NewLogger() *Logger {
	flags := log.LstdFlags
	return &Logger{
		infoLog:  log.New(os.Stdout, blue+" ✓ [INFO] ", flags),
		warnLog:  log.New(os.Stdout, yellow+" ! [WARN] ", flags),
		errorLog: log.New(os.Stderr, red+" ⨯ [ERROR] ", flags),
	}
}

func (l *Logger) Info(v ...any) {
	filename, line := l.moreInfo()
	l.infoLog.Printf("%s:%d - %+v %s\n", filename, line, v, reset)
}

func (l *Logger) Warn(v ...any) {
	filename, line := l.moreInfo()
	l.warnLog.Printf("%s:%d - %+v %s\n", filename, line, v, reset)
}

func (l *Logger) Error(v ...any) {
	filename, line := l.moreInfo()
	l.errorLog.Printf("%s:%d - %+v %s\n", filename, line, v, reset)
}

func (l *Logger) moreInfo() (string, int) {
	_, filename, line, _ := runtime.Caller(2)
	lastSlashIdx := strings.LastIndex(filename, "/")
	filename = filename[lastSlashIdx+1:]
	return filename, line
}
