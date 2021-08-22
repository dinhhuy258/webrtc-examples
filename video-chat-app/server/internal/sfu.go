package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
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
			RoomId: RoomManagerInstance.CreateRoom(),
		},
	)
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomIDQuery, ok := r.URL.Query()["roomID"]
	if !ok {
		return
	}
	roomID := roomIDQuery[0]

	if !RoomManagerInstance.HasRoom(roomID) {
		return
	}
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
}
