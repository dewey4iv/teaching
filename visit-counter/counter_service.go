package main

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type CounterService struct {
	Visits int64 `json:"visits"`
	logger *Logger
}

func (cs *CounterService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	atomic.AddInt64(&cs.Visits, 1)

	respBytes, err := json.Marshal(VisitResponse{cs.Visits})
	if err != nil {
		cs.logger.Err("Err marshaling response:")
		cs.logger.Err(err)
	}

	if _, err := w.Write(respBytes); err != nil {
		cs.logger.Err("Err writing to response:")
		cs.logger.Err(err)
	}

}

func (cs *CounterService) GetVisits() int64 {
	return cs.Visits
}

func (cs *CounterService) SetVisits(visits int64) {
	cs.Visits = visits
}

type VisitResponse struct {
	Visitors int64 `json:"visitors"`
}
