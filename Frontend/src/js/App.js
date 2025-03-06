import React, { useState, useEffect } from "react";

import '../css/App.css';

import StartScreen from './StartScreen';
import LobbyScreen from './LobbyScreen';
import GameScreen from './GameScreen';
import RulesScreen from './RulesScreen';
import * as Enums from './constants/Enums';
import * as Constants from './constants/Constants';

function App() {
  const [gameState, setGameState] = useState(Enums.GameStates.MainMenu)
  const [webSocket, setWebSocket] = useState(null);
  const [gameInfo, setGameInfo] = useState({});
  const [cardAction, setCardAction] = useState(null);
  const [cards, setCards] = useState(null);
  const [showRules, setShowRules] = useState(false);
  const [architorturePlayerId, setArchitorturePlayerId] = useState(null);
  const [undoObject, setUndoObject] = useState(null);
  const [architectureMemoryLossDrawn, setArchitectureMemoryLossDrawn] = useState(false);

  useEffect(() => {
    const cleanup = () => {
      if (webSocket !== null) {
        webSocket.close();
      }
    };

    window.addEventListener('beforeunload', cleanup);
    return () => {
      window.removeEventListener('beforeunload', cleanup);
    };
  }, [webSocket]);

  function getMessageBody(requestType, data) {
    return {
      requestType: requestType,
      body: data,
    };
  }

  function handleStartGameCallback(username, roomId, expansion) {
    fetch(Constants.ApiBaseUrl + "cards/" + expansion, {
      method: 'GET'
    })
    .then((response) => response.json())
    .then((responseJson) => {
      setCards(responseJson);
    });

    let ws = new WebSocket(Constants.WsBaseUrl+"?userName="+username+"&roomId="+roomId+"&expansion="+expansion);
    setWebSocket(ws);
    ws.onopen= () => {
    };

    ws.onclose = e => {
    };

    ws.onerror = err => {
      console.error("Socket encountered error: ", err.message, "Closing Socket");
      ws.close();
    };

    ws.onmessage = msg => {
      let data = JSON.parse(msg.data);
      if (data.messageType === Enums.MessageTypes.GameInfoUpdate ||
        data.messageType === Enums.MessageTypes.TurnStart ||
        data.messageType === Enums.MessageTypes.LobbyStart ||
        data.messageType === Enums.MessageTypes.GameStart ||
        data.messageType === Enums.MessageTypes.EndLoading) {
        setGameInfo(data);
      } else if (data.messageType === Enums.MessageTypes.CardAction) {
        setCardAction(data);
      } else if (data.messageType === Enums.MessageTypes.Architorture) {
        if (!data.player) {
          setArchitorturePlayerId(null);
        } else {
          setArchitorturePlayerId(data.player.id);
        }
      } else if (data.messageType === Enums.MessageTypes.Undo) {
        setUndoObject(data);
      } else if (data.messageType === Enums.MessageTypes.ArchitectureMemoryLoss) {
        setArchitectureMemoryLossDrawn(true);
      }

      if (data.messageType === Enums.MessageTypes.LobbyStart) {
        setGameState(Enums.GameStates.Lobby);
      } else if (data.messageType === Enums.MessageTypes.GameStart) {
        setGameState(Enums.GameStates.GameActive);
      }
    };
  }

  function handleLobbyReady(isReady) {
    let readyMessage = {
      IsReady: isReady,
    };
    webSocket.send(JSON.stringify(getMessageBody(Enums.RequestTypes.Ready, JSON.stringify(readyMessage))));
  }

  function messagePassthrough(message, messageType) {
    webSocket.send(JSON.stringify(getMessageBody(messageType, JSON.stringify(message))));
  }

  function clearCardAction() {
    setCardAction(null);
  }

  function clearUndoObject() {
    setUndoObject(null);
  }

  function resetArchitectureMemoryLossDrawn() {
    setArchitectureMemoryLossDrawn(false);
  }

  function renderGame() {
    if (gameState === Enums.GameStates.MainMenu) {
      return (
        <StartScreen startGameCallback={(username, roomId, expansion) => handleStartGameCallback(username, roomId, expansion)}/>
      );
    } else if (gameState === Enums.GameStates.Lobby) {
      return (
        <LobbyScreen parentGameInfo={gameInfo} handleLobbyCallback={(isReady) => handleLobbyReady(isReady)} />
      );
    } else if (gameState === Enums.GameStates.GameActive) {
      return (
        <GameScreen parentGameInfo={gameInfo} 
          messagePassthroughCallback={(message, messageType) => messagePassthrough(message, messageType)} 
          parentCardAction={cardAction} 
          clearCardActionCallback={() => clearCardAction()} 
          allCards={cards}
          architorturePlayerId={architorturePlayerId}
          undoObject={undoObject}
          clearUndoObject={() => clearUndoObject()}
          architectureMemoryLossDrawn={architectureMemoryLossDrawn}
          resetArchitectureMemoryLossDrawn={() => resetArchitectureMemoryLossDrawn()}/>
      );
    } else {
      return (
        <div>
          An Error Occurred
        </div>
      );
    }
  }

  return (
    <div>
      {(gameState !== Enums.GameStates.MainMenu && cards && !showRules) &&
        <div className="RulesMenuButton" onClick={() => setShowRules(true)}>
          <p>Show Rules</p>
        </div>
      }
      {(showRules) &&
        <RulesScreen closeCallback={() => setShowRules(false)} cards={cards}/>
      }
      {renderGame()}
    </div>
  );
}

export default App;
