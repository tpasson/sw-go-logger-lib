# Logger Package
The `logger` package provides a logging utility that allows you to log messages with different severity levels. It includes features for structured logging.
## Usage

1. In your project, open a terminal or command prompt and navigate to the root directory of your Go project.

2. Use the `go get` command followed by the import path to download the package and its dependencies. For example, if the import path is `github.com/passon-engineering/sw-go-logger-lib`, run the following command:
   ```
   go get github.com/passon-engineering/sw-go-logger-lib
   ```

3. After running `go get`, Go will download the package and its dependencies and store them in your local Go workspace.

4. In your Go code, you can import the package using the import path you identified earlier. For example:
   ```go
   import "github.com/passon-engineering/sw-go-logger-lib"
   ```

5. You can now use the functions, types, and other elements provided by the external package in your own code.

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