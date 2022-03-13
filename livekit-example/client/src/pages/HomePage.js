import React from "react";
import { useHistory } from 'react-router-dom'

function HomePage() {
  const history = useHistory()

  function connectToRoom(token) {
    const params = {
      token,
    }

    history.push({
      pathname: '/room',
      search: "?" + new URLSearchParams(params).toString()
    })
  }

  return (
    <div className="prejoin">
      <main>
        <h2>LiveKit Video</h2>
        <hr />

        <button onClick={connectToRoom.bind(this, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDk3MzM5NjUsImlzcyI6IkFQSXVQZXhQeDlZYVl1NiIsImp0aSI6Ikh1eSBEdW9uZyAtIDEiLCJuYmYiOjE2NDcxNDE5NjUsInN1YiI6Ikh1eSBEdW9uZyAtIDEiLCJ2aWRlbyI6eyJyb29tIjoiRGVtbyBsaXZla2l0Iiwicm9vbUpvaW4iOnRydWV9fQ.Q67CDTdPQPY3yliB7B3hhASZYuSk3Ab31yea_jv07Y0")} >dinhhuy258</button>
        <button onClick={connectToRoom.bind(this, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDk3MzQwMDAsImlzcyI6IkFQSXVQZXhQeDlZYVl1NiIsImp0aSI6Ikh1eSBEdW9uZyAtIDIiLCJuYmYiOjE2NDcxNDIwMDAsInN1YiI6Ikh1eSBEdW9uZyAtIDIiLCJ2aWRlbyI6eyJyb29tIjoiRGVtbyBsaXZla2l0Iiwicm9vbUpvaW4iOnRydWV9fQ.pOWpwvwa1_RuEW283lnjS9Vg9--rylz_MdAVnGgPi94")} >huyduong147</button>
      </main>
    </div>
  );
}

export default HomePage;

