package logger

// The format defines how much information is being logged and in which order. Has to be defined while initalizing the logger
// Possible key fields for the format are:
/*
	STATUS
	PRE_TEXT
	ID
	SOURCE
	INFO
	DATA
	ERROR
	PROCESSING_TIME
	TIMESTAMP
	HTTP_REQUEST
	PROCESSED_DATA
*/
type LogFormat int

const (
	FORMAT_STATUS LogFormat = iota
	FORMAT_PRE_TEXT
	FORMAT_ID
	FORMAT_SOURCE
	FORMAT_INFO
	FORMAT_DATA
	FORMAT_ERROR
	FORMAT_PROCESSING_TIME
	FORMAT_TIMESTAMP
	FORMAT_HTTP_REQUEST
	FORMAT_PROCESSED_DATA
)
