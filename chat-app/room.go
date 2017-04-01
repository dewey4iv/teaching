package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Room holds all of the connections
type Room struct {
	mux     sync.Mutex
	clients map[*websocket.Conn]string // string is the username
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
