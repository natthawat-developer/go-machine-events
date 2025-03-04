package logger

import (
	"log"
	"os"
)

// Logger เป็นโครงสร้างที่จัดการ Logging
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// **NewLogger สร้าง Logger ใหม่**
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// **Info ใช้สำหรับ Log ข้อความทั่วไป**
func (l *Logger) Info(message string, args ...interface{}) {
	l.infoLogger.Printf(message, args...)
}

// **Error ใช้สำหรับ Log ข้อความ Error**
func (l *Logger) Error(message string, args ...interface{}) {
	l.errorLogger.Printf(message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.errorLogger.Printf(message, args...)
	os.Exit(1) // ออกจากโปรแกรม
}