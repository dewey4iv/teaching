package main

import (
	"encoding/json"
	"time"
)

func NewPayload(t string, from string, message string) Payload {
	return Payload{
		Type:      t,
		From:      from,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

type Payload struct {
	Type      string      `json:"type"`
	From      string      `json:"from"`
	Message   interface{} `json:"message"`
	CreatedAt time.Time   `json:"created_at"`
}

func (p Payload) MarshalJSON() ([]byte, error) {
	// Mon Jan 2 15:04:05 MST 2006
	createdAt := p.CreatedAt.Format(time.Kitchen)

	var msg = struct {
		Type      string      `json:"type"`
		From      string      `json:"from"`
		Message   interface{} `json:"message"`
		CreatedAt string      `json:"created_at"`
	}{p.Type, p.From, p.Message, createdAt}

	return json.Marshal(msg)
}
