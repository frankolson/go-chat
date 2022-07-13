package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		// Read message from browser
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Broadcast message to all connected clients
		log.Println("Received message: ", string(p))
	}
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/websocket", websocketEndpoint)
}

func main() {
	log.Println("Starting http server...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":4444", nil))
}
