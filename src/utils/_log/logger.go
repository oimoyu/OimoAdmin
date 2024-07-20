package _log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

type LoggerStruct struct {
	*log.Logger
	LogPrefix string

	FileLogPath string
}

func NewLogger(logPrefix string, fileLogPath string) *LoggerStruct {
	return &LoggerStruct{
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		LogPrefix:   logPrefix,
		FileLogPath: fileLogPath,
	}
}

func (l *LoggerStruct) Debug(format string, v ...interface{}) {
	l.SetPrefix(fmt.Sprintf("%s[%s DEBUG]: %s", Blue, l.LogPrefix, Reset))
	l.Printf(format, v...)
}

func (l *LoggerStruct) Error(format string, v ...interface{}) {
	l.SetPrefix(fmt.Sprintf("%s[%s ERROR]: %s", Red, l.LogPrefix, Reset))

	l.Printf(format, v...)
}

func (l *LoggerStruct) Info(format string, v ...interface{}) {
	l.SetPrefix(fmt.Sprintf("%s[%s INFO]: %s", Green, l.LogPrefix, Reset))
	l.Printf(format, v...)
}

func (l *LoggerStruct) Warn(format string, v ...interface{}) {
	l.SetPrefix(fmt.Sprintf("%s[%s WARN]: %s", Yellow, l.LogPrefix, Reset))

	l.Printf(format, v...)
}

func (l *LoggerStruct) FileLog(message string) error {
	file, err := os.OpenFile(l.FileLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_err := fmt.Errorf("failed to open log file: %w", err)
		l.Error(_err.Error())
		return _err
	}
	defer file.Close()

	sanitizedMessage := strings.ReplaceAll(message, "\n", " ")

	logMessage := fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), sanitizedMessage)
	if _, err := file.WriteString(logMessage); err != nil {
		_err := fmt.Errorf("failed to write log message: %w", err)
		l.Error(_err.Error())
		return _err
	}
	return nil
}
