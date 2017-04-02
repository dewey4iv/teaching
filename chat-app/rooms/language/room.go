package room

import (
	"fmt"
	"strings"

	room "github.com/dewey4iv/teaching/chat-app/rooms/simple"
	"github.com/gorilla/websocket"
)

// New takes a set of options and returns a new Room
func New(opts ...Option) (*Room, error) {
	var r Room

	simple, err := room.New()
	if err != nil {
		return nil, err
	}

	r.Room = simple

	for _, opt := range opts {
		if err := opt.Apply(&r); err != nil {
			return nil, err
		}
	}

	return &r, nil
}

// Room wraps a simple room and makes
// sure no bad words are broadcasted
type Room struct {
	filterWords []string
	*room.Room
}

// Broadcast checks that a word isn't in the list and then broadcasts the message
func (r *Room) Broadcast(conn *websocket.Conn, msg string) error {
	for i := range r.filterWords {
		if strings.Contains(msg, r.filterWords[i]) {
			return fmt.Errorf("can't use word: %s", r.filterWords[i])
		}
	}

	return r.Room.Broadcast(conn, msg)
}
