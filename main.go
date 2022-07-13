package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var rdb *redis.Client
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

	sendPreviousMessages(conn)
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

		log.Println("Received message: ", string(chatMessage.Text))
		broadcast <- chatMessage
	}
}

func sendPreviousMessages(conn *websocket.Conn) {
	messages, err := rdb.LRange("messages", 0, -1).Result()
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messages {
		var chatMessage ChatMessage
		err := json.Unmarshal([]byte(message), &chatMessage)
		if err != nil {
			panic(err)
		}

		messageClient(conn, chatMessage)
	}
}

func handleMessages() {
	for {
		chatMessage := <-broadcast

		storeMessageInDatabase(chatMessage)
		messageClients(chatMessage)
	}
}

func storeMessageInDatabase(message ChatMessage) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	if err := rdb.RPush("messages", jsonMessage).Err(); err != nil {
		panic(err)
	}
}

func messageClients(message ChatMessage) {
	for client := range clients {
		messageClient(client, message)
	}
}

func messageClient(conn *websocket.Conn, message ChatMessage) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println(err)
		conn.Close()
		delete(clients, conn)
	}
}

func setupDatabase() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/websocket", websocketEndpoint)
}

func main() {
	log.Println("Starting http server...")

	setupDatabase()
	setupRoutes()
	go handleMessages()

	log.Fatal(http.ListenAndServe(":4444", nil))
}
