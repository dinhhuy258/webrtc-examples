import React, { useEffect, useRef } from "react";

const JoinRoom = (props) => {
  const webSocketRef = useRef();
  const peerConnectionRef = useRef();
  const userStream = useRef();

  const getUserMedia = async () => {
    try {
      return await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    getUserMedia().then((stream) => {
      document.getElementById('localVideo').srcObject = stream
      userStream.current = stream;

      webSocketRef.current = new WebSocket(
        `ws://localhost:8080/join?roomID=${props.match.params.roomID}`
      );

      webSocketRef.current.addEventListener("message", async (e) => {
        const message = JSON.parse(e.data);
        if (!message) {
          return console.log('Failed to parse msg')
        }

        if (message.event === "offer") {
          let offer = JSON.parse(message.data)
          if (!offer) {
            return console.log('failed to parse answer')
          }

          handleVideoOfferMsg(offer)
        }
        else if (message.event === "candidate") {
          let candidate = JSON.parse(message.data)
          if (!candidate) {
            return console.log('failed to parse candidate')
          }

          handleNewICECandidate(candidate)
        }
      });
    });
  });

  const createPeerConnection = () => {
    const peerConnection = new RTCPeerConnection({});
    peerConnection.onicecandidate = handleICECandidateEvent;
    peerConnection.ontrack = handleTrackEvent;

    return peerConnection;
  }

  const handleICECandidateEvent = (e) => {
    if (!e.candidate) {
      // All candidates have been sent
      return
    }

    webSocketRef.current.send(
      JSON.stringify({ event: 'candidate', data: JSON.stringify(e.candidate) })
    );
  }

  const handleNewICECandidate = async (candidate) => {
    peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(candidate))
  }

  const handleVideoOfferMsg = async (offer) => {
    peerConnectionRef.current = createPeerConnection()

    await peerConnectionRef.current.setRemoteDescription(
      new RTCSessionDescription(offer)
    )

    userStream.current.getTracks().forEach((track) => {
      peerConnectionRef.current.addTrack(track, userStream.current);
    });

    const answer = await peerConnectionRef.current.createAnswer();
    await peerConnectionRef.current.setLocalDescription(answer);

    webSocketRef.current.send(
      JSON.stringify({ event: 'answer', data: JSON.stringify(answer) })
    );
  }

  const handleTrackEvent = (e) => {
    console.log("Received track")
    console.log(e)
    if (e.track.kind === 'audio') {
      return
    }

    let el = document.createElement(e.track.kind)
    el.srcObject = e.streams[0]
    el.autoplay = true
    el.controls = true
    document.getElementById('remoteVideos').appendChild(el)
    e.track.onmute = function(e) {
      el.play()
    }

    e.streams[0].onremovetrack = ({ track }) => {
      if (el.parentNode) {
        el.parentNode.removeChild(el)
      }
    }
  };

  return (
    <div>
      <h3> Local Video </h3>
      <video id="localVideo" width="160" height="120" autoPlay muted></video> <br />

      <h3> Remote Video </h3>
      <div id="remoteVideos"></div> <br />
    </div>
  );
};

export default JoinRoom;

