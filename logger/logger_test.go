package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerOutput(t *testing.T) {
	// Delete previous generated log files
	err := deleteLogFiles()
	if err != nil {
		fmt.Println("Error deleting log files:", err)
	} else {
		fmt.Println("Log files deleted successfully")
	}

	candidateOne(t)
	candidateTwo(t)
	candidateThree(t)
	candidateFour(t)
	candidateFive(t)
}

func candidateOne(t *testing.T) {
	// Create a reference timestamp
	ts := time.Now()

	// Create a new logger with desired format
	logger, err := NewLogger(
		[]LogFormat{
			FORMAT_TIMESTAMP,
			FORMAT_STATUS,
			FORMAT_PRE_TEXT,
			FORMAT_HTTP_REQUEST,
			FORMAT_ID,
			FORMAT_SOURCE,
			FORMAT_INFO,
			FORMAT_DATA,
			FORMAT_ERROR,
			FORMAT_PROCESSING_TIME,
			FORMAT_PROCESSED_DATA,
		}, Options{
			OutputStdout:     true,
			OutputFile:       true,
			OutputFolderPath: "",
		}, Container{
			Status:    STATUS_INFO,
			Info:      "System Logger succesfully started! Awaiting logger tasks...",
			Timestamp: ts,
		})
	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	// Create a mock HTTP request for testing
	request, _ := http.NewRequest("GET", "https://example.com", nil)
	// Set the remote address
	request.RemoteAddr = "192.168.0.1:12345"

	data := map[string]interface{}{
		"name":     "John Doe",
		"age":      30,
		"isActive": true,
		"tags":     []string{"go", "programming", "dummy"},
	}

	// Create a log entry container
	container := Container{
		Timestamp:      ts,
		Status:         STATUS_INFO,
		PreText:        "SERVER1",
		HttpRequest:    request,
		Id:             "5f322ac4ba",
		Source:         "handler/user",
		Info:           "This is an information message",
		Data:           "233",
		Error:          "something went wrong",
		ProcessingTime: 1 * time.Millisecond,
		ProcessedData:  data,
	}

	// Redirect STDOUT to capture the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the Entry method to log the container
	logger.Entry(container)

	duration := 20 * time.Millisecond
	time.Sleep(duration)

	// Reset STDOUT
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output from the pipe
	var capturedOutput strings.Builder
	io.Copy(&capturedOutput, r)

	// Verify the captured output
	res1 := ts.Format(time.RFC3339) + " INFO System Logger succesfully started! Awaiting logger tasks... >Processed Data:\nnull\n"
	res2 := ts.Format(time.RFC3339) + " INFO SERVER1 192.168.0.1:12345 GET https://example.com 5f322ac4ba handler/user This is an information message 233 something went wrong [1 ms]"
	res3 := " >Processed Data:\n{\n  \"age\": 30,\n  \"isActive\": true,\n  \"name\": \"John Doe\",\n  \"tags\": [\n    \"go\",\n    \"programming\",\n    \"dummy\"\n  ]\n}\n"
	expected := res1 + res2 + res3
	actual := capturedOutput.String()

	if string(actual) != string(expected) {
		t.Errorf("Unexpected result.\nExpected:\n%#v\nGot:\n%#v", expected, actual)
	}
}

func candidateTwo(t *testing.T) {
	// Create a reference timestamp
	ts := time.Now()

	// Create a new logger with desired format
	logger, err := NewLogger(
		[]LogFormat{
			FORMAT_STATUS,
			FORMAT_PRE_TEXT,
			FORMAT_HTTP_REQUEST,
			FORMAT_ID,
			FORMAT_SOURCE,
			FORMAT_INFO,
			FORMAT_DATA,
			FORMAT_ERROR,
			FORMAT_PROCESSING_TIME,
			FORMAT_PROCESSED_DATA,
		}, Options{
			OutputStdout:     true,
			OutputFile:       true,
			OutputFolderPath: "",
		}, Container{
			Status:    STATUS_INFO,
			Info:      "System Logger succesfully started! Awaiting logger tasks...",
			Timestamp: ts,
		})
	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	// Create a log entry container
	container := Container{
		Status:         STATUS_INFO,
		PreText:        "SERVER1",
		Id:             "5f322ac4ba",
		Source:         "handler/user",
		Info:           "This is an information message",
		Data:           "233",
		Error:          "something went wrong",
		ProcessingTime: 1 * time.Millisecond,
	}

	// Redirect STDOUT to capture the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the Entry method to log the container
	logger.Entry(container)

	duration := 20 * time.Millisecond
	time.Sleep(duration)

	// Reset STDOUT
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output from the pipe
	var capturedOutput strings.Builder
	io.Copy(&capturedOutput, r)

	// Verify the captured output
	expected := "INFO System Logger succesfully started! Awaiting logger tasks... >Processed Data:\nnull\nINFO SERVER1 5f322ac4ba handler/user This is an information message 233 something went wrong [1 ms] >Processed Data:\nnull\n"
	actual := capturedOutput.String()

	if string(actual) != string(expected) {
		t.Errorf("Unexpected result.\nExpected:\n%#v\nGot:\n%#v", expected, actual)
	}
}

func candidateThree(t *testing.T) {
	// Create a reference timestamp
	ts := time.Now()

	// Create a new logger with desired format
	logger, err := NewLogger([]LogFormat{}, Options{
		OutputStdout:     true,
		OutputFile:       true,
		OutputFolderPath: "",
	}, Container{
		Status:    STATUS_INFO,
		Info:      "System Logger succesfully started! Awaiting logger tasks...",
		Timestamp: ts,
	})
	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	// Create a log entry container
	container := Container{
		Status:         STATUS_INFO,
		PreText:        "SERVER1",
		Id:             "5f322ac4ba",
		Source:         "handler/user",
		Info:           "This is an information message",
		Data:           "233",
		Error:          "something went wrong",
		ProcessingTime: 1 * time.Millisecond,
	}

	// Redirect STDOUT to capture the output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the Entry method to log the container
	logger.Entry(container)

	duration := 20 * time.Millisecond
	time.Sleep(duration)

	// Reset STDOUT
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output from the pipe
	var capturedOutput strings.Builder
	io.Copy(&capturedOutput, r)

	// Verify the captured output
	expected := ""
	actual := capturedOutput.String()

	if string(actual) != string(expected) {
		t.Errorf("Unexpected result.\nExpected:\n%#v\nGot:\n%#v", expected, actual)
	}
}

func candidateFour(t *testing.T) {
	// Create a reference timestamp
	ts := time.Now()

	// Create a new logger with desired format
	logger, err := NewLogger(
		[]LogFormat{
			FORMAT_STATUS,
			FORMAT_ID,
		}, Options{
			OutputStdout:     true,
			OutputFile:       true,
			OutputFolderPath: "",
		}, Container{
			Status:    STATUS_INFO,
			Info:      "System Logger succesfully started! Awaiting logger tasks...",
			Timestamp: ts,
		})
	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	// Create a log entry container
	containerInfo := Container{
		Status:         STATUS_INFO,
		Id:             "21BTC",
		ProcessingTime: 1 * time.Millisecond,
	}

	containerError := Container{
		Status:         STATUS_ERROR,
		Id:             "21BTC",
		ProcessingTime: 1 * time.Millisecond,
	}

	containerFatal := Container{
		Status:         STATUS_FATAL,
		Id:             "21BTC",
		ProcessingTime: 1 * time.Millisecond,
	}

	containerTrace := Container{
		Status:         STATUS_TRACE,
		Id:             "21BTC",
		ProcessingTime: 1 * time.Millisecond,
	}

	containerWarn := Container{
		Status:         STATUS_WARN,
		Id:             "21BTC",
		ProcessingTime: 1 * time.Millisecond,
	}

	// Call the Entry method to log the container
	logger.Entry(containerInfo)
	logger.Entry(containerInfo)
	logger.Entry(containerInfo)
	logger.Entry(containerInfo)
	logger.Entry(containerInfo)
	logger.Entry(containerError)
	logger.Entry(containerError)
	logger.Entry(containerError)
	logger.Entry(containerError)
	logger.Entry(containerFatal)
	logger.Entry(containerFatal)
	logger.Entry(containerFatal)
	logger.Entry(containerTrace)
	logger.Entry(containerTrace)
	logger.Entry(containerWarn)

	duration := 100 * time.Millisecond
	time.Sleep(duration)

	// Verify the captured output
	expected := "Log Level Counters: [INFO: 6] [WARN: 1] [TRACE: 2] [ERROR: 4] [FATAL: 3]"
	actual := logger.GetLogStatusCounters()

	if string(actual) != string(expected) {
		t.Errorf("Unexpected result.\nExpected:\n%#v\nGot:\n%#v", expected, actual)
	}
}

func candidateFive(t *testing.T) {
	// Create a reference timestamp
	ts := time.Now()

	// Create a new logger with desired format
	_, err := NewLogger(
		[]LogFormat{
			FORMAT_STATUS,
			FORMAT_ID,
		}, Options{
			OutputStdout:     true,
			OutputFile:       true,
			OutputFolderPath: "folder/not/existing/",
		}, Container{
			Status:    STATUS_INFO,
			Info:      "System Logger succesfully started! Awaiting logger tasks...",
			Timestamp: ts,
		})
	if err == nil {
		t.Errorf("Unexpected result: Code should throw an error here")
	}

	// Verify the captured output
	expected := "open folder/not/existing/testfile.tmp: no such file or directory"
	actual := err.Error()

	if string(actual) != string(expected) {
		t.Errorf("Unexpected result.\nExpected:\n%#v\nGot:\n%#v", expected, actual)
	}
}

func deleteLogFiles() error {
	dir, err := os.Getwd() // Get current working directory
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.log")) // Get list of all log files
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(file) // Delete each log file
		if err != nil {
			return err
		}
	}

	return nil
}
