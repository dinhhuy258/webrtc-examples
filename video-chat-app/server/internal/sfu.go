package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// Helper to make Gorilla Websockets threadsafe
type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(data interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(data)
}

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
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := &ThreadSafeWriter{unsafeConn, sync.Mutex{}}
	defer conn.Close()

	peerConnection := createPeerConnection()
	if peerConnection == nil {
		return
	}

	RoomManagerInstance.Join(roomID, conn, peerConnection)
}

func createPeerConnection() *webrtc.PeerConnection {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Printf("failed to create peer connection %v", err)
		return nil
	}

	defer peerConnection.Close()

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Printf("failed to add transceiver from kind %v", err)
			return nil
		}
	}

	return peerConnection
}
