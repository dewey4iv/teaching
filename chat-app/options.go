package main

import "github.com/gorilla/websocket"

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

// WithRoom sets the room on the ChatServer
func WithRoom(room Room) Option {
	return &withRoom{room}
}

type withRoom struct {
	room Room
}

func (opt *withRoom) Apply(c *ChatServer) error {
	c.room = opt.room

	return nil
}
