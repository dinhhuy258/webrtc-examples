package main

import (
	"video-chat-app/internal"

	"log"
	"net/http"
)

func main() {
	internal.RoomManagerInstance.Init()

	http.HandleFunc("/create-room", internal.CreateRoomRequestHandler)
	http.HandleFunc("/join", internal.JoinRoomRequestHandler)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
