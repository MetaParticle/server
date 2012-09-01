package logger

type loglevel int

const (
	PANIC loglevel = 0
	ERROR
	NORMAL
)

const (
	SEVERE loglevel = 1+iota
	MAJOR
	WARNING
	MINOR
	IMPORTANT // Should be logged during normal operation
	INFO
)
