package utils

import (
	"log"
	"os"
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
	flags := log.LstdFlags | log.Lshortfile
	return &Logger{
		infoLog:  log.New(os.Stdout, blue+" ✓ [INFO] ", flags),
		warnLog:  log.New(os.Stdout, yellow+" ! [WARN] ", flags),
		errorLog: log.New(os.Stderr, red+" ⨯ [ERROR] ", flags),
	}
}

func (l *Logger) Info(v ...any) {
	l.infoLog.Printf("%+v %s\n", v, reset)
}

func (l *Logger) Warn(v ...any) {
	l.warnLog.Printf("%+v %s\n", v, reset)
}

func (l *Logger) Error(v ...any) {
	l.errorLog.Printf("%+v %s\n", v, reset)
}
