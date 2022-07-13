package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan ChatMessage)

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()
	clients[conn] = true
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		var chatMessage ChatMessage
		err := conn.ReadJSON(&chatMessage)
		if err != nil {
			log.Println(err)
			return
		}

		// Broadcast message to all connected clients
		log.Println("Received message: ", string(chatMessage.Text))
		broadcast <- chatMessage
	}
}

func handleMessages() {
	for {
		chatMessage := <-broadcast
		for client := range clients {
			err := client.WriteJSON(chatMessage)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/websocket", websocketEndpoint)
}

func main() {
	log.Println("Starting http server...")

	setupRoutes()
	go handleMessages()

	log.Fatal(http.ListenAndServe(":4444", nil))
}
