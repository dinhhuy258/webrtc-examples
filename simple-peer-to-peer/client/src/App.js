import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import SelectUser from "./components/SelectUser"
import JoinCall from "./components/JoinCall"

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Switch>
          <Route path="/" exact component={SelectUser}></Route>
          <Route path="/room/:roomID" component={JoinCall}></Route>
        </Switch>
      </BrowserRouter>
    </div>
  );
}

export default App;
