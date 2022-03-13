import React from "react";
import { useHistory, useLocation } from 'react-router-dom'

import { Room, VideoPresets } from 'livekit-client'
import { DisplayContext, LiveKitRoom } from 'livekit-react'

function RoomPage() {
  const history = useHistory()
  const query = new URLSearchParams(useLocation().search)
  const token = query.get('token')

  function onLeave() {
    history.push({
      pathname: '/',
    })
  }

  async function onConnected(room: Room) {
    await room.localParticipant.setCameraEnabled(true);
  }

  let displayOptions = {
    stageLayout: 'grid',
    showStats: false,
  }

  return (
    <DisplayContext.Provider value={displayOptions}>
      <LiveKitRoom
        url='ws://localhost:7880'
        token={token}
        onConnected={room => {
          onConnected(room);
        }}
        connectOptions={{
          adaptiveStream: true,
          dynacast: true,
          videoCaptureDefaults: {
            resolution: VideoPresets.hd.resolution,
          },
          publishDefaults: {
            videoEncoding: VideoPresets.hd.encoding,
            simulcast: true,
          },
          logLevel: 'debug',
        }}
        onLeave={onLeave}
      />
    </DisplayContext.Provider>
  );
}

export default RoomPage;

