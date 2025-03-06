package GameLogic

import (
	"Architorture-Backend/Constants"
	"Architorture-Backend/DataLayer/DataModels/CardTypeEnum"
	"Architorture-Backend/GameLogic/MessageType"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Player struct {
	Username           string    `json:"username"`
	IsReady            bool      `json:"isReady"`
	Id                 uuid.UUID `json:"id"`
	PlayerNumber       int       `json:"playerNumber"`
	ArchiveIncreases   int       `json:"archiveIncreases"`
	HandIncreases      int       `json:"handIncreases"`
	CardCount          int       `json:"cardCount"`
	Eliminated         bool      `json:"eliminated"`
	disconnected       bool
	hand               []CardModel
	archive            []CardModel
	handMax            int
	archiveMax         int
	webSocket          *websocket.Conn
	room               string
	hub                *hub
	send               chan []byte
	skipCount          int
	timeExtensionUsage map[uuid.UUID]uuid.UUID
	drawCount          int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request, roomId string, username string, expansion int, hub *hub) {
	webSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	user := &Player{
		webSocket:          webSocket,
		Id:                 uuid.New(),
		Username:           username,
		hub:                hub,
		room:               roomId,
		send:               make(chan []byte),
		IsReady:            false,
		hand:               []CardModel{},
		archive:            []CardModel{},
		handMax:            7,
		archiveMax:         3,
		skipCount:          0,
		ArchiveIncreases:   0,
		HandIncreases:      0,
		Eliminated:         false,
		disconnected:       false,
		PlayerNumber:       -1,
		timeExtensionUsage: make(map[uuid.UUID]uuid.UUID),
		drawCount:          1,
	}

	regInfo := registerPlayer{
		player:    user,
		expansion: expansion,
	}

	hub.register <- regInfo
	go user.writePump()
	go user.readPump()
}

func (this *Player) AddCardToHand(card CardModel) {
	this.hand = append(this.hand, card)
	if card.CardTypeId == CardTypeEnum.Architorture {
		game := this.hub.rooms[this.room]
		if !game.HandleArchitortureCard(this, card) {
			this.EliminatePlayer()
			this.hub.HandleGameUpdate(this.room, MessageType.GameInfoUpdate)
		}
	}
}

func (this *Player) AddCardToArchive(card CardModel) {
	this.archive = append(this.archive, card)
}

func (this *Player) EliminatePlayer() {
	game := this.hub.rooms[this.room]
	this.Eliminated = true
	if game.players[game.currentTurnData.currentPlayerIndex].Id == this.Id {
		game.endTurn()
	}
}

func (this *Player) SetRoomId(roomId string) {
	this.room = roomId
}

func (this *Player) SetReadyStatus(status bool) {
	this.IsReady = status
}

func (this *Player) ArchiveCard(archiveCardId uuid.UUID, unarchiveCardId uuid.UUID) bool {
	if archiveCardId == uuid.Nil && unarchiveCardId == uuid.Nil {
		println("-1")
		return false
	} else if archiveCardId == uuid.Nil && len(this.hand) >= this.handMax {
		println("-2")
		return false
	} else if unarchiveCardId == uuid.Nil && len(this.archive) >= this.archiveMax {
		println("-3")
		return false
	} else if archiveCardId != uuid.Nil && !this.HasCards([]uuid.UUID{archiveCardId}) {
		println("-4")
		return false
	} else if unarchiveCardId != uuid.Nil && !this.HasCardInArchive(unarchiveCardId) {
		println("-5")
		return false
	}

	println("--1")

	if archiveCardId != uuid.Nil {
		archiveCard := this.GetCard(archiveCardId)
		if !archiveCard.Archivable {
			return false
		}
	}

	println("--2")
	if unarchiveCardId != uuid.Nil {
		unarchiveCard := this.GetCardFromArchive(unarchiveCardId)
		if !unarchiveCard.Archivable {
			return false
		}
	}

	println("--3")
	if unarchiveCardId == uuid.Nil {

		println("---1")
		card := this.RemoveCard(archiveCardId)
		this.AddCardToArchive(card)
	} else if archiveCardId == uuid.Nil {

		println("---2")
		card := this.RemoveArchivedCard(unarchiveCardId)
		this.AddCardToHand(card)
	} else {
		println("---3")
		archiveCard := this.RemoveCard(archiveCardId)
		unarchiveCard := this.RemoveArchivedCard(unarchiveCardId)
		this.AddCardToHand(unarchiveCard)
		this.AddCardToArchive(archiveCard)
	}

	return true
}

func (this *Player) HasCards(cards []uuid.UUID) bool {
	cardMap := make(map[uuid.UUID]bool)
	for _, card := range this.hand {
		cardMap[card.Id] = true
	}

	valid := true
	for _, card := range cards {
		valid = valid && cardMap[card]
	}

	return valid
}

func (this *Player) HasCardInArchive(unarchiveCardId uuid.UUID) bool {
	for _, card := range this.archive {
		if card.Id == unarchiveCardId {
			return true
		}
	}

	return false
}

func (this *Player) GetCardDBType(cardId uuid.UUID) int {
	for _, card := range this.hand {
		if card.Id == cardId {
			return card.DbId
		}
	}

	return -1
}

func (this *Player) RemoveCard(cardId uuid.UUID) CardModel {
	var card CardModel
	cardIndex := -1
	for i, c := range this.hand {
		if c.Id == cardId {
			cardIndex = i
			card = c
			break
		}
	}

	this.hand = append(this.hand[:cardIndex], this.hand[cardIndex+1:]...)

	if architortureCardId, exists := this.timeExtensionUsage[cardId]; exists {
		architortureCard := this.RemoveCard(architortureCardId)
		delete(this.timeExtensionUsage, cardId)
		this.AddCardToHand(architortureCard)
	}

	return card
}

func (this *Player) RemoveArchivedCard(cardId uuid.UUID) CardModel {
	var card CardModel
	cardIndex := -1
	for i, c := range this.archive {
		if c.Id == cardId {
			cardIndex = i
			card = c
			break
		}
	}

	this.archive = append(this.archive[:cardIndex], this.archive[cardIndex+1:]...)
	return card
}

func (this *Player) GetCard(cardId uuid.UUID) CardModel {
	var card CardModel
	for _, c := range this.hand {
		if c.Id == cardId {
			card = c
			break
		}
	}

	return card
}

func (this *Player) GetCardFromArchive(cardId uuid.UUID) CardModel {
	var card CardModel
	for _, c := range this.archive {
		if c.Id == cardId {
			card = c
			break
		}
	}

	return card
}

func (this *Player) DecrementSkipCounter() {
	this.skipCount--
	if this.skipCount < 0 {
		this.skipCount = 0
	}
}

func (this *Player) SetSkipCount(skipCount int) {
	this.skipCount = skipCount
}

func (this *Player) IncrementUSBMemory(incrementAmount int) {
	this.archiveMax += incrementAmount
	this.ArchiveIncreases++
}

func (this *Player) IncrementHandSize(incrementAmount int) {
	this.handMax += incrementAmount
	this.HandIncreases++
}

func (this *Player) DecrementHandSize(incrementAmount int) {
	this.handMax -= incrementAmount
	this.HandIncreases--
}

func (this *Player) ClearHand() []CardModel {
	cards := this.hand
	this.hand = []CardModel{}
	return cards
}

func (this *Player) CanDecrementUSBMemory(decrementAmount int) bool {
	if this.ArchiveIncreases < decrementAmount {
		return false
	}

	return true
}

func (this *Player) DecrementUSBMemory(decrementAmount int) {
	for i := 0; i < decrementAmount; i++ {
		this.ArchiveIncreases--
		this.archiveMax--
	}
}

func (this *Player) RemoveCardByIndex(index int) CardModel {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	cardIndex := r.Intn(len(this.hand))
	card := this.hand[cardIndex]
	return this.RemoveCard(card.Id)
}

func (this *Player) ValidateIndex(index int) bool {
	return index >= 0 && index < len(this.hand)
}

func (this *Player) GetCardByName(cardRequestName string) (CardModel, bool) {
	for _, card := range this.hand {
		if card.Name == cardRequestName {
			return this.RemoveCard(card.Id), true
		}
	}

	return CardModel{}, false
}

func (this *Player) UpdateCardCount() {
	this.CardCount = len(this.hand)
}

func (this *Player) UpdatePlayerNumber(num int) {
	this.PlayerNumber = num
}

func (this *Player) AddTimeExtensionUsage(timeExtensionId uuid.UUID, architortureId uuid.UUID) {
	this.timeExtensionUsage[timeExtensionId] = architortureId
}

func (this *Player) HasUndoCard() bool {
	for _, card := range this.hand {
		if card.DbId == 27 ||
			card.DbId == 28 ||
			card.DbId == 34 {
			return true
		}
	}

	return false
}

func (this *Player) HasNotInThisLifetimeCard() bool {
	for _, card := range this.hand {
		if card.DbId == 34 {
			return true
		}
	}

	return false
}

func (this *Player) writePump() {
	ticker := time.NewTicker(Constants.PingPeriod)
	defer func() {
		ticker.Stop()
		this.webSocket.Close()
	}()

	for {
		select {
		case message, ok := <-this.send:
			if !ok {
				this.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := this.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := this.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (this *Player) readPump() {
	defer func() {
		this.hub.unregister <- unregisterPlayer{
			room:     this.room,
			playerId: this.Id,
		}
		this.webSocket.Close()
	}()

	this.webSocket.SetReadLimit(Constants.MaxMessageSize)
	this.webSocket.SetReadDeadline(time.Now().Add(Constants.PongWait))
	this.webSocket.SetPongHandler(func(string) error {
		this.webSocket.SetReadDeadline(time.Now().Add(Constants.PongWait))
		return nil
	})

	for {
		_, msg, err := this.webSocket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var data GameRequest
		json.Unmarshal(msg, &data)
		data.playerId = this.Id
		data.roomId = this.room
		this.hub.broadcast <- data
	}
}

func (this *Player) write(mt int, payload []byte) error {
	this.webSocket.SetWriteDeadline(time.Now().Add(Constants.WriteWait))
	return this.webSocket.WriteMessage(mt, payload)
}
