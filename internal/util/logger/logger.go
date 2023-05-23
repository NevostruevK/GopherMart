package logger

import (
	"log"
	"os"
)

var logWriter = os.Stdout

func NewLogger(name string, flags int) *log.Logger {
	return log.New(logWriter, name, flags)
}
