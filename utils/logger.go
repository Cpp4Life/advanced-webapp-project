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
	l.infoLog.Println(v, reset)
}

func (l *Logger) Warn(v ...any) {
	l.warnLog.Println(v, reset)
}

func (l *Logger) Error(v ...any) {
	_, filename, line, _ := runtime.Caller(1)
	lastSlashIdx := strings.LastIndex(filename, "/")
	filename = filename[lastSlashIdx+1:]
	l.errorLog.Printf("%s, line #%d, %+v %s\n", filename, line, v, reset)
}
