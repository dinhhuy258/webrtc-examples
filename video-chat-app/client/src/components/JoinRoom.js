import React, { useEffect, useRef } from "react";

const JoinRoom = (props) => {
  const webSocketRef = useRef();

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

      webSocketRef.current = new WebSocket(
        `ws://localhost:8080/join?roomID=${props.match.params.roomID}`
      );

      webSocketRef.current.addEventListener("message", async (e) => {
        const message = JSON.parse(e.data);
        if (!message) {
          return console.log('Failed to parse msg')
        }

        console.log(message)
      });
    });
  });

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

