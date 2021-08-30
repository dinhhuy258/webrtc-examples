import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import SelectUser from "./components/SelectUser"
import VideoCall from "./components/VideoCall"

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Switch>
          <Route path="/" exact component={SelectUser}></Route>
          <Route path="/video-call/:username" component={VideoCall}></Route>
        </Switch>
      </BrowserRouter>
    </div>
  );
}

export default App;
