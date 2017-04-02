package payload

import (
	"encoding/json"
	"time"
)

// NewServerMsg returns a payload
func NewServerMsg(t string, message string) Payload {
	return Payload{
		Type:      t,
		From:      "Chat Bot",
		Message:   message,
		CreatedAt: time.Now(),
	}
}

// NewUserMsg returns a payload
func NewUserMsg(username string, message string) Payload {
	return Payload{
		Type:      "chat-message",
		From:      username,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

// New returns a raw payload
func New(t string, from string, message string) Payload {
	return Payload{
		Type:      t,
		From:      from,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

// Payload defines the basic message passed to and from the client
type Payload struct {
	Type      string      `json:"type"`
	From      string      `json:"from"`
	Message   interface{} `json:"message"`
	CreatedAt time.Time   `json:"created_at"`
}

// MarshalJSON implments json.Marshaller
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
