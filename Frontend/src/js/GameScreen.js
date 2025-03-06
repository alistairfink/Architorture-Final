import React, { useState, useEffect } from "react";
import * as Enums from './constants/Enums';

import '../css/GameScreen.css';
import '../css/Shared.css';
import CardBack from '../res/cards/CardBack.png';
import CardBlank from '../res/cards/CardBlank.png';
import Arrow from '../res/arrow.png';
import USB from '../res/USB.png';
import MemExpansion from '../res/Memory-Expansion.png';
import Loading from '../res/loading.svg';
import Trash from '../res/trash.svg';
import Archive from '../res/archive.svg';
import PlayerEliminated from '../res/PlayerEliminated.png';
import PlayerWon from '../res/PlayerWon.png';

function GameScreen({
  parentGameInfo, 
  messagePassthroughCallback, 
  parentCardAction, 
  clearCardActionCallback, 
  allCards, 
  architorturePlayerId, 
  undoObject, 
  clearUndoObject,
  architectureMemoryLossDrawn,
  resetArchitectureMemoryLossDrawn
}) {

  const [gameInfo, setgameInfo] = useState(parentGameInfo);
  const [archiveFlowObject, setArchiveFlowObject] = useState(null);
  const [cardAction, setCardAction] = useState(null);
  const [vendettaFlowObject, setVendettaFlowObject] = useState(null);
  const [foundUSBFlowObject, setFoundUSBFlowObject] = useState(null);
  const [sharePreviewFlowObject, setSharePreviewFlowObject] = useState(null);
  const [knowledgeFlowObject, setKnowledgeFlowObject] = useState(null);
  const [assistanceFlowObject, setAssistanceFlowObject] = useState(null);
  const [hoverCard, setHoverCard] = useState(null);
  const [thisPlayer, setThisPlayer] = useState(null);
  const [numbCardSelected, setNumbCardSelected] = useState(null);
  const [loading, setLoading] = useState(false);
  const [shownEliminated, setShownEliminated] = useState(false);
  const [eliminated, setEliminated] = useState(false);
  const [winner, setWinner] = useState(false);

  useEffect(() => {
    let player = null;
    let elimCount = 0;
    for (let i = 0; i < parentGameInfo.players.length; i++) {
      if (parentGameInfo.players[i].id === parentGameInfo.playerId) {
        if (parentGameInfo.players[i].eliminated && !shownEliminated) {
          setEliminated(true);
          setShownEliminated(true);
        }

        setThisPlayer(parentGameInfo.players[i]);
        player = parentGameInfo.players[i];
      }

      if (parentGameInfo.players[i].eliminated) {
        elimCount++;
      }
    }

    if (elimCount === parentGameInfo.players.length-1 && !player.eliminated) {
      setWinner(true);
    }

    if (parentGameInfo.messageType === Enums.MessageTypes.EndLoading) {
      setLoading(false);
    }

    setCardAction(parentCardAction);
    setgameInfo(parentGameInfo);
    setAssistanceFlowObject(null);
  }, [parentGameInfo, parentCardAction, clearCardActionCallback, setEliminated, setWinner, shownEliminated, setShownEliminated]);

  function actionInvalid() {
    return (cardAction !== null && cardAction.requiresResponse) || 
      vendettaFlowObject !== null || 
      foundUSBFlowObject !== null || 
      sharePreviewFlowObject !== null || 
      knowledgeFlowObject !== null ||
      assistanceFlowObject !== null;
  }

  function handleCardClick(index) {
    if (actionInvalid()) {
      return;
    }

    let card = gameInfo.hand[index];
    if (gameInfo.currentPlayer !== gameInfo.playerId &&
      card.id !== 27 &&
      card.id !== 28 &&
      card.id !== 34) {
      return;
    }

    switch(card.dbId) {
      case 1:
      case 2:
      case 3:
      case 4:
      case 27:
      case 28:
      case 37:
      case 35:
        return;
      case 19:
        setVendettaFlowObject({
          UseTarget: false,
          Target: "00000000-0000-0000-0000-000000000000",
          CardToPick: "00000000-0000-0000-0000-000000000000",
          playedCard: card.id, 
        });
        break;
      case 31:
      case 32:
        setVendettaFlowObject({
          UseTarget: true,
          Target: "00000000-0000-0000-0000-000000000000",
          CardToPick: "00000000-0000-0000-0000-000000000000",
          playedCard: card.id,
        });
        break;
      case 17:
      case 18:
        setFoundUSBFlowObject({
          playedCard: card.id,
        });
        break;
      case 36:
        setSharePreviewFlowObject({
          playedCard: card.id,
        });
        break;
      case 38:
        setAssistanceFlowObject({
          PlayedCard: [card.id],
        });
        break;
      case 29:
      case 30:
      case 39:
      case 40:
      case 41:
      case 42:
      case 43:
      case 44:
      case 45:
      case 46:
      case 47:
      case 48:
      case 49:
      case 50:
      case 51:
      case 52:
      case 53:
      case 54:
      case 55:
      case 56:
      case 57:
        let allCardsMap = {};
        let allCardsArr = [];
        for (let i = 0; i < allCards.length; i++) {
          if (!allCardsMap[allCards[i].name]) {
            allCardsArr.push(allCards[i])
          }

          allCardsMap[allCards[i].name] = true;
        }

        let flowObject = {
          OriginalCardDbId: card.dbId,
          Target: "00000000-0000-0000-0000-000000000000",
          Cards: {},
          CardToRequest: "",
          CurrentStep: Enums.KnowledgeSteps.SelectPlayerCards,
          AllCards: allCardsArr,
        };
        flowObject.Cards[card.id] = card;
        setKnowledgeFlowObject(flowObject);
        break;
      default:
        let playCardObject = {
          Cards: [card.id],
          TargetPlayerId: "00000000-0000-0000-0000-000000000000",
        } 
        messagePassthroughCallback(playCardObject, Enums.RequestTypes.PlayCard);
        break;
    }
  }

  function openArchive() {
    setArchiveFlowObject({
      CardToArchive: "00000000-0000-0000-0000-000000000000",
      CardToUnarchive: "00000000-0000-0000-0000-000000000000",
      CurrentStep: Enums.ArchiveSteps.ViewArchive
    });
  }

  function closeArchive() {
    setArchiveFlowObject(null);
  }

  function handleArchive(card) {
    messagePassthroughCallback({
      archiveCardId: card.id,
      unarchiveCardId: "00000000-0000-0000-0000-000000000000",
    }, Enums.RequestTypes.Archive);
    setArchiveFlowObject(null);
  }

  function handleArchiveCard() {
    messagePassthroughCallback({
      archiveCardId: archiveFlowObject.CardToArchive,
      unarchiveCardId: archiveFlowObject.CardToUnarchive,
    }, Enums.RequestTypes.Archive);
    setArchiveFlowObject(null);
  }

  function handleDrawCard() {
    if (actionInvalid()) {
      return;
    }
    
    messagePassthroughCallback({}, Enums.RequestTypes.DrawCard);
  }

  function HandleCardAction() {
    switch(cardAction.cardDbId) {
      case 20:
      case 21:
      case 36:
        return (
          <div>
            <div className="GameScreen-Back">
              <div onClick={() => clearCardActionCallback()}></div>
            </div>
            {(cardAction.cardDbId === 36 && gameInfo.currentPlayer !== gameInfo.playerId) &&
              <h1>A Player has Shared a Preview</h1>
            }
            <h1>Preview of Next {cardAction.cards.length} Cards:</h1>
            <div className="GameScreen-CardActionCards">
              {cardAction.cards.map(function(c, i) {
                return (
                  <div key={i}>
                    <h2>{i+1}</h2>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        );
      case 25:
      case 26:
        if (cardAction.cards[0].order === null) {
          let cAction = cardAction;
          for (let i = 0; i < cardAction.cards.length; i++) {
            cAction.cards[i].order = i+1;
          }
          setCardAction({...cAction});
        }

        return (
          <div>
            <h1>Select Order of Next {cardAction.cards.length} Cards</h1>
            {(cardAction.Error !== "" && cardAction.Error !== null) &&
              <p className="errorText">{cardAction.Error}</p>
            }
            <div className="GameScreen-CardActionCards">
              {cardAction.cards.map(function(c, i) {
                return (
                  <div className="GameScreen-CardActionCard-Col" key={i}>
                    <input type="number" onChange={e => {
                      let cAction = cardAction;
                      cAction.cards[i].order = e.target.value;
                      setCardAction({...cAction});
                    }} value={cardAction.cards[i].order} min="1" max={cardAction.cards.length}/>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
            <div onClick={() => handleChangeOfPlansResponse()} className="button">
              <h3 className="StartScreen-Submit-Text noselect">Submit</h3>
            </div>
          </div>
        );
      case 5:
      case 6:
        return (
          <div>
            <h1>You've Drawn and Architecture Dump Card</h1>
            <h1>Select a Card to Discard</h1>
            <div className="GameScreen-CardActionCards">
              {gameInfo.hand.map(function(c, i) {
                return (
                  <div onClick={() => handleArchitectureDumpSelection(i)} className="GameScreen-CardActionCardSelectable" key={i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        );
      case 38:
        return (
          <div>
            <h1>A Player has "Asked" You to Assist Them</h1>
            <h1>Select a Card to Give to Them</h1>
            <div className="GameScreen-CardActionCards">
              {gameInfo.hand.map(function(c, i) {
                return (
                  <div onClick={() => handleAssitanceSelection(c.id)} className="GameScreen-CardActionCardSelectable" key={i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        );
      case 7:
      case 8:
        let tempArr = [];
        let handSize = 0;
        for (let i = 0; i < gameInfo.players.length; i++) {
          if (gameInfo.players[i].id === gameInfo.currentPlayer) {
            handSize = gameInfo.players[i].cardCount;
            break;
          }
        }

        for (let i = 0; i < handSize; i++) {
          tempArr.push("");
        } 

        return (
          <div className="GameScreen-ArchitectureMemoryLoss">
            <h1>The Current Player has Drawn an Architecture Memory Loss Card</h1>
            <h1>Select a Card From their Hand to Discard</h1>
            <div className="GameScreen-CardActionCards">
              {tempArr.map(function(c, i) {
                return (
                  <div key={"CardActionCard-" + i} onClick={() => {
                    submitArchitectureMemoryLoss(i);
                  }} className="GameScreen-CardActionCardSelectable">
                    <img className="GameScreen-CardActionCard"  src={CardBack} alt={"All_Cards_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        );
      case 3:
      case 4:
        let architortureCount = -1;
        for (let i = 0; i < gameInfo.hand.length; i++) {
          if (gameInfo.hand[i].cardTypeId === 1) {
            architortureCount++;
          }
        }

        if (numbCardSelected === null) {
          return (
            <div>
              <h1>Select a Save Card to Use</h1>
              <div className="GameScreen-CardActionCards">
                {gameInfo.hand.map(function(c, i) {
                  if (c.cardTypeId !== 2) {
                    return (null);
                  }

                  if ((c.dbId === 35 || c.dbId === 49) && architortureCount > 0) {
                      architortureCount--;
                      return (null);
                  }

                  return (
                    <div onClick={() => handleSaveCardSelect(c)} className="GameScreen-CardActionCardSelectable" key={i}>
                      <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        }

        return (
          <div>
            <h1>Select a Target</h1>
            <div className="GameScreen-SharePreviewUsers">
              {gameInfo.players.map(function(p, i) {
                if (p.id !== gameInfo.playerId) {
                  return (
                    <div key={"SharePreviewPlayer-"+i} onClick={() => handleNumbCard(p.id)}>
                      <div className="GameScreen-Player GameScreen-SharePreviewUser">
                        <h3>{p.username}</h3>
                        <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                      </div>
                    </div>
                  );
                } else {
                  return (null);
                }
              })}
            </div>
          </div>
        );
      default:
        return (null);
    }
  }

  function HandleSharePreviewUI() {
    return (
      <div>
        <div className="GameScreen-Back-Normal">
          <div onClick={() => {
            setSharePreviewFlowObject(null);
          }}></div>
        </div>
        <h1>Select a Player to Share the Preview with:</h1>
        <div className="GameScreen-SharePreviewUsers">
          {gameInfo.players.map(function(p, i) {
            if (p.id !== gameInfo.playerId) {
              return (
                <div key={"SharePreviewPlayer-"+i} onClick={() => handleSharePreviewSelection(i)}>
                  <div className="GameScreen-Player GameScreen-SharePreviewUser">
                    <h3>{p.username}</h3>
                    <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                  </div>
                </div>
              );
            } else {
              return (null);
            }
          })}
        </div>
      </div>
    );
  }

  function HandleAssistanceUI() {
    return (
      <div>
        <div className="GameScreen-Back-Normal">
          <div onClick={() => {
            setAssistanceFlowObject(null);
          }}></div>
        </div>
        <h1>Select a Player to Assist You</h1>
        <div className="GameScreen-SharePreviewUsers">
          {gameInfo.players.map(function(p, i) {
            if (p.id !== gameInfo.playerId) {
              return (
                <div key={"SharePreviewPlayer-"+i} onClick={() => submiteAsistanceFlow(p.id)}>
                  <div className="GameScreen-Player GameScreen-SharePreviewUser">
                    <h3>{p.username}</h3>
                    <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                  </div>
                </div>
              );
            } else {
              return (null);
            }
          })}
        </div>
      </div>
    );
  }

  function HandleVendettaUI() {
    return (
      <div> 
        {vendettaFlowObject.CardToPick === "00000000-0000-0000-0000-000000000000" &&
          <div className="GameScreen-CardActionCards">
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                setVendettaFlowObject(null);
              }}></div>
            </div>
            <h1>Select A Card to Give</h1>
            <div className="GameScreen-CardActionCards">
              {gameInfo.hand.map(function(c, i) {
                if (c.id === vendettaFlowObject.playedCard) {
                  return(null);
                }

                return (
                  <div onClick={() => {
                        let data = vendettaFlowObject;
                        data.CardToPick = c.id;
                        setVendettaFlowObject({...data});
                        if (!data.UseTarget) {
                          handleVendettaTargetSelection();
                        }
                      }} className="GameScreen-CardActionCardSelectable" key={i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        }
        {(vendettaFlowObject.CardToPick !== "00000000-0000-0000-0000-000000000000" && vendettaFlowObject.UseTarget) &&
          <div className="GameScreen-Vendetta-SelectTarget">
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                  let data = vendettaFlowObject;
                  data.CardToPick = "00000000-0000-0000-0000-000000000000";
                  setVendettaFlowObject({...data});
              }}></div>
            </div>
            <h1>Select a Player to Target</h1>
            <div className="GameScreen-ShowUsers-Vendetta">
              {gameInfo.players.map(function(p, i) {
                if (p.id !== gameInfo.playerId) {
                  return (
                    <div key={"SharePreviewPlayer-"+i} onClick={() => {
                      let data = vendettaFlowObject;
                      data.Target = p.id;
                      setVendettaFlowObject({...data});
                      handleVendettaTargetSelection();
                    }}>
                      <div className="GameScreen-Player GameScreen-SharePreviewUser">
                        <h3>{p.username}</h3>
                        <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                      </div>
                    </div>
                  );
                } else {
                  return (null);
                }
              })}
            </div>
          </div>
        }
      </div>
    );
  }

  function HandleFoundUSBUI() {
    return (
      <div>
        <div className="GameScreen-Back-Normal">
          <div onClick={() => {
            setFoundUSBFlowObject(null);
          }}></div>
        </div>
        <h1>Select Player to Steal USB From</h1>
        <div className="GameScreen-SharePreviewUsers">
          {gameInfo.players.map(function(p, i) {
            if (p.id !== gameInfo.playerId && p.archiveIncreases > 0) {
              return (
                <div key={"SharePreviewPlayer-"+i} onClick={() => handleFoundUSBTargetSelection(i)}>
                  <div className="GameScreen-Player GameScreen-SharePreviewUser">
                    <h3>{p.username}</h3>
                    <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                  </div>
                </div>
              );
            } else {
              return (null);
            }
          })}
        </div>
      </div>
    );
  }

  function submitKnowledgeFlow() {
    if (Object.keys(knowledgeFlowObject.Cards).length === 3) {
      let cardIds = [];
      for (const [key] of Object.entries(knowledgeFlowObject.Cards)) {
        cardIds.push(key);
      }

      let dataObject = {
        Cards: cardIds,
        TargetPlayerId: knowledgeFlowObject.Target,
        cardrequestName: knowledgeFlowObject.CardToRequest,
      };

      setKnowledgeFlowObject(null);      
      messagePassthroughCallback(dataObject, Enums.RequestTypes.PlayCard);
    } else if (Object.keys(knowledgeFlowObject.Cards).length === 2) { 
      let cardIds = [];
      for (const [key] of Object.entries(knowledgeFlowObject.Cards)) {
        cardIds.push(key);
      }

      let dataObject = {
        cards: cardIds,
        targetPlayerId: knowledgeFlowObject.Target,
        targetCardIndex: knowledgeFlowObject.TargetCardIndex,
      };

      setKnowledgeFlowObject(null);      
      messagePassthroughCallback(dataObject, Enums.RequestTypes.PlayCard);
      setLoading(true);
    }
  }

  function knowledgeCardIsValid(originalCardDbId, newCardDbId) {
    if (newCardDbId === 29 || newCardDbId === 30) {
      return true;
    }

    switch(originalCardDbId) {
      case 39:
      case 40:
      case 41:
      case 42:
        return newCardDbId === 39 || newCardDbId === 40 || newCardDbId === 41 || newCardDbId === 42;
      case 43:
      case 44:
      case 45:
        return newCardDbId === 43 || newCardDbId === 44 || newCardDbId === 45;
      case 46:
      case 47:
      case 48:
        return newCardDbId === 46 || newCardDbId === 47 || newCardDbId === 48;
      case 49:
      case 50:
      case 51:
        return newCardDbId === 49 || newCardDbId === 50 || newCardDbId === 51;
      case 52:
      case 53:
      case 54:
        return newCardDbId === 52 || newCardDbId === 53 || newCardDbId === 54;
      case 55:
      case 56:
      case 57:
        return newCardDbId === 55 || newCardDbId === 56 || newCardDbId === 57;
      default:
        return false;
    }
  }

  function HandleKnowledgeFlowUI() {
    let topCount = 0;
    for (let i = 0; i < gameInfo.players.length; i++) {
      if (gameInfo.players[i].cardCount > topCount) {
        topCount = gameInfo.players[i].cardCount;
      }
    }

    let emptyArr = [];
    for (let i = 0; i < topCount; i++) {
      emptyArr.push(i);
    }

    return (
      <div>
        {(knowledgeFlowObject.CurrentStep === Enums.KnowledgeSteps.SelectPlayerCards) &&
          <div className="GameScreen-ArchiveFlow-Outer">
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                setKnowledgeFlowObject(null);
              }}></div>
            </div>
            <div className="GameScreen-Archive-ArchiveButton">
              <div onClick={() => {
                if (Object.keys(knowledgeFlowObject.Cards).length !== 2 && Object.keys(knowledgeFlowObject.Cards).length !== 3) {
                  return;
                }

                let data = knowledgeFlowObject;
                data.CurrentStep = Enums.KnowledgeSteps.SelectTarget;
                setKnowledgeFlowObject({...data});
              }} className="button">
                <h3 className="StartScreen-Submit-Text noselect">Submit</h3>
              </div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select Cards to Use</h1>
              <div className="GameScreen-CardActionCards">
                {gameInfo.hand.map(function(c, i) {
                  if (knowledgeCardIsValid(knowledgeFlowObject.OriginalCardDbId, c.dbId)) {
                    return (
                      <div key={"KnowledgeCard-"+i} onClick={() => {
                        if (Object.keys(knowledgeFlowObject.Cards).length === 3 && !knowledgeFlowObject.Cards[c.id]) {
                          return;
                        }

                        let data = knowledgeFlowObject;
                        if (data.Cards === null) {
                          data.Cards = {};
                        }

                        if (data.Cards[c.id]) {
                          delete data.Cards[c.id];
                        } else {
                          data.Cards[c.id] = c;
                        }

                        setKnowledgeFlowObject({...data});
                      }} className={knowledgeFlowObject.Cards[c.id] ? "GameScreen-KnowledgeCardSelected" : "GameScreen-CardActionCardSelectable"}>
                        <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                      </div>
                    );
                  }                    
                  return(null);
                })}
              </div>
            </div>
          </div>
        }
        {(knowledgeFlowObject.CurrentStep === Enums.KnowledgeSteps.SelectTarget) &&
          <div>
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                let data = knowledgeFlowObject;
                data.Cards = {};
                data.CurrentStep = Enums.KnowledgeSteps.SelectPlayerCards;
                setKnowledgeFlowObject({...data});
              }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select Player to Target</h1>
              <div className="GameScreen-SharePreviewUsers">
                {gameInfo.players.map(function(p, i) {
                  if (p.id !== gameInfo.playerId) {
                    return (
                      <div key={"SharePreviewPlayer-"+i} onClick={() => {
                        let data = knowledgeFlowObject;
                        data.Target = p.id;
                        data.CurrentStep = Enums.KnowledgeSteps.SelectTargetCard;
                        if (Object.keys(knowledgeFlowObject.Cards).length === 2) {
                          data.CurrentStep = Enums.KnowledgeSteps.SelectTargetcardBlind;
                          data.Count = p.cardCount;
                          // submitKnowledgeFlow();
                        }
                        setKnowledgeFlowObject({...data});
                      }}>
                        <div className="GameScreen-Player GameScreen-SharePreviewUser">
                          <h3>{p.username}</h3>
                          <img src={require("../res/players/player_"+p.playerNumber+".png").default} alt={"player" + p.playerNumber}/>
                        </div>
                      </div>
                    );
                  } else {
                    return (null);
                  }
                })}
             </div>
            </div>
          </div>
        }
        {(knowledgeFlowObject.CurrentStep === Enums.KnowledgeSteps.SelectTargetCard) &&
          <div>
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                let data = knowledgeFlowObject;
                data.Target = "00000000-0000-0000-0000-000000000000";
                data.CurrentStep = Enums.KnowledgeSteps.SelectTarget;
                setKnowledgeFlowObject({...data});
              }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select a Card Type to Request</h1>
              <div className="GameScreen-CardActionCards">
                {knowledgeFlowObject.AllCards.map(function(c, i) {
                  return (
                    <div onClick={() => {
                      let data = knowledgeFlowObject;
                      data.CardToRequest = c.name;
                      setKnowledgeFlowObject({...data});
                      submitKnowledgeFlow();
                    }} className="GameScreen-CardActionCardSelectable" key={i}>
                      <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"All_Cards_"+i}/>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        }
        {(knowledgeFlowObject.CurrentStep === Enums.KnowledgeSteps.SelectTargetcardBlind) &&
          <div>
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                let data = knowledgeFlowObject;
                data.Target = "00000000-0000-0000-0000-000000000000";
                data.CurrentStep = Enums.KnowledgeSteps.SelectTarget;
                setKnowledgeFlowObject({...data});
              }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select a Card to Request</h1>
              <div className="GameScreen-CardActionCards">
                {emptyArr.map(function(c, i) {
                  if (i >= knowledgeFlowObject.Count) {
                    return (null);
                  }

                  return (
                    <div onClick={() => submitKnowledgeFlow()} className="GameScreen-CardActionCardSelectable" key={i}>
                      <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={CardBack} alt={"All_Cards_"+i}/>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        }
      </div>
    );
  }

  function ArchiveUI() {
    return (
      <div>
        {(archiveFlowObject.CurrentStep === Enums.ArchiveSteps.ViewArchive) && 
          <div className="GameScreen-ArchiveFlow-Outer">
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
                closeArchive();
              }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select A Card to Unarchive</h1>
              <h3>{gameInfo.archive.length} / {gameInfo.archiveMax}</h3>
              <div className="GameScreen-CardActionCards">
                {gameInfo.archive.map(function(c, i) {
                  return (
                    <div onClick={() => {
                      let data = archiveFlowObject;
                      data.CardToArchive = "00000000-0000-0000-0000-000000000000";
                      data.CardToUnarchive = c.id;
                      data.CurrentStep = Enums.ArchiveSteps.Archive;
                      if (gameInfo.hand.length < gameInfo.handMax) {
                        data.CurrentStep = Enums.ArchiveSteps.ChooseSwap;
                      }
                      setArchiveFlowObject({...data});
                    }} className="GameScreen-CardActionCardSelectable" key={"CardActionCard-" + i}>
                      <img className="GameScreen-CardActionCard" src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Archive_"+i}/>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        }
        {(archiveFlowObject.CurrentStep === Enums.ArchiveSteps.ChooseSwap) &&
          <div>
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
              let data = archiveFlowObject;
              data.CardToUnarchive = "00000000-0000-0000-0000-000000000000";
              data.CurrentStep = Enums.ArchiveSteps.ViewArchive;
              setArchiveFlowObject({...data});
            }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Swap with Card in Hand?</h1>
              <div className="button" onClick={() => {
                  let data = archiveFlowObject;
                  data.CurrentStep = Enums.ArchiveSteps.Archive;
                  setArchiveFlowObject({...data});
              }}>
                <h3>Yes</h3>
              </div>
              <div className="button" onClick={() => handleArchiveCard()}>
                <h3>No</h3>
              </div>
            </div>
          </div>
        }
        {(archiveFlowObject.CurrentStep === Enums.ArchiveSteps.Archive) && 
          <div>
            <div className="GameScreen-Back-Normal">
              <div onClick={() => {
              let data = archiveFlowObject;
              data.CardToUnarchive = "00000000-0000-0000-0000-000000000000";
              data.CurrentStep = Enums.ArchiveSteps.ViewArchive;
              setArchiveFlowObject({...data});
            }}></div>
            </div>
            <div className="GameScreen-ArchiveFlow">
              <h1>Select A Card to Archive</h1>
              <div className="GameScreen-CardActionCards">
                {gameInfo.hand.map(function(c, i) {
                  return (
                    <div onClick={() => {
                      let data = archiveFlowObject;
                      data.CardToArchive = c.id;
                      setArchiveFlowObject({...data});
                      handleArchiveCard();
                    }} className="GameScreen-CardActionCardSelectable">
                      <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        }
      </div>
    );
  }

  function UndoUI() {
    return (
      <div className="GameScreen-ArchiveFlow">
        <div className="GameScreen-Back-Normal">
          <div onClick={() => {
            let response = {
              cardId: "00000000-0000-0000-0000-000000000000",
              use: false,
            };
            messagePassthroughCallback(response, Enums.RequestTypes.Undo);
            clearUndoObject();
          }}></div>
        </div>
        {(undoObject.stage === Enums.UndoStates.Initial) &&
          <div>
            <h1>A Player is Attacking You With the Following Card(s)</h1>
            <div className="GameScreen-CardActionCards">
              {undoObject.card.map(function(c, i) {
                return (
                  <div key={"Undo_Used_"+i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
            <h1>Select a Card to Counter With</h1>
            <div className="GameScreen-CardActionCards">
              {gameInfo.hand.map(function(c, i) {
                if (c.dbId !== 27 &&
                  c.dbId !== 28 &&
                  c.dbId !== 34) {
                  return (null);
                }

                return (
                  <div onClick={() => {
                    let response = {
                      cardId: c.id,
                      use: true,
                    };
                    messagePassthroughCallback(response, Enums.RequestTypes.Undo);
                    clearUndoObject();
                  }} className="GameScreen-CardActionCardSelectable" key={"Undo_Card_"+i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        }        
        {(undoObject.stage === Enums.UndoStates.CanUndoUndo) &&
          <div>
            <h1>A Player has Undone Your Action</h1>
            <h1>Select a Card to Counter With</h1>
            <div className="GameScreen-CardActionCards">
              {gameInfo.hand.map(function(c, i) {
                if (c.dbId !== 34) {
                  return (null);
                }

                return (
                  <div onClick={() => {
                    let response = {
                      cardId: c.id,
                      use: true,
                    };
                    messagePassthroughCallback(response, Enums.RequestTypes.Undo);
                    clearUndoObject();
                  }} className="GameScreen-CardActionCardSelectable" key={"Undo_Card_"+i}>
                    <img className="GameScreen-CardActionCard" key={"CardActionCard-" + i} src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={"Card_"+i}/>
                  </div>
                );
              })}
            </div>
          </div>
        }
      </div>
    );
  }

  function ArchitectureMemoryLossDrawnUI() {
    return (
      <div>
        <div className="GameScreen-Back-Normal">
          <div onClick={() => {
            resetArchitectureMemoryLossDrawn();
          }}></div>
        </div>
        <h1>You've drawn an Architecture Memory Loss Card!</h1>
        <h1>The next player is choosing a card from your hand to discard.</h1>
      </div>
    );
  }

  function handleChangeOfPlansResponse() {
    let cAction = cardAction;
    cAction.Error = "";
    setCardAction({...cAction});
    let cards = [-1, -1];
    for (let i = 0; i < 2; i++) {
      if (!cardAction.cards[i].order || cardAction.cards[i].order > 2) {
        let cAction = cardAction;
        cAction.Error = "Invalid Input";
        setCardAction({...cAction});
        return;
      }

      cards[cardAction.cards[i].order - 1] = cardAction.cards[i].id;
    }

    for (let i = 0; i < 2; i++) {
      if (cards[i] === -1) {
        let cAction = cardAction;
        cAction.Error = "Invalid Input";
        setCardAction({...cAction});
        return;
      }
    }

    let response = {
      actionCardId: cardAction.cardDbId,
      cards: cards, 
    }
    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function handleVendettaTargetSelection() {
    let playCardObject = {
      cards: [vendettaFlowObject.playedCard],  
      targetPlayerId: vendettaFlowObject.Target,
      cardsToGive: [vendettaFlowObject.CardToPick],
    };

    messagePassthroughCallback(playCardObject, Enums.RequestTypes.PlayCard);
    setVendettaFlowObject(null);
    setLoading(true);
  }

  function handleFoundUSBTargetSelection(idx) {
    let playCardObject = {
      cards: [foundUSBFlowObject.playedCard],
      targetPlayerId: gameInfo.players[idx].id,
    };

    messagePassthroughCallback(playCardObject, Enums.RequestTypes.PlayCard);
    setFoundUSBFlowObject(null);
    setLoading(true);
  }

  function handleSharePreviewSelection(idx) {
    let playCardObject = {
      cards: [sharePreviewFlowObject.playedCard],
      targetPlayerId: gameInfo.players[idx].id,
    };

    messagePassthroughCallback(playCardObject, Enums.RequestTypes.PlayCard);
    setSharePreviewFlowObject(null);
  } 

  function handleArchitectureDumpSelection(idx) {
    let response = {
      actionCardId: cardAction.cardDbId,
      cards: [gameInfo.hand[idx].id],
    }

    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function submiteAsistanceFlow(targetId) {
    let playCardObject = {
      Cards: assistanceFlowObject.PlayedCard,
      TargetPlayerId: targetId,
    };
    messagePassthroughCallback(playCardObject, Enums.RequestTypes.PlayCard);
    setAssistanceFlowObject(null);
    setLoading(true);
  }

  function handleAssitanceSelection(cardId) {    
    let response = {
      actionCardId: cardAction.cardDbId,
      cards: [cardId], 
    };
    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function handleSaveCardSelect(card) {
    if (card.dbId === 37) {
      setNumbCardSelected(card);
      return;
    }

    let response = {
      actionCardId: cardAction.cardDbId,
      cards: [card.id, cardAction.cards[0].id],
    };
    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function handleNumbCard(targetId) {
    let response = {
      actionCardId: cardAction.cardDbId,
      cards: [numbCardSelected.id, cardAction.cards[0].id],
      target: targetId
    };

    setNumbCardSelected(null);
    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function submitArchitectureMemoryLoss(i) {
    let response = {
      actionCardId: cardAction.cardDbId,
      cardIndex: i,
    };
    clearCardActionCallback();
    messagePassthroughCallback(response, Enums.RequestTypes.CardActionResponse);
  }

  function popupShouldBeOpen() {
    return cardAction !== null || 
      sharePreviewFlowObject !== null ||
      assistanceFlowObject !== null ||
      vendettaFlowObject !== null ||
      foundUSBFlowObject !== null ||
      knowledgeFlowObject !== null ||
      archiveFlowObject !== null ||
      undoObject !== null ||
      architectureMemoryLossDrawn;
  }

  function handleDiscard(card) {
    let playCardObject = {
      Cards: [card.id],
      TargetPlayerId: "00000000-0000-0000-0000-000000000000",
    };

    messagePassthroughCallback(playCardObject, Enums.RequestTypes.Discard);
  }

  return (
    <div>
      {(loading) &&
        <div className="GameScreen-Loading">
          <img src={Loading} alt="Loading"/>
          <h2>Loading</h2>
        </div>
      }
      {(eliminated) &&
        <div className="GameScreen-PlayerEliminated">
          <div className="GameScreen-Back-Eliminated">
            <div onClick={() => setEliminated(false)}></div>
          </div>
          <img src={PlayerEliminated} alt="Player Eliminated"/>
        </div>
      }
      {(winner) &&
        <div className="GameScreen-PlayerEliminated">
          <img src={PlayerWon} alt="Winner"/>
        </div>
      }
      {(popupShouldBeOpen()) &&
        <div className="GameScreen-Popup">
          <div>      
            {(cardAction !== null) && 
              <HandleCardAction />
            }
            {(sharePreviewFlowObject !== null) &&
              <HandleSharePreviewUI />
            }
            {(assistanceFlowObject !== null) &&
              <HandleAssistanceUI />
            }
            {(vendettaFlowObject !== null) &&
              <HandleVendettaUI />
            }
            {(foundUSBFlowObject !== null) &&
              <HandleFoundUSBUI />
            }
            {(knowledgeFlowObject !== null) &&
              <HandleKnowledgeFlowUI />
            }
            {(archiveFlowObject !== null) &&
              <ArchiveUI />
            }
            {(undoObject !== null) && 
              <UndoUI />
            }
            {(architectureMemoryLossDrawn) &&
              <ArchitectureMemoryLossDrawnUI />
            }
          </div>
        </div>
      }
      <div className="GameScreen">
        <div className="GameScreen-Players">
          {gameInfo.players.map(function(p, i) {
            if (p.id !== gameInfo.playerId) {
              return (
                <div key={i}>
                  {(p.id === gameInfo.currentPlayer) &&
                    <div className="GameScreen-TurnIndicator">
                      <img src={Arrow} alt={"Turn_Indicator"}/>
                    </div>
                  }
                  {(p.archiveIncreases > 0) &&
                    <div className="GameScreen-USBIndicator">
                      <img src={USB} alt={"Archive_Indicator"}/>
                    </div>
                  }
                  {(p.handIncreases > 0) &&
                    <div className="GameScreen-HandIndicator">
                      <img src={MemExpansion} alt={"Hand_Indicator"}/>
                    </div>
                  }
                  <div className="GameScreen-Player">
                    <h3>{p.username}</h3>
                    <img className={architorturePlayerId === p.id ? "GameScreen-ImgShake" : ""} src={require("../res/players/player_"+p.playerNumber+(p.eliminated ? "_eliminated" : "")+".png").default} alt={"player" + p.playerNumber}/>
                    <div className="GameScreen-Player-Cards">
                      {([...Array(p.cardCount)].map(function(e, i) {
                        return(
                          <img key={p.id + "_Card_" + i} src={CardBack} alt={"Card_Back_"+i}/>
                        );
                      }))
                      }
                    </div>
                  </div>
                </div>
              );
            } else {
              return (null);
            }
          })}
        </div>
        <div className="GameScreen-Middle">
          <div className="GameScreen-HoverCard">
            {(hoverCard !== null) &&
              <img src={require("../res/cards/" + hoverCard.dbId + "_" + hoverCard.expansionId + "_" + hoverCard.cardNumber + ".png").default} alt={"Big Card"}/>
            }
          </div>
          <div className="GameScreen-DraftingBoard">
            <img className="GameScreen-DrawPile" src={CardBack} onClick={() => handleDrawCard()} alt={"Draw"}/>
            {(gameInfo.lastPlayed.id !== "00000000-0000-0000-0000-000000000000") 
              ? (<img 
                  src={require("../res/cards/" + gameInfo.lastPlayed.dbId + "_" + gameInfo.lastPlayed.expansionId + "_" + gameInfo.lastPlayed.cardNumber + ".png").default}
                  onMouseEnter={() => setHoverCard(gameInfo.lastPlayed)}
                  onMouseLeave={() => setHoverCard(null)}
                  alt={"Last_Played"}
                />
              ) : <img src={CardBlank} alt={"Last_Played"}/>
            }
          </div>
          <div className="GameScreen-Archive">
            {(thisPlayer !== null && thisPlayer.archiveIncreases > 0) &&
              <div className="GameScreen-USBIndicator-CurrentPlayer">
                <img src={USB} alt={"Player_Archive_Indicator"}/>
              </div>
            }
            {(thisPlayer !== null && thisPlayer.handIncreases > 0) &&
              <div className="GameScreen-USBIndicator-CurrentPlayer">
                <img src={MemExpansion} alt={"Player_Hand_Indicator"}/>
              </div>
            }
            <div className="GameScreen-Archive-Bottom">
              {(gameInfo.currentPlayer === gameInfo.playerId) &&
                <div className="GameScreen-TurnIndicator-CurrentPlayer">
                  <img src={Arrow} alt={"Player_Turn_Indicator"}/>
                </div>
              }
              <div onClick={() => openArchive()} className="button">
                <h3 className="StartScreen-Submit-Text noselect">View Archive</h3>
              </div>
            </div>
            <div className="GameScreen-HandCount">
              <h3>{gameInfo.hand.length} / {gameInfo.handMax}</h3>
            </div>
          </div>
        </div>
        <div className="GameScreen-Bottom">
          <div className="GameScreen-Hand">
            {gameInfo.hand.map(function(c, i) {
              return (
                <div className="GameScreen-Hand-Card" key={"Player_Hand_"+i}>
                  <img className="GameScreen-Hand-Card-Image" 
                    src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} 
                    alt={"Hand_"+i}
                    onMouseEnter={() => setHoverCard(c)}
                    onMouseLeave={() => setHoverCard(null)}
                    onClick={() => handleCardClick(i)}/>
                  <div className="GameScreen-Hand-Card-Bottom">
                    <img src={Archive} alt={"Discard_"+i} onClick={() => handleArchive(c)}/>
                    <img src={Trash} alt={"Discard_"+i} onClick={() => handleDiscard(c)}/>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

export default GameScreen;