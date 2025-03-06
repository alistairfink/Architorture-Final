package GameLogic

import (
	"Architorture-Backend/GameLogic/MessageType"
	"Architorture-Backend/GameLogic/RequestType"
	"Architorture-Backend/GameLogic/UndoRequestState"
	"github.com/google/uuid"
)

type GameRequest struct {
	RequestType RequestType.RequestType `json:"requestType"`
	Body        string                  `json:"body"`
	playerId    uuid.UUID
	roomId      string
}

type ReadyRequest struct {
	IsReady bool `json:"isReady"`
}

type ArchiveRequest struct {
	ArchiveCardId   uuid.UUID `json:"archiveCardId"`
	UnarchiveCardId uuid.UUID `json:"unarchiveCardId"`
}

type PlayCardRequest struct {
	Cards           []uuid.UUID `json:"cards"`
	TargetPlayerId  uuid.UUID   `json:"targetPlayerId"`
	CardsToGive     []uuid.UUID `json:"cardsToGive"`
	CardrequestName string      `json:"cardrequestName"`
	TargetCardIndex int         `json:"targetCardIndex"`
}

type DrawCardRequest struct {
}

type GameInfoUpdate struct {
	MessageType   MessageType.MessageType `json:"messageType"`
	RoomId        string                  `json:"roomId"`
	Players       []*Player               `json:"players"`
	Hand          []CardModel             `json:"hand"`
	PlayerId      uuid.UUID               `json:"playerId"`
	CurrentPlayer uuid.UUID               `json:"currentPlayer"`
	Archive       []CardModel             `json:"archive"`
	HandMax       int                     `json:"handMax"`
	ArchiveMax    int                     `json:"archiveMax"`
	LastPlayed    CardModel               `json:"lastPlayed"`
	Expansions    []int                   `json:"expansions"`
}

type CardActionData struct {
	MessageType      MessageType.MessageType `json:"messageType"`
	CardDbId         int                     `json:"cardDbId"`
	Cards            []CardModel             `json:"cards"`
	RequiresResponse bool                    `json:"requiresResponse"`
}

type CardResponse struct {
	ActionCardId int         `json:"actionCardId"`
	Cards        []uuid.UUID `json:"cards"`
	CardIndex    int         `json"cardIndex"`
	Target       uuid.UUID   `json:"target"`
}

type ArchitortureCardDrawn struct {
	MessageType MessageType.MessageType `json:"messageType"`
	Player      *Player                 `json:"player"`
}

type UndoCardRequest struct {
	MessageType MessageType.MessageType           `json:"messageType"`
	Stage       UndoRequestState.UndoRequestState `json:"stage"`
	Cards       []CardModel                       `json:"card"`
	Undone      bool                              `json:"undone"`
}

type UndoCardResponse struct {
	CardId uuid.UUID `json:"cardId"`
	Use    bool      `json:"use"`
}
