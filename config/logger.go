package config

import (
	"io"
	"log"
	"os"
)

var AppLogger *log.Logger

func InitLogger() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	
	AppLogger = log.New(multiWriter, "[ALUMNI-API] ", log.Ldate|log.Ltime|log.Lshortfile)
	
	AppLogger.Println("=== Alumni API Started ===")
}

func GetLogger() *log.Logger {
	if AppLogger == nil {
		InitLogger()
	}
	return AppLogger
}
