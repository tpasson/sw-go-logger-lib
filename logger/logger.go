package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogSeverity int

const (
	NONE LogSeverity = iota
	PRODUCTION
	DEBUG
)

type Logger struct {
	Severity LogSeverity
	LogChan  chan Container
}

type Container struct {
	PreText        string
	Id             string
	Source         string
	Info           string
	Data           string
	Error          string
	ProcessingTime time.Duration
	Timestamp      time.Time
	HttpRequest    *http.Request
	ProcessedData  *interface{}
}

var severityMap = map[string]LogSeverity{
	"NONE":       NONE,
	"PRODUCTION": PRODUCTION,
	"DEBUG":      DEBUG,
}



// NewLogger creates a new Logger instance with the specified severity.
// Valid log severities are: NONE, PRODUCTION, DEBUG.
func NewLogger(severity string) (*Logger, error) {
	logSeverity, ok := severityMap[severity]
	if !ok {
		return nil, fmt.Errorf("invalid log severity: %s. valid log severities: NONE, PRODUCTION, DEBUG", severity)
	}

	logger := &Logger{
		Severity: logSeverity,
		LogChan:  make(chan Container),
	}

	go logger.processLogs()

	return logger, nil
}



// Entry logs a message based on the severity level and the provided container.
func (l *Logger) Entry(c Container) {
	if l.Severity == NONE {
		return
	}

	if c.Timestamp.IsZero() {
		c.Timestamp = generateTimestamp()
	}

	l.LogChan <- c
}



// generateTimestamp creates the current timestamp
func generateTimestamp() time.Time {
	return time.Now()
}



// processLogs handles log entries asynchronously
func (l *Logger) processLogs() {
	for c := range l.LogChan {
		// Process the log entry
		if l.Severity == NONE {
			continue
		}

		// Create buffer
		var result strings.Builder

		result.WriteString(formatTimestamp(c.Timestamp))

		// Append fields from the container to the log message
		appendFields(&result, c)

		// Append error message, if present, to the log message
		appendError(&result, c)

		// Append HTTP request details, if present, to the log message
		appendHttpRequest(&result, c)

		// Append processing time, if non-zero, to the log message
		appendProcessingTime(&result, c)

		switch l.Severity {
		case PRODUCTION:
			// Log the message for production severity
			// Do nothing
		case DEBUG:
			// Log the message with additional JSON data for debug severity
			logWithJson(result.String(), c)
		}

		writeLog(result.String(), &c)
	}
}



// formatTimestamp formats the current timestamp and returns it as a string
func formatTimestamp(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
}



// appendFields appends non-empty fields from the container to the log message
func appendFields(result *strings.Builder, c Container) {
	fields := []string{
		c.PreText,
		c.Id,
		c.Source,
		c.Info,
		c.Data,
	}

	for _, field := range fields {
		if field != "" {
			result.WriteString(" " + field)
		}
	}
}



// appendError appends the error message with color formatting to the log message
func appendError(result *strings.Builder, c Container) {
	if c.Error != "" {
		result.WriteString(" \033[31mError: " + c.Error + "\033[0m")
	}
}



// appendHttpRequest appends HTTP request details to the log message
func appendHttpRequest(result *strings.Builder, c Container) {
	if c.HttpRequest != nil {
		result.WriteString(" " + c.HttpRequest.RemoteAddr + " " + c.HttpRequest.Method + " " + c.HttpRequest.URL.String())
	}
}



// appendProcessingTime appends the processing time with color formatting to the log message
func appendProcessingTime(result *strings.Builder, c Container) {
	processingTimeMs := c.ProcessingTime.Microseconds() / 1000
	formattedTime := strconv.FormatInt(processingTimeMs, 10) + " ms"

	if processingTimeMs != 0 {
		result.WriteString(" [" + formattedTime + "]")
	}
}



// logWithJson logs the message with additional JSON data
func logWithJson(message string, c Container) {
	wJsonBytes, err := json.MarshalIndent(c.ProcessedData, "", "  ")
	if err != nil {
		log.Println("Error marshaling to JSON:", err)
	}
	wJsonData := string(wJsonBytes)

	rJsonData := ""
	if c.HttpRequest != nil {
		rJsonBytes, err := json.MarshalIndent(c.HttpRequest.Body, "", "  ")
		if err != nil {
			log.Println("Error marshaling to JSON:", err)
		}
		rJsonData = string(rJsonBytes)
	}

	message += fmt.Sprintf(" BODY JSON: %s WRITE/PROCESSED JSON: %s", rJsonData, wJsonData)
}



// writeLog writes the log entry
func writeLog(message string, c *Container) {
	// Format the log file name as YYYY_MM_DD.log based on the log event timestamp
	logFileName := c.Timestamp.Format("2006_01_02") + ".log"

	// Open the log file in append mode, create if it doesn't exist
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	// Write the log message to the file
	_, err = fmt.Fprintln(file, message)
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
}
