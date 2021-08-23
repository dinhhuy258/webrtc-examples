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

        if (message.event == "offer") {
          let offer = JSON.parse(message.data)
          if (!offer) {
            return console.log('failed to parse answer')
          }

          handleVideoOfferMsg(offer)
        }
      });
    });
  });

  const createPeerConnection = () => {
    const peerConnection = new RTCPeerConnection({});

    return peerConnection
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

