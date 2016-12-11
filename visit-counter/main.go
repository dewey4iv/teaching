package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	if err := ReadVisits(visitFile, counterService); err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		panic(err)
	}

	l.Info("Setting up HTTP Server")
	mux := http.NewServeMux()
	loggerMiddleware := NewMddlLogger(l)
	mux.Handle("/visit", loggerMiddleware(counterService))

	go func(mux *http.ServeMux) {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			panic(err)
		}
	}(mux)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case _ = <-sigCh:
		l.Info("Received Shutdown Signal")
		l.Info("Writing Visits to Disk")
		if err := SaveVisits(visitFile, counterService); err != nil {
			panic(err)
		}
		l.Info("Shutting Down")
	}
}
