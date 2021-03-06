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

type requestCallMessage struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

type signallingMessage struct {
	Target string `json:"target"`
	Sdp    string `json:"sdp"`
}

type iceCandidateMessage struct {
	Target    string `json:"target"`
	Candidate string `json:"candidate"`
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
		case "ice-candidate":
			log.Println("Received ICE candidate")

			iceCandidateMessage := &iceCandidateMessage{}
			if err := json.Unmarshal([]byte(message.Data), &iceCandidateMessage); err != nil {
				return
			}

			if targetWebsocket, ok := webSocketUsersMap[iceCandidateMessage.Target]; ok {
				targetWebsocket.WriteJSON(&websocketMessage{
					Event: "ice-candidate",
					Data:  iceCandidateMessage.Candidate,
				})
			} else {
				log.Println(iceCandidateMessage.Target + " is not online")
			}

			break
		case "offer", "answer":
			signallingMessage := &signallingMessage{}
			if err := json.Unmarshal([]byte(message.Data), &signallingMessage); err != nil {
				return
			}

			if callerWebsocket, ok := webSocketUsersMap[signallingMessage.Target]; ok {
				callerWebsocket.WriteJSON(&websocketMessage{
					Event: message.Event,
					Data:  signallingMessage.Sdp,
				})
			} else {
				log.Println(signallingMessage.Target + " is not online")
			}

			break
		case "request-call":
			requestCallMessage := &requestCallMessage{}
			if err := json.Unmarshal([]byte(message.Data), &requestCallMessage); err != nil {
				return
			}

			if calleeWebsocket, ok := webSocketUsersMap[requestCallMessage.Callee]; ok {
				log.Println(requestCallMessage.Caller + " call to " + requestCallMessage.Callee)

				calleeWebsocket.WriteJSON(&websocketMessage{
					Event: "request-call",
					Data:  message.Data,
				})
			} else {
				log.Println(requestCallMessage.Callee + " is not online")

				webSocketUsersMap[requestCallMessage.Caller].WriteJSON(&websocketMessage{
					Event: "message",
					Data:  "User " + requestCallMessage.Callee + " is not online",
				})
			}

			break
		}
	}
}
