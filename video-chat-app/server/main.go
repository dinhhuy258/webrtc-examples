package main

import (
	"time"
	"video-chat-app/internal"

	"log"
	"net/http"
)

func main() {
	internal.RoomManagerInstance.Init()

	http.HandleFunc("/create-room", internal.CreateRoomRequestHandler)
	http.HandleFunc("/join", internal.JoinRoomRequestHandler)

	// request a keyframe every 3 seconds
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			internal.DispatchKeyFrames()
		}
	}()

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
