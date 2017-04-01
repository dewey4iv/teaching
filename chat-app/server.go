package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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
