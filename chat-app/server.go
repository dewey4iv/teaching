package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dewey4iv/teaching/chat-app/payload"
	room "github.com/dewey4iv/teaching/chat-app/rooms/simple"
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
		room, err := room.New()
		if err != nil {
			return nil, err
		}

		c.room = room
	}

	return &c, nil
}

// ChatServer is the Chat ChatServer
type ChatServer struct {
	upgrader websocket.Upgrader
	room     Room
}

func (s *ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ensure unique username

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleErr(w, fmt.Errorf("error upgrading to websockets: %s", err.Error()), http.StatusConflict)
		return
	}

	if err := s.setup(w, r, conn); err != nil {
		conn.WriteJSON(payload.NewServerMsg("close-message", err.Error()))
		return
	}

	for {
		// TODO: maybe this should be a payload??
		var txt struct {
			Message string `json:"message"`
		}

		if err := conn.ReadJSON(&txt); err != nil {
			log.Printf("error reading from connenction:\n %s", err.Error())
			s.cleanupConn(conn)
			return
		}

		s.room.Broadcast(conn, txt.Message)
	}
}

func (s *ChatServer) setup(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) error {
	username, err := getUsername(r)
	if err != nil {
		return err
	}

	if err := s.room.Add(conn, username); err != nil {
		payload := payload.NewServerMsg("close-message", err.Error())
		conn.WriteJSON(payload)
		s.cleanupConn(conn)
		return err
	}

	welcome := payload.New("welcome", "Chat Bot", "")

	if err := conn.WriteJSON(welcome); err != nil {
		s.cleanupConn(conn)
		return err
	}

	if err := s.room.Broadcast(conn, fmt.Sprintf("%s has joined the room!", username)); err != nil {
		s.cleanupConn(conn)
		return err
	}

	return nil
}

type Room interface {
	Broadcast(conn *websocket.Conn, msg string) error
	Add(conn *websocket.Conn, username string) error
	Remove(conn *websocket.Conn) error
}

func (s *ChatServer) cleanupConn(conn *websocket.Conn) {
	log.Printf("closing connection")
	conn.Close()
	s.room.Remove(conn)
}

func handleErr(w http.ResponseWriter, err error, status int) {
	if status == 0 {
		status = http.StatusNotFound
	}

	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
}

func getUsername(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	last := parts[len(parts)-1]
	if last == "websocket" {
		return "", fmt.Errorf("can't get username")
	}

	return last, nil
}
