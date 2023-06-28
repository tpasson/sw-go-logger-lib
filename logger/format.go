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
	STATUS LogFormat = iota
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
)