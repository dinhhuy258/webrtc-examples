package internal

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomManager struct {
	Mutex sync.RWMutex
	Rooms map[string][]Participant
}

func (rm *RoomManager) CreateRoom() string {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	// Generate room id
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	roomID := string(b)

	// Add to RoomManager
	rm.Rooms[roomID] = []Participant{}

	return roomID
}
