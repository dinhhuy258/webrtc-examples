import React, { useEffect, useRef, useState } from "react";

const usernames = ["dinhhuy258", "huyduong147"]

const VideoCall = (props) => {
  const webSocketRef = useRef();
  const peerConnectionRef = useRef()
  const [receivingCall, setReceivingCall] = useState(false);
  const [caller, setCaller] = useState("");
  const [callerSignal, setCallerSignal] = useState();
  const [callAccepted, setCallAccepted] = useState(false);
  const username = props.match.params.username
  const anotherUsername = usernames.find(u => u != username)

  if (!webSocketRef.current) {
    webSocketRef.current = new WebSocket(
      `ws://localhost:8080/connect?username=${username}`
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
    document.getElementById('remoteVideo').srcObject = e.streams[0]
  };

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

  const createPeerConnection = () => {
    const peerConnection = new RTCPeerConnection({});
    peerConnection.onicecandidate = handleICECandidateEvent;
    peerConnection.ontrack = handleTrackEvent;

    return peerConnection;
  }

  const handleNewICECandidate = async (candidate) => {
    console.log("Receive ice candidate")
    console.log(candidate)

    peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(candidate))
  }

  const handleAnswerMessage = async (answer) => {
    console.log("Callee has answered")

    var answerDesc = new RTCSessionDescription(JSON.parse(answer))
    peerConnectionRef.current.setRemoteDescription(new RTCSessionDescription(answerDesc))
  }

  const handleVideoOfferMsg = async (offer) => {
    console.log("Received offer")

    var offerDesc = new RTCSessionDescription(JSON.parse(offer));
    await peerConnectionRef.current.setRemoteDescription(
      new RTCSessionDescription(offerDesc)
    )

    const answer = await peerConnectionRef.current.createAnswer();
    await peerConnectionRef.current.setLocalDescription(answer);

    webSocketRef.current.send(
      JSON.stringify(
        {
          event: 'answer',
          data: JSON.stringify(
            {
              caller: anotherUsername,
              callee: username,
              sdp: JSON.stringify(peerConnectionRef.current.localDescription)
            }
          )
        })
    );
  }

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

      stream.getTracks().forEach((track) => {
        peerConnectionRef.current.addTrack(track, stream);
      });
    });

    webSocketRef.current.onmessage = (event) => {
      const message = JSON.parse(event.data);

      switch (message.event) {
        case "call":
          setReceivingCall(true);

          const messageData = JSON.parse(message.data);
          setCaller(messageData.caller);
          setCallerSignal(messageData.sdp);

          break;
        case "answer":
          handleAnswerMessage(message.data)

          break;
        case "candidate":
          let candidate = JSON.parse(message.data)
          if (!candidate) {
            return console.log('failed to parse candidate')
          }

          handleNewICECandidate(candidate)
          break;
        case "close":
          console.log("close");
          break;
        case "message":
          alert(message.data)
          break;
        default:
          break;
      }
    };
  }, []);

  function acceptCall() {
    setReceivingCall(false);

    handleVideoOfferMsg(callerSignal);
  }

  function rejectCall() {
    setReceivingCall(false);
  }

  const handleCall = async (e) => {
    e.preventDefault();

    peerConnectionRef.current.createOffer().then(async (offer) => {
      await peerConnectionRef.current.setLocalDescription(offer);
    }).then(() => {
      webSocketRef.current.send(JSON.stringify({
        event: "call",
        data: JSON.stringify(
          {
            caller: username,
            callee: anotherUsername,
            sdp: JSON.stringify(peerConnectionRef.current.localDescription)
          }
        )
      }));
    });
  }

  let incomingCall;
  if (receivingCall) {
    incomingCall = (
      <div>
        <h1>{caller} is calling you</h1>
        <button onClick={acceptCall}>Accept</button>
        <button onClick={rejectCall}>Reject</button>
      </div>
    )
  }

  return (
    <div>
      <h1>{username}</h1>
      <button onClick={handleCall}>Call {anotherUsername}</button>
      {incomingCall}
      <h3> Local Video </h3>
      <video id="localVideo" width="160" height="120" autoPlay muted></video> <br />

      <h3> Remote Video </h3>
      <video id="remoteVideo" autoPlay muted></video> <br />
    </div>
  );
};

export default VideoCall;
