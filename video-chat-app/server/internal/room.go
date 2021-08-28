package internal

import (
	"math/rand"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

var RoomManagerInstance = RoomManager{}

type Participant struct {
	Websocket      *ThreadSafeWriter
	PeerConnection *webrtc.PeerConnection
  DataChannel    *webrtc.DataChannel
}

type Room struct {
	ID           string
	Mutex        sync.RWMutex
	Participants []Participant
	TrackLocals  map[string]*webrtc.TrackLocalStaticRTP
}

type RoomManager struct {
	Mutex sync.RWMutex
	Rooms map[string]*Room
}

func (rm *RoomManager) Init() {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	rm.Rooms = make(map[string]*Room)
}

func (rm *RoomManager) CreateRoom() string {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	// Generate room id
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	roomID := string(b)

	// Add to RoomManager
	rm.Rooms[roomID] = &Room{
		ID:           roomID,
		Mutex:        sync.RWMutex{},
		Participants: []Participant{},
		TrackLocals:  map[string]*webrtc.TrackLocalStaticRTP{},
	}

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

func (rm *RoomManager) GetRoom(roomID string) *Room {
	rm.Mutex.RLock()
	defer rm.Mutex.RUnlock()

	return rm.Rooms[roomID]
}

func (rm *RoomManager) Join(
	roomID string,
	websocket *ThreadSafeWriter,
	peerConnection *webrtc.PeerConnection,
  dataChannel *webrtc.DataChannel,
) *Room {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	participant := Participant{
		websocket,
		peerConnection,
    dataChannel,
	}

	rm.Rooms[roomID].Mutex.Lock()
	defer rm.Rooms[roomID].Mutex.Unlock()
	rm.Rooms[roomID].Participants = append(rm.Rooms[roomID].Participants, participant)

	return rm.Rooms[roomID]
}
