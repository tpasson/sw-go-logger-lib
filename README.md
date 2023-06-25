# Logger Package
The `logger` package provides a logging utility that allows you to log messages with different severity levels. It includes features for structured logging and supports various log severities, such as NONE, PRODUCTION, and DEBUG.

## Installation
To use the `logger` package, you need to have Go installed and set up. Then, you can install the package by running the following command:

```bash
go get github.com/<username>/<repository>/logger
```

## Usage
Import the logger package in your Go code:

```go
import "github.com/<username>/<repository>/logger"
```

## Creating a Logger
To create a new logger instance, use the `NewLogger` function:

```go
severity := "PRODUCTION" // Specify the log severity
logger, err := logger.NewLogger(severity)
if err != nil {
    // Handle the error
}
```

Valid log severities are `NONE`, `PRODUCTION`, and `DEBUG`. An error will be returned if an invalid severity is provided.

## Logging Messages
To log a message, use the `Entry` method of the logger instance:

```go
container := logger.Container{
    PreText:        "Some pre-text",
    Id:             "123",
    Source:         "example.go",
    Info:           "Some information",
    Data:           "Some data",
    Error:          "An error occurred",
    ProcessingTime: time.Duration(500 * time.Millisecond),
    HttpRequest:    req,
    ProcessedData:  data,
}

logger.Entry(container)
```

The `Container` struct contains the necessary information for the log entry, including pre-text, ID, source, information, data, error message, processing time, HTTP request (optional), and processed data. You can customize the values based on your specific use case.

The log message will be printed according to the severity level set during the logger creation. If the severity is set to `NONE`, the message will be ignored.

## Log Output
The log output will be printed to the console or standard output, depending on the severity level:

* If the severity is set to PRODUCTION, the log message will be printed using the log.Println function.
* If the severity is set to DEBUG, the log message will be printed using the log.Println function, and additional details such as the processed data and request body will be logged in JSON format.

## Contributing
Contributions to the logger package are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.