import React, { useEffect, useRef } from "react";

const usernames = ["dinhhuy258", "huyduong147"]

const VideoCall = (props) => {
  const webSocketRef = useRef();
  const username = props.match.params.username
  const anotherUsername = usernames.find(u => u != username)

  webSocketRef.current = new WebSocket(
    `ws://localhost:8080/connect?username=${username}`
  );

  webSocketRef.current.addEventListener("message", async (e) => {
    const message = JSON.parse(e.data);
    if (!message) {
      return console.log('Failed to parse msg')
    }

    switch (message.event) {
      case "call":
        console.log(message.data)
        break;
      case "message":
        alert(message.data)
        break;
    }
  });

  const handleCall = async (e) => {
    e.preventDefault();

    webSocketRef.current.send(JSON.stringify({
      event: "call", data: '{ "caller": "' + username + '", "callee": "' + anotherUsername + '" }'
    }));
  };

  return (
    <div>
      <h1>{username}</h1>
      <button onClick={handleCall}>Call {anotherUsername}</button>
      <h3> Local Video </h3>
      <video id="localVideo" width="160" height="120" autoPlay muted></video> <br />

      <h3> Remote Video </h3>
      <video id="remoteVideo" autoPlay muted></video> <br />
    </div>
  );
};

export default VideoCall;

