package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// New takes a set of options and returns a new ChatServer instance
func New(opts ...Option) (*ChatServer, error) {
	var c ChatServer

	for _, opt := range opts {
		if err := opt.Apply(&c); err != nil {
			return nil, err
		}
	}

	if c.upgrader.ReadBufferSize == 0 || c.upgrader.WriteBufferSize == 0 {
		c.upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
	}

	if c.room == nil {
		c.room = &Room{
			mux:     sync.Mutex{},
			clients: make(map[*websocket.Conn]string),
		}
	}

	return &c, nil
}

// ChatServer is the Chat ChatServer
type ChatServer struct {
	upgrader websocket.Upgrader
	room     *Room
}

func (s *ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	last := parts[len(parts)-1]
	if last == "websocket" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "you must provide a username"}`))
		return
	}

	username := last

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "your browser doesn't seem to support websockets... fool"}`))
		return
	}

	if err := s.room.Add(conn, username); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
	}

	welcome := NewPayload("welcome", "Chat Bot", "")

	if err := conn.WriteJSON(welcome); err != nil {
		// log.Printf("error writing json to connection: %s", err.Error())
		return
	}

	if err := s.room.Broadcast(conn, fmt.Sprintf("%s has joined the room!", username)); err != nil {
		return
	}

	for {
		var txt struct {
			Message string `json:"message"`
		}

		if err := conn.ReadJSON(&txt); err != nil {
			log.Printf("attempting to read from conn")
			return
		}

		s.room.Broadcast(conn, txt.Message)
	}
}

// Room holds all of the connections
type Room struct {
	mux     sync.Mutex
	clients map[*websocket.Conn]string // string is the username
}

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

// Broadcast takes a message and sends it to all in the room
func (r *Room) Broadcast(conn *websocket.Conn, msg string) error {
	username := r.clients[conn]

	log.Printf("Got message from %s \n Message: \n\t %s", username, msg)

	message := NewPayload("chat-message", username, msg)

	for conn, _ := range r.clients {
		if err := conn.WriteJSON(message); err != nil {
			return err
		}
	}

	return nil
}

// Add takes a username and connection and adds it to the list of clients
func (r *Room) Add(conn *websocket.Conn, username string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.clients[conn] = username

	return nil
}

// Remove removes the connection from the room
func (r *Room) Remove(conn *websocket.Conn) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.clients, conn)

	return nil
}

// Option is anything that has an Apply(*ChatServer) error
type Option interface {
	Apply(*ChatServer) error
}

// WithUpgrader takes an upgrader and sets it on the ChatServer
func WithUpgrader(upgrader websocket.Upgrader) Option {
	return &withUpgrader{upgrader}
}

type withUpgrader struct {
	upgrader websocket.Upgrader
}

func (opt *withUpgrader) Apply(c *ChatServer) error {
	return nil
}
