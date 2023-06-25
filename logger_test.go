package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestGeneral(t *testing.T) {
	// Delete previous generated log files
	err := deleteLogFiles()
	if err != nil {
		fmt.Println("Error deleting log files:", err)
	} else {
		fmt.Println("Log files deleted successfully")
	}

	// Create a logger instance
	log, err := NewLogger("DEBUG")
	if err != nil {
		fmt.Println("Failed to create logger:", err)
		return
	}

	// WaitGroup to synchronize the completion of all goroutines
	var wg sync.WaitGroup
	numWorkers := 3
	numIterations := 200

	wg.Add(numWorkers)

	// Start timer
	startTime := time.Now()

	// Start worker goroutines
	go workerA(numIterations, log, &wg)
	go workerB(numIterations, log, &wg)
	go workerC(numIterations, log, &wg)


	// Wait for all worker goroutines to complete
	wg.Wait()

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Println("All workers completed.")
	fmt.Printf("Elapsed Time: %s\n", elapsedTime)
}



func workerA(numIterations int, log *Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < numIterations; i++ {
		start := time.Now()

		duration := 1 * time.Millisecond
		time.Sleep(duration)

		// Create a log entry
		container := Container{
			PreText:        "Worker A",
			Id:             fmt.Sprintf("ID%d", i+1),
			Source:         fmt.Sprintf("Source%d", i+1),
			Info:           fmt.Sprintf("Info%d", i+1),
			Data:           fmt.Sprintf("Data%d", i+1),
			Error:          "",
			ProcessingTime: time.Since(start),
			HttpRequest:    nil,
			ProcessedData:  nil,
		}

		// Pass the log entry to the logger
		log.Entry(container)
	}
}

func workerB(numIterations int, log *Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < numIterations; i++ {
		// Create a log entry
		container := Container{
			PreText:        "Worker B",
			Id:             fmt.Sprintf("ID%d", i+1),
			Source:         fmt.Sprintf("Source%d", i+1),
			Info:           fmt.Sprintf("Info%d", i+1),
			Data:           fmt.Sprintf("Data%d", i+1),
			Error:          "",
			ProcessingTime: 0,
			Timestamp:      generateTimestamp(),
			HttpRequest:    nil,
			ProcessedData:  nil,
		}

		// Pass the log entry to the logger
		log.Entry(container)
	}
}

func workerC(numIterations int, log *Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < numIterations; i++ {
		// Create a log entry
		container := Container{
			PreText:        "Worker C",
			Id:             fmt.Sprintf("ID%d", i+1),
			Source:         fmt.Sprintf("Source%d", i+1),
			Info:           fmt.Sprintf("Info%d", i+1),
			Data:           fmt.Sprintf("Data%d", i+1),
			Error:          "",
			ProcessingTime: 0,
			HttpRequest:    nil,
			ProcessedData:  nil,
		}

		// Pass the log entry to the logger
		log.Entry(container)
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