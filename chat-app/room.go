package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// NewRoom returns a new instance of a Room
func NewRoom() (*Room, error) {
	return &Room{
		mux:            sync.Mutex{},
		connectionsMap: make(map[*websocket.Conn]int),
		usernamesMap:   make(map[string]int),
	}, nil
}

// Room holds all of the connections
type Room struct {
	mux            sync.Mutex
	connections    []*websocket.Conn
	connectionsMap map[*websocket.Conn]int
	usernames      []string
	usernamesMap   map[string]int
}

// Broadcast takes a message and sends it to all in the room
func (r *Room) Broadcast(conn *websocket.Conn, msg string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	pos := r.connectionsMap[conn]
	username := r.usernames[pos]

	log.Printf("Got message from %s \n Message: \n\t %s", username, msg)

	message := NewUserMsg(username, msg)

	for _, conn := range r.connections {
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

	// check that username doesn't already exist
	if _, exists := r.usernamesMap[username]; exists {
		return fmt.Errorf(`username "%s" already taken`, username)
	}

	r.connections = append(r.connections, conn)
	r.connectionsMap[conn] = len(r.connections) - 1
	r.usernames = append(r.usernames, username)
	r.usernamesMap[username] = len(r.usernames) - 1

	return nil
}

// Remove removes the connection from the room
func (r *Room) Remove(conn *websocket.Conn) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	pos, exists := r.connectionsMap[conn]
	if !exists {
		return nil
	}

	r.connections = append(r.connections[:pos], r.connections[pos+1:]...)
	delete(r.connectionsMap, conn)
	username := r.usernames[pos]
	r.usernames = append(r.usernames[:pos], r.usernames[pos+1:]...)
	delete(r.usernamesMap, username)

	return nil
}
