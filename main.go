package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/gorilla/websocket"
)

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Websockets are not supported yet")
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/websocket", websocketEndpoint)
}

// http server from html file
func main() {
	log.Println("Starting http server...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":4444", nil))
}
