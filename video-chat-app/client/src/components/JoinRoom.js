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
      console.log(stream)
      webSocketRef.current = new WebSocket(
        `ws://localhost:8080/join?roomID=${props.match.params.roomID}`
      );
    });
  });

  return (
    <div>
      JoinRoom
    </div>
  );
};

export default JoinRoom;

