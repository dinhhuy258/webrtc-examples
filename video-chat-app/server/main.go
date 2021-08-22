package main

import (
	"encoding/json"
	"log"
	"net/http"
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
			RoomId: "123123123",
		},
	)
}

func main() {
	http.HandleFunc("/create-room", CreateRoomRequestHandler)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
