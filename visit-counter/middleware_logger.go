package main

import (
	"net/http"
	"time"
)

func NewMddlLogger(l *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("Receiving Request")
			start := time.Now()
			// l.Info(r)

			next.ServeHTTP(w, r)

			l.Info("Request Finished")
			l.Info("Took:", time.Since(start))
		})
	}
}
