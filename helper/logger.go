package helper

import (
	"log"
	"os"
	"runtime"
	"strings"
)

type logger struct {
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
}

func NewLogger() *logger {
	flags := log.LstdFlags
	return &logger{
		infoLog:  log.New(os.Stdout, "\033[34m ✓ [INFO] ", flags),
		warnLog:  log.New(os.Stdout, "\033[33m ! [WARN] ", flags),
		errorLog: log.New(os.Stderr, "\033[31m ⨯ [ERROR] ", flags),
	}
}

func (l *logger) Info(v ...interface{}) {
	l.infoLog.Println(v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.warnLog.Println(v...)
}

func (l *logger) Error(v ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	lastSlashIdx := strings.LastIndex(filename, "/")
	filename = filename[lastSlashIdx+1:]
	l.errorLog.Printf("%s, line #%d, %+v", filename, line, v)
}
