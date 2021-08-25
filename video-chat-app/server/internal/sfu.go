package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
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

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RoomResponseData struct {
	RoomId string `json:"room_id"`
}

func DispatchKeyFrames() {
  RoomManagerInstance.Mutex.Lock()
  defer RoomManagerInstance.Mutex.Unlock()

  for _, room := range RoomManagerInstance.Rooms {
    dispatchKeyFrame(room)
  }
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

	defer peerConnection.Close()

	room := RoomManagerInstance.Join(roomID, conn, peerConnection)

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			// all candidates have been sent
			return
		}

		candidatestring, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println("failed to parsr ice candidate %w", err)
			return
		}

		conn.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidatestring),
		})
	})

	peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		switch pcs {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			// If PeerConnection is closed remove it from global list
			signalPeerConnections(room)
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := addTrack(room, t)
		defer removeTrack(room, trackLocal)

		buf := make([]byte, 1500)

		for {
			i, _, err := t.Read(buf)
			if err != nil {
        log.Println(err)
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
        log.Println(err)
				return
			}
		}
	})

	signalPeerConnections(room)

	message := &websocketMessage{}
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}

		switch message.Event {
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Println(err)
				return
			}
			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Printf("failed to add ice candidate %v", err)
				return
			}
		}
	}
}

func createPeerConnection() *webrtc.PeerConnection {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Printf("failed to create peer connection %v", err)
		return nil
	}

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

func addTrack(room *Room, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	room.Mutex.Lock()
	defer func() {
		room.Mutex.Unlock()
    signalPeerConnections(room)
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	room.TrackLocals[t.ID()] = trackLocal

	return trackLocal
}

func removeTrack(room *Room, t *webrtc.TrackLocalStaticRTP) {
	room.Mutex.Lock()
	defer func() {
		room.Mutex.Unlock()
    signalPeerConnections(room)
	}()

	delete(room.TrackLocals, t.ID())
}

func dispatchKeyFrame(room *Room) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	for i := range room.Participants {
		peerConnection := room.Participants[i].PeerConnection
		for _, receiver := range peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}
			_ = peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}

// signalPeerConnections updates each PeerConnection so that it getting all the expected media tracks
func signalPeerConnections(room *Room) {
	room.Mutex.Lock()
	defer func() {
		room.Mutex.Unlock()
    dispatchKeyFrame(room)
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range room.Participants {
			participant := room.Participants[i]

			if participant.PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				// Delete closed connection
				room.Participants = append(room.Participants[:i], room.Participants[i+1:]...)
				return true // We modified the slice, start from the beginning
			}

			// map of sender we already are seanding, so we don't double send
			existingSenders := map[string]bool{}
			for _, sender := range participant.PeerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := room.TrackLocals[sender.Track().ID()]; !ok {
					if err := participant.PeerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range participant.PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range room.TrackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := participant.PeerConnection.AddTrack(room.TrackLocals[trackID]); err != nil {
						return true
					}
				}
			}

			sdp, err := participant.PeerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = participant.PeerConnection.SetLocalDescription(sdp); err != nil {
				return true
			}

			// Send offer to Websocket
			offerString, err := json.Marshal(sdp)
			if err != nil {
				return true
			}

			if err = participant.Websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}

		return false
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				signalPeerConnections(room)
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
