import React from "react";

const SelectUser = (props) => {
  const select = async (username, e) => {
    e.preventDefault();

    props.history.push(`/video-call/${username}`)
  };

  return (
    <div>
      <div>
        Select user to login:
      </div>
      <div>
        <button onClick={select.bind(this, "dinhhuy258")}>dinhhuy258</button>
        <button onClick={select.bind(this, "huyduong147")}>huyduong147</button>
      </div>
    </div>
  );
};

export default SelectUser;
