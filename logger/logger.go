package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	Format []LogFormat
	LogChan  chan Container
	StatusCounters map[LogStatus]int
}

type Container struct {
	Status				 LogStatus
	PreText        string
	Id             string
	Source         string
	Info           string
	Data           string
	Error          string
	ProcessingTime time.Duration
	Timestamp      time.Time
	HttpRequest    *http.Request
	ProcessedData  any
}

// Creates a new Logger instance with the specified ontent.
//
// Parameters:
//	- format: string - the desired format
//
// Returns:
//	- *Logger: the created Logger instance
func NewLogger(format []LogFormat) (*Logger) {
	logger := &Logger{
		Format: format,
		LogChan:  make(chan Container),
		// Initialize the LevelCounters map
		StatusCounters: make(map[LogStatus]int),
	}

	go logger.processLogs()

	return logger
}



// Logs a message based on the severity level and the provided container.
//
// If the logger's severity level is set to NONE, the log entry will be skipped.
//
// If the timestamp of the provided container is zero, it will be set to the current
// timestamp using the generateTimestamp function.
//
// The log entry is then sent to the logger's LogChan channel for further processing.
//
// Parameters:
//	- c: Container - the log entry container containing the log message and metadata 
func (l *Logger) Entry(c Container) {
	// Check for element - if empty: logger disabled
	if len(l.Format) == 0 {
		return
	}

	if c.Timestamp.IsZero() {
		c.Timestamp = generateTimestamp()
	}

	l.LogChan <- c
}



// Creates the current timestamp.
//
// Returns:
//	- time.Time: the current timestamp 
func generateTimestamp() time.Time {
	return time.Now()
}



// Formats the given timestamp and returns it as a string.
//
// Parameters:
//	- timestamp: time.Time - the timestamp to format
//
// Returns:
//	- string: the formatted timestamp 
func formatTimestamp(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
}



// Processes logs from the log channel and writes them to the log file.
//
// It is a method of the Logger type and is executed as a goroutine. It continuously reads log entries
// from the log channel (`l.LogChan`) and processes each log entry by formatting it based on the configured
// log format items. The formatted log message is then written to the log file and also printed to STDOUT.
//
// This method uses various helper functions to format different log components based on the configured format items.
// It also trims any trailing spaces from the formatted log message before writing it to the log file.
func (l *Logger) processLogs() {
	for c := range l.LogChan {


		// Create buffer
		var result strings.Builder

		for _, formatItem := range l.Format {
			switch formatItem {
			case STATUS:
				if str := logStatustoString[c.Status]; str != "" {
					// Increment the log level counter
					incrementLogStatusCounter(l, c.Status)
					result.WriteString(str + " ")
				}
			case PRE_TEXT:
				if c.PreText != "" {
					result.WriteString(c.PreText + " ")
				}
			case ID:
				if c.Id != "" {
					result.WriteString(c.Id + " ")
				}
			case SOURCE:
				if c.Source != "" {
					result.WriteString(c.Source + " ")
				}
			case INFO:
				if c.Info != "" {
					result.WriteString(c.Info + " ")
				}
			case DATA:
				if c.Data != "" {
					result.WriteString(c.Data + " ")
				}
			case ERROR:
				if c.Error != "" {
					result.WriteString(c.Error + " ")
				}
			case PROCESSING_TIME:
				if str := getProcessingTime(c.ProcessingTime); str != "" {
					result.WriteString(str + " ")
				}
			case TIMESTAMP:
				if str := formatTimestamp(c.Timestamp); str != "" {
					result.WriteString(str + " ")
				}
			case HTTP_REQUEST:
				if str := getHttpRequest(c.HttpRequest); str != "" {
					result.WriteString(str + " ")
				}
			case PROCESSED_DATA:
				if str := getProcessedData(c.ProcessedData); str != "" {
					result.WriteString(str + " ")
				}
			}
		}

		trimmedResult := strings.TrimRight(result.String(), " ")

		writeLog(trimmedResult, &c)
	}
}



// Returns a formatted string representation of an HTTP request.
//
// It takes an *http.Request object as input and returns a string containing the remote address,
// HTTP method, and URL of the request. If the provided HTTP request is nil, an empty string is returned.
//
// Parameters:
//    - httpRequest: *http.Request - the HTTP request object to format
//
// Returns:
//    - string: the formatted HTTP request
//
// Example:
//    request := &http.Request{
//        Method: "GET",
//        URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
//    }
//    result := getHttpRequest(request)
//    // result will be "192.168.0.1:12345 GET https://example.com/"
func getHttpRequest(httpRequest *http.Request) string {
	if httpRequest != nil {
		return(httpRequest.RemoteAddr + " " + httpRequest.Method + " " + httpRequest.URL.String())
	}
	return ""
}



// Returns the processing time as a formatted string.
//
// It takes a time.Duration value representing the processing time as input. The function
// converts the processing time to milliseconds and formats it as "[X ms]", where X is the
// number of milliseconds. If the processing time is zero, an empty string is returned.
//
// Parameters:
//    - processingTime: time.Duration - the processing time to format
//
// Returns:
//    - string: the formatted processing time
//
// Example:
//    processingTime := time.Duration(500 * time.Millisecond)
//    result := getProcessingTime(processingTime)
//    // result will be "[500 ms]"
func getProcessingTime(processingTime time.Duration) string {
	processingTimeMs := processingTime.Microseconds() / 1000
	formattedTime := strconv.FormatInt(processingTimeMs, 10) + " ms"

	if processingTimeMs != 0 {
		result := "[" + formattedTime + "]"
		return result
	}
	return ""
}



// Serializes the provided data to JSON format.
//
// It takes any value as the input data and marshals it into JSON format using
// the json.MarshalIndent function. The data is indented with two spaces per level.
// If an error occurs during the marshaling process, the error message is returned.
// Otherwise, the marshaled data is returned as a string.
//
// Parameters:
//    - processedData: any - the data to be processed
//
// Returns:
//    - string: the processed data in JSON format, or an error message
//
// Example:
//    result := getProcessedData(data)
//
// Output:
//    {"name": "John Doe", "age": 30}
func getProcessedData(processedData any) string {
	wJsonBytes, err := json.MarshalIndent(processedData, "", "  ")	
	if err != nil {
		return (err.Error())
	}
	wJsonData := string(wJsonBytes)

	return wJsonData
}



// Writes the log message to a log file and also prints it to STDOUT.
//
// It formats the log file name as "YYYY_MM_DD.log" based on the log event timestamp.
// The log file is opened in append mode and created if it doesn't exist.
// The log message is written to the file and also printed to STDOUT.
//
// Parameters:
//    - message: string - the log message to write
//    - c: *Container - the log entry container
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

	// Write the log message to STDOUT
	fmt.Println(message)
}