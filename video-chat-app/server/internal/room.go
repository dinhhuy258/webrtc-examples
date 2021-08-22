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

func (rm *RoomManager) Init() {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	rm.Rooms = make(map[string][]Participant)
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

func (rm *RoomManager) HasRoom(roomID string) bool {
	rm.Mutex.RLock()
	defer rm.Mutex.RUnlock()

	if _, ok := rm.Rooms[roomID]; ok {
		return true
	}

	return false
}

func (rm *RoomManager) Join(roomID string, host bool, conn *websocket.Conn) {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	rm.Rooms[roomID] = append(rm.Rooms[roomID], Participant{
		host,
		conn,
	})
}
