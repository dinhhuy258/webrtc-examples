package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
  defer websocket.Close()

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
  }
}

