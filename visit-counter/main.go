package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	visitFile = "./visits"
)

func main() {
	// create new logger
	l, err := NewLogger("./logs")
	if err != nil {
		panic(err)
	}

	counterService := &CounterService{0, l}

	l.Info("Attempting to set Visits")
	if err := ReadVisits(visitFile, counterService); err != nil && os.IsNotExist(err) {
		panic(err)
	}

	l.Info("Setting up HTTP Server")
	mux := http.NewServeMux()
	loggerMiddleware := NewMddlLogger(l)
	mux.Handle("/visit", loggerMiddleware(counterService))

	httpChan := make(chan error)

	go func(httpChan chan error, mux *http.ServeMux, l *Logger) {
		l.Info("Starting HTTP Server")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			httpChan <- err
		}
	}(httpChan, mux, l)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpChan:
		l.Err("Error w/ HTTP Server")
		l.Err(err)
	case _ = <-sigCh:
		l.Info("Received Shutdown Signal")
		l.Info("Writing Visits to Disk")
		if err := SaveVisits(visitFile, counterService); err != nil {
			panic(err)
		}
		l.Info("Shutting Down")
	}
}
