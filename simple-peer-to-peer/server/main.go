package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var webSocketUsersMap = make(map[string]*websocket.Conn)

func main() {
	http.HandleFunc("/connect", connectHandler)

	err := http.ListenAndServe(":8080", nil)
	log.Println("Starting server on port 8080")

	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type callMessage struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	username, ok := r.URL.Query()["username"]
	if !ok {
		return
	}

	log.Println("Login with username " + username[0])

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	webSocketUsersMap[username[0]] = websocket

	defer func() {
		websocket.Close()
		webSocketUsersMap[username[0]] = nil
	}()

	message := &websocketMessage{}
	for {
		_, raw, err := websocket.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}

		switch message.Event {
		case "call":
			callMessage := &callMessage{}
			if err := json.Unmarshal([]byte(message.Data), &callMessage); err != nil {
				return
			}

			if calleeWebsocket, ok := webSocketUsersMap[callMessage.Callee]; ok {
				log.Println(callMessage.Caller + " call to " + callMessage.Callee)

				calleeWebsocket.WriteJSON(&websocketMessage{
					Event: "call",
					Data:  callMessage.Caller,
				})
			} else {
				log.Println(callMessage.Callee + " is not online")

				webSocketUsersMap[callMessage.Caller].WriteJSON(&websocketMessage{
					Event: "message",
					Data:  "User " + callMessage.Callee + " is not online",
				})
			}

			break
		}
	}
}
