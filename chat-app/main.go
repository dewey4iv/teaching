package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	log.Printf("Starting the Chat Server")

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("./public/")).ServeHTTP(w, r)
	})

	chatServer, err := New(WithUpgrader(websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
	}))

	if err != nil {
		panic(err)
	}

	router.Handle("/websocket", chatServer)
	router.Handle("/websocket/", chatServer)

	http.ListenAndServe(":8080", router)
}
