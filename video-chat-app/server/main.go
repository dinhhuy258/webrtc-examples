package main

import (
	"video-chat-app/internal"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	roomManager = internal.RoomManager{}

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type RoomResponseData struct {
	RoomId string `json:"room_id"`
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		RoomResponseData{
			RoomId: roomManager.CreateRoom(),
		},
	)
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomIDQuery, ok := r.URL.Query()["roomID"]
	if !ok {
		return
	}
	roomID := roomIDQuery[0]

	if !roomManager.HasRoom(roomID) {
    return
	} 
  _, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    panic(err)
  }
}

func main() {
	roomManager.Init()

	http.HandleFunc("/create-room", CreateRoomRequestHandler)
	http.HandleFunc("/join", JoinRoomRequestHandler)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
