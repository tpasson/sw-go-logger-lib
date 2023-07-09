package logger

import (
	"fmt"
	"sort"
	"strings"
)

// The status which will be displayed in the message e.g. [WARN]
type LogStatus int

const (
	STATUS_INFO LogStatus = iota
	STATUS_WARN
	STATUS_TRACE
	STATUS_ERROR
	STATUS_FATAL
)

var logStatustoString = map[LogStatus]string{
	STATUS_INFO:  "INFO",
	STATUS_WARN:  "WARN",
	STATUS_TRACE: "TRACE",
	STATUS_ERROR: "ERROR",
	STATUS_FATAL: "FATAL",
}

// Increments the log level counter for the given log status.
//
// It is a function that takes a Logger instance and a Container pointer as arguments. The function increments
// the log level counter for the log status specified in the Container. The log level counters are maintained
// within the Logger instance.
//
// Example:
//
//	logger := logger.NewLogger()
//	container := &Container{
//	    Status: STATUS_INFO,
//	}
//	incrementLogStatusCounter(logger, container.Status)
//	fmt.Println(logger.GetLogStatusCounters())
//
// Output:
//
//	Log Level Counters:
//	  INFO: 1
func incrementLogStatusCounter(l *Logger, ls LogStatus) {
	l.StatusCounters[ls]++
}

// Returns a formatted string representing the log level counters.
//
// It is a method of the Logger type and is used to retrieve the current count of log entries for each log level.
// The log level counters are represented as a formatted string containing the log level names and their corresponding
// count values.
//
// The method iterates over the log level counters stored in the `l.StatusCounters` map. It sorts the keys (log levels)
// in ascending order and retrieves the count value for each log level. The log level names and count values are then
// formatted and appended to a strings.Builder. The resulting formatted string represents the log level counters.
//
// Example:
//
//	logger := logger.NewLogger()
//	// Log some entries
//	counters := logger.GetLogStatusCounters()
//	fmt.Println(counters)
//	// Output example: Log Level Counters: [DEBUG: 2] [INFO: 5] [WARNING: 3] [ERROR: 1]
func (l *Logger) GetLogStatusCounters() string {
	var builder strings.Builder
	builder.WriteString("Log Level Counters:")

	// Sort the keys of the log level counters
	keys := make([]int, 0, len(l.StatusCounters))
	for status := range l.StatusCounters {
		keys = append(keys, int(status))
	}
	sort.Ints(keys)

	// Iterate over the sorted keys and retrieve the counter values
	for _, status := range keys {
		count := l.StatusCounters[LogStatus(status)]
		builder.WriteString(fmt.Sprintf(" [%s: %d]", logStatustoString[LogStatus(status)], count))
	}

	return builder.String()
}
