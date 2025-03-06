import React, { useState, useEffect } from "react";

import '../css/LobbyScreen.css';
import '../css/Shared.css';

function LobbyScreen({parentGameInfo, handleLobbyCallback}) {
  const [gameInfo, setgameInfo] = useState(parentGameInfo);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    setgameInfo(parentGameInfo);
  }, [parentGameInfo]);

  function toggleReady() {
    handleLobbyCallback(!isReady);
    setIsReady(!isReady);
  }

  return (
    <div className="LobbyScreen">
      <div className="LobbyScreen-Top">
        <h1 className="LobbyScreen-ExpansionText">Expansion: {Math.max(...gameInfo.expansions)}</h1>
        <div className="LobbyScreen-Players">
          {gameInfo.players.map(function(p, idx) {
            return (
              <div key={idx} className="LobbyScreen-Player">
                <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                <h3 className={p.isReady ? "LobbyScreen-Player-Ready" : "LobbyScreen-Player-NotReady"}>{p.username}</h3>
              </div>
            );
          })}
        </div>
      </div>
      <div className="LobbyScreen-Bottom">
        <div className="LobbyScreen-Bottom-Top">
          <h2 className="LobbyScreen-Bottom-RoomCode">Room Code: {gameInfo.roomId.toUpperCase()}</h2>
          <div className="LobbyScreen-Bottom-Top-Right">
            <span>
              <div className="button" onClick={() => toggleReady()}>
                <h3 className="noselect">Ready</h3>
              </div>
            </span>
          </div>
        </div>
        <div>
        </div>
      </div>
    </div>
  );
}

export default LobbyScreen;
