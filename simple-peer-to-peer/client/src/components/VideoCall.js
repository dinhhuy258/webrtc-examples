import React, { useEffect, useRef, useState } from "react";

const usernames = ["dinhhuy258", "huyduong147"]

const VideoCall = (props) => {
  const webSocketRef = useRef();
  const peerConnectionRef = useRef()
  const [receivingCall, setReceivingCall] = useState(false);
  const [caller, setCaller] = useState("");
  const username = props.match.params.username
  const anotherUsername = usernames.find(u => u !== username)

  if (!webSocketRef.current) {
    webSocketRef.current = new WebSocket(
      `ws://localhost:8080/connect?username=${username}`
    );
  }

  const handleTrackEvent = (e) => {
    if (e.track.kind === 'audio') {
      return
    }

    console.log("Received remote video track")
    document.getElementById('remoteVideo').srcObject = e.streams[0]
  };

  const handleICECandidateEvent = (e) => {
    if (!e.candidate) {
      // All candidates have been sent
      return
    }

    webSocketRef.current.send(
      JSON.stringify({
        event: 'ice-candidate', data: JSON.stringify(
          {
            target: anotherUsername,
            candidate: JSON.stringify(e.candidate)
          }
        )
      })
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

  const handleOfferMessage = async (offer) => {
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
              target: anotherUsername,
              sdp: JSON.stringify(peerConnectionRef.current.localDescription)
            }
          )
        })
    )
  }

  const handleAnswerMessage = async (answer) => {
    var answerDesc = new RTCSessionDescription(JSON.parse(answer))
    peerConnectionRef.current.setRemoteDescription(new RTCSessionDescription(answerDesc))
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
        case "request-call":
          setReceivingCall(true);
          const messageData = JSON.parse(message.data);
          setCaller(messageData.caller);

          break;
        case "offer":
          console.log("Receive offer")
          handleOfferMessage(message.data)

          break;
        case "answer":
          console.log("Receive answer")
          handleAnswerMessage(message.data)

          break;
        case "ice-candidate":
          let candidate = JSON.parse(message.data)
          if (!candidate) {
            return console.log('failed to parse candidate')
          }

          handleNewICECandidate(candidate)

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

    // Create an offer and send it to the caller
    peerConnectionRef.current.createOffer().then(async (offer) => {
      await peerConnectionRef.current.setLocalDescription(offer);
    }).then(() => {
      webSocketRef.current.send(
        JSON.stringify(
          {
            event: 'offer',
            data: JSON.stringify(
              {
                target: anotherUsername,
                sdp: JSON.stringify(peerConnectionRef.current.localDescription)
              }
            )
          })
      );
    });
  }

  function rejectCall() {
    setReceivingCall(false);
  }

  const handleCall = async (e) => {
    e.preventDefault();

    webSocketRef.current.send(JSON.stringify({
      event: "request-call",
      data: JSON.stringify(
        {
          caller: username,
          callee: anotherUsername,
        }
      )
    }));
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
      <video id="localVideo" width="128" height="128" autoPlay muted></video> <br />

      <h3> Remote Video </h3>
      <video id="remoteVideo" width="256" height="256" autoPlay muted></video> <br />
    </div>
  );
};

export default VideoCall;
