import React, { useEffect, useRef } from "react";

const JoinRoom = (props) => {
  const webSocketRef = useRef();
  const peerConnectionRef = useRef();
  const dataChannelRef = useRef();

  const getUserMedia = async () => {
    try {
      return await navigator.mediaDevices.getUserMedia({ video: true, audio: false })
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    getUserMedia().then((stream) => {
      document.getElementById('localVideo').srcObject = stream

      peerConnectionRef.current = createPeerConnection()
      dataChannelRef.current = createDataChannel()
      peerConnectionRef.current.ondatachannel = dataChannel => {
        dataChannel.onmessage = handleMessageOnDataChannel
      }

      stream.getTracks().forEach((track) => {
        peerConnectionRef.current.addTrack(track, stream);
      });

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

  const createDataChannel = () => {
    const dataChannel = peerConnectionRef.current.createDataChannel("ClientMessage")

    return dataChannel
  }

  const handleMessageOnDataChannel = (e) => {
    console.log(e)
  }

  const handleICECandidateEvent = (e) => {
    if (!e.candidate) {
      // All candidates have been sent
      return
    }

    console.log("Sending ice candidate")
    webSocketRef.current.send(
      JSON.stringify({ event: 'candidate', data: JSON.stringify(e.candidate) })
    );
  }

  const handleNewICECandidate = async (candidate) => {
    console.log("Receive ice candidate")
    console.log(candidate)

    peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(candidate))
  }

  const handleVideoOfferMsg = async (offer) => {
    console.log("Received offer")
    console.log(offer)

    await peerConnectionRef.current.setRemoteDescription(
      new RTCSessionDescription(offer)
    )

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
        console.log("Remove track")
        el.parentNode.removeChild(el)
      }
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    const message = document.getElementById('messageBox').value
    if (message == "") {
      return
    }

    document.getElementById('messageBox').value = ""
    dataChannelRef.current.send(message)
  };

  return (
    <div>
      <h3> Chat </h3>
      <ul id="messages"></ul>
      <input id="messageBox" />
      <button onClick={sendMessage}>Send message</button>

      <h3> Local Video </h3>
      <video id="localVideo" width="160" height="120" autoPlay muted></video> <br />

      <h3> Remote Video </h3>
      <div id="remoteVideos"></div> <br />
    </div>
  );
};

export default JoinRoom;

