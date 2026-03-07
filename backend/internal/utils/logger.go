package utils

import (
	"log"
	"os"
)

// Logger 日志记录器
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLogger 创建新的日志记录器
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info 记录信息日志
func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Warn 记录警告日志
func (l *Logger) Warn(msg string) {
	l.warnLogger.Println(msg)
}

// Error 记录错误日志
func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

// Fatal 记录致命错误并退出
func (l *Logger) Fatal(msg string) {
	l.errorLogger.Fatal(msg)
}
