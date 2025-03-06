import React, { useState } from "react";
import * as Constants from './constants/Constants';

import '../css/StartScreen.css';
import '../css/Shared.css';
import ArchitortureLogo from '../res/logo.png';

function StartScreen({ startGameCallback }) {
  const [newGame, setNewGame] = useState(false);
  const [joinRoom, setJoinRoom] = useState(false);
  const [username, setUsername] = useState("");
  const [error, setError] = useState("");
  const [roomId, setRoomId] = useState("");
  const [expansion, setExpansion] = useState(1);

  function handleNewGame() {
    if (newGame || joinRoom) {
      return;
    }

    setUsername("");
    setRoomId("");
    setNewGame(true);
  }

  function handleJoinRoom() {
    if (newGame || joinRoom) {
      return;
    }

    setUsername("");
    setRoomId("");
    setJoinRoom(true);
  }

  function handleSubmitNewGame() {
    if (username === "") {
      setError("Username Invalid");
      return;
    }

    startGameCallback(username, "", expansion);
  }

  function handleSubmitJoinRoom() {
    if (username === "") {
      setError("Username Invalid");
      return;
    } else if (roomId.length !== 5) {
      setError("Room Code Must Be 5 Characters");
      return;
    }

    fetch(Constants.ApiBaseUrl + "CheckRoomId/" + roomId.toLowerCase(), {
      method: 'GET'
    })
    .then((response) => response.json())
    .then((responseJson) => {
      if (!responseJson.response) {
        setError("Invalid Room Code");
      } else {
        startGameCallback(username, roomId.toLowerCase(), -1);
      }
    });
  }

  function handleBack() {
    setNewGame(false);
    setJoinRoom(false);
  }

  function renderMenuScreen() {
    if (!newGame && !joinRoom) {
      return (
        <div className="StartScreen-Bottom">
          <div className="StartScreen-Title-Button" onClick={() => handleNewGame()}>
            <h2 className="StartScreen-Title-Button-Text noselect">New Game</h2>
          </div>
          <div className="StartScreen-Title-Button" onClick={() => handleJoinRoom()}>
            <h2 className="StartScreen-Title-Button-Text noselect">Join Room</h2>
          </div>
        </div>
      )
    } else if (newGame) {
      return (
        <div className="StartScreen-Bottom">
          <div className="StartScreen-Bottom-Back-Block">
            <div className="StartScreen-Back" onClick={() => handleBack()}></div>
          </div>
          <div>
            <h2 className="StartScreen-Expansion-Title">Expansion:</h2>
            <div className="StartScreen-Expansion-Numbers">
              <div className={expansion === 1 ? "StartScreen-Expansion-Numbers-Selected" : ""} onClick={() => setExpansion(1)}>
                <h2>1</h2>
              </div>
              <div className={expansion === 2 ? "StartScreen-Expansion-Numbers-Selected" : ""} onClick={() => setExpansion(2)}>
                <h2>2</h2>
              </div>
              <div className={expansion === 3 ? "StartScreen-Expansion-Numbers-Selected" : ""} onClick={() => setExpansion(3)}>
                <h2>3</h2>
              </div>
              <div className={expansion === 4 ? "StartScreen-Expansion-Numbers-Selected" : ""} onClick={() => setExpansion(4)}>
                <h2>4</h2>
              </div>
            </div>
          </div>
          <input type="text" onChange={e => setUsername(e.target.value)} value={username} className="StartScreen-Input" placeholder="Enter username"/>
          <div onClick={() => handleSubmitNewGame()} className="button">
            <h3 className="StartScreen-Submit-Text noselect">Submit</h3>
          </div>
          {(error) &&
            <p className="errorText">{error}</p>
          }
        </div>
      )
    } else {
      return (
        <div className="StartScreen-Bottom">
          <div className="StartScreen-Bottom-Back-Block">
            <div className="StartScreen-Back" onClick={() => handleBack()}></div>
          </div>
          <input type="text" onChange={e => setUsername(e.target.value)} value={username} className="StartScreen-Input" placeholder="Enter username"/>
          <input type="text" onChange={e => setRoomId(e.target.value)} value={roomId} className="StartScreen-Input" placeholder="Enter room id"/>
          <div onClick={() => handleSubmitJoinRoom()} className="button">
            <h3 className="noselect">Submit</h3>
          </div>
          {(error) &&
            <p className="StartScreen-ErrorText">{error}</p>
          }
        </div>
      )
    }
  }

  return (
    <div className="StartScreen">
      <div className="StartScreen-TitleBlock">
        <img src={ArchitortureLogo} className="StartScreen-Title" alt={"Logo"}/>
      </div>
      {renderMenuScreen()}
    </div>
  );
}

export default StartScreen;
