package main

import (
	"io"
	"log"
	"os"
)

// NewLogger takes the path that we want to start logging and creates a logger
// with a file in that location
func NewLogger(path string) (*Logger, error) {
	logFile, err := os.OpenFile(path+"/access.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	errFile, err := os.OpenFile(path+"/error.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var logger Logger

	logOutput := io.MultiWriter(os.Stdout, logFile)
	logger.log = log.New(logOutput, "[ INFO ] \t", log.LUTC|log.LstdFlags)

	errOutput := io.MultiWriter(os.Stderr, errFile)
	logger.err = log.New(errOutput, "[ ERROR ] \t", log.LUTC|log.LstdFlags|log.Lshortfile)

	return &logger, nil
}

type Logger struct {
	log *log.Logger
	err *log.Logger
}

func (l *Logger) Info(in ...interface{}) {
	l.log.Println(in...)
}

func (l *Logger) Err(in ...interface{}) {
	l.err.Println(in...)
}
