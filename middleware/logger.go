package middleware

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

var requestLogger *log.Logger

// InitRequestLogger initializes the request logger
func InitRequestLogger() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		// Use console only if can't create log file
		requestLogger = log.New(os.Stdout, "[REQUEST] ", log.Ldate|log.Ltime)
		return
	}

	// Open log file
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		// Use console only if can't open log file
		requestLogger = log.New(os.Stdout, "[REQUEST] ", log.Ldate|log.Ltime)
		return
	}

	// Create multi writer (console + file)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	requestLogger = log.New(multiWriter, "[REQUEST] ", log.Ldate|log.Ltime)
}

func LoggerMiddleware(c *fiber.Ctx) error {
	// Initialize logger if not already done
	if requestLogger == nil {
		InitRequestLogger()
	}
	
	start := time.Now()
	
	// Process request
	err := c.Next()
	
	// Log request details
	duration := time.Since(start)
	status := c.Response().StatusCode()
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	userAgent := c.Get("User-Agent")
	
	// Format log message
	requestLogger.Printf("%s %s [%d] - %v - %s - %s", 
		method, path, status, duration, ip, userAgent)
	
	return err
}
