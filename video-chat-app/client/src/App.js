import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import CreateRoom from "./components/CreateRoom"
import JoinRoom from "./components/JoinRoom"

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Switch>
          <Route path="/" exact component={CreateRoom}></Route>
          <Route path="/room/:roomID" component={JoinRoom}></Route>
        </Switch>
      </BrowserRouter>
    </div>
  );
}

export default App;
