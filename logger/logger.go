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



/* 
NewLogger creates a new Logger instance with the specified severity.
Valid log severities are: NONE, PRODUCTION, DEBUG.

Parameters:
	- severity: string - the desired severity level (NONE, PRODUCTION, DEBUG)

Returns:
	- *Logger: the created Logger instance
	- error: an error if the severity level is invalid 
*/
func NewLogger(severity string) (*Logger, error) {
	if !isValidSeverity(severity) {
		return nil, fmt.Errorf("invalid log severity: %s. valid log severities: %s", severity, getValidSeverities())
	}

	logSeverity := severityMap[severity]

	logger := &Logger{
		Severity: logSeverity,
		LogChan:  make(chan Container),
	}

	go logger.processLogs()

	return logger, nil
}



/* 
Entry logs a message based on the severity level and the provided container.

If the logger's severity level is set to NONE, the log entry will be skipped.

If the timestamp of the provided container is zero, it will be set to the current
timestamp using the generateTimestamp function.

The log entry is then sent to the logger's LogChan channel for further processing.

Parameters:
	- c: Container - the log entry container containing the log message and metadata 
*/
func (l *Logger) Entry(c Container) {
	if l.Severity == NONE {
		return
	}

	if c.Timestamp.IsZero() {
		c.Timestamp = generateTimestamp()
	}

	l.LogChan <- c
}



/* 
isValidSeverity checks if the given severity string is valid.

Parameters:
	- severity: string - the severity string to check

Returns:
	- bool: true if the severity string is valid, false otherwise 
*/
func isValidSeverity(severity string) bool {
	_, ok := severityMap[severity]
	return ok
}



/*
getValidSeverities returns a string containing all valid severities.

Returns:
	- string: a comma-separated string of valid severities
*/
func getValidSeverities() string {
	validSeverities := make([]string, 0, len(severityMap))
	for severity := range severityMap {
		validSeverities = append(validSeverities, severity)
	}
	return strings.Join(validSeverities, ", ")
}



/* 
generateTimestamp creates the current timestamp.

Returns:
	- time.Time: the current timestamp 
*/
func generateTimestamp() time.Time {
	return time.Now()
}



/* 
formatTimestamp formats the given timestamp and returns it as a string.

Parameters:
	- timestamp: time.Time - the timestamp to format

Returns:
	- string: the formatted timestamp 
*/
func formatTimestamp(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
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



/* 
appendFields appends non-empty fields from the container to the log message.

Parameters:
	- result: *strings.Builder - the string builder to append the fields to
	- c: Container - the log entry container 
*/
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



/* 
appendError appends the error message with color formatting to the log message.

Parameters:
	- result: *strings.Builder - the string builder to append the error message to
	- c: Container - the log entry container 
*/
func appendError(result *strings.Builder, c Container) {
	if c.Error != "" {
		result.WriteString(" \033[31mError: " + c.Error + "\033[0m")
	}
}



/* 
appendHttpRequest appends HTTP request details to the log message.

Parameters:
	- result: *strings.Builder - the string builder to append the HTTP request details to
	- c: Container - the log entry container 
*/
func appendHttpRequest(result *strings.Builder, c Container) {
	if c.HttpRequest != nil {
		result.WriteString(" " + c.HttpRequest.RemoteAddr + " " + c.HttpRequest.Method + " " + c.HttpRequest.URL.String())
	}
}



/* 
appendProcessingTime appends the processing time with color formatting to the log message.

Parameters:
	- result: *strings.Builder - the string builder to append the processing time to
	- c: Container - the log entry container 
*/
func appendProcessingTime(result *strings.Builder, c Container) {
	processingTimeMs := c.ProcessingTime.Microseconds() / 1000
	formattedTime := strconv.FormatInt(processingTimeMs, 10) + " ms"

	if processingTimeMs != 0 {
		result.WriteString(" [" + formattedTime + "]")
	}
}



/* 
logWithJson logs the message with additional JSON data.

Parameters:
	- message: string - the log message to append the JSON data to
	- c: Container - the log entry container 
*/
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



/* 
writeLog writes the log entry.

Parameters:
	- message: string - the log message to write
	- c: *Container - the log entry container 
*/
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
