package middleware

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

var requestLogger *log.Logger

func InitRequestLogger() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		requestLogger = log.New(os.Stdout, "[REQUEST] ", log.Ldate|log.Ltime)
		return
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		requestLogger = log.New(os.Stdout, "[REQUEST] ", log.Ldate|log.Ltime)
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	requestLogger = log.New(multiWriter, "[REQUEST] ", log.Ldate|log.Ltime)
}

func LoggerMiddleware(c *fiber.Ctx) error {
	if requestLogger == nil {
		InitRequestLogger()
	}
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)
	status := c.Response().StatusCode()
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	userAgent := c.Get("User-Agent")
	
	requestLogger.Printf("%s %s [%d] - %v - %s - %s", 
		method, path, status, duration, ip, userAgent)
	
	return err
}
