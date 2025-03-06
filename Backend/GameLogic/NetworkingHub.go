package GameLogic

import (
	"Architorture-Backend/DataLayer"
	"Architorture-Backend/GameLogic/GameState"
	"Architorture-Backend/GameLogic/MessageType"
	"Architorture-Backend/GameLogic/RequestType"
	"Architorture-Backend/Models"
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

type hub struct {
	rooms      map[string]*gameStateMachine
	broadcast  chan GameRequest
	register   chan registerPlayer
	unregister chan unregisterPlayer
	db         *DataLayer.DatabaseConnection
}

type unregisterPlayer struct {
	room     string
	playerId uuid.UUID
}

type registerPlayer struct {
	player    *Player
	expansion int
}

func InitHub(db *DataLayer.DatabaseConnection) *hub {
	return &hub{
		broadcast:  make(chan GameRequest),
		register:   make(chan registerPlayer),
		unregister: make(chan unregisterPlayer),
		rooms:      make(map[string]*gameStateMachine),
		db:         db,
	}
}

func (this *hub) Run() {
	for {
		select {
		case regInfo := <-this.register:
			this.handleRegister(regInfo)
		case sub := <-this.unregister:
			this.handleUnregister(sub)
		case message := <-this.broadcast:
			if message.RequestType == RequestType.Ready {
				var data ReadyRequest
				json.Unmarshal([]byte(message.Body), &data)
				this.handlePlayerReady(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.Archive {
				var data ArchiveRequest
				json.Unmarshal([]byte(message.Body), &data)
				this.handleArchive(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.PlayCard {
				var data PlayCardRequest
				json.Unmarshal([]byte(message.Body), &data)
				this.handlePlayCard(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.DrawCard {
				var data DrawCardRequest
				json.Unmarshal([]byte(message.Body), &data)
				this.handleDrawCard(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.CardActionResponse {
				var data CardResponse
				json.Unmarshal([]byte(message.Body), &data)
				this.handleCardResponse(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.Discard {
				var data PlayCardRequest
				json.Unmarshal([]byte(message.Body), &data)
				this.handleDiscardCard(message.playerId, message.roomId, data)
			} else if message.RequestType == RequestType.Undo {
				var data UndoCardResponse
				json.Unmarshal([]byte(message.Body), &data)
				this.handleUndoCard(message.playerId, message.roomId, data)
			}
		}
	}
}

func (this *hub) CheckRoomId(roomId string) Models.CheckRoomIdResultModel {
	room, exists := this.rooms[roomId]
	return Models.CheckRoomIdResultModel{
		Response: exists && room.gameState == GameState.Lobby,
	}
}

func (this *hub) GetAllRooms() []string {
	roomIds := make([]string, len(this.rooms))
	itterator := 0
	for key, _ := range this.rooms {
		roomIds[itterator] = key
		itterator++
	}

	return roomIds
}

func (this *hub) GetCards(roomId string) []string {
	game, exists := this.rooms[roomId]
	if !exists {
		return []string{}
	}

	return game.GetAvailableCardNames()
}

func (this *hub) GetCardsByExpansion(expansionInt int) []CardModel {
	expansions := []int{}
	for i := 1; i <= expansionInt; i++ {
		expansions = append(expansions, i)
	}

	dbCards := this.db.GetAvailableCards(expansions)
	cards := []CardModel{}
	for _, dbCard := range dbCards {
		for i := 0; i < dbCard.Quantity; i++ {
			cards = append(cards, CardModel{
				Id:              uuid.Nil,
				CardTypeId:      dbCard.CardTypeId,
				Name:            dbCard.Name,
				Description:     dbCard.Description,
				PlayImmediately: dbCard.PlayImmediately,
				Archivable:      dbCard.Archivable,
				DbId:            dbCard.Id,
				CardNumber:      i,
				ExpansionId:     dbCard.ExpansionId,
			})
		}
	}

	return cards
}

func (this *hub) GetDrawPile(roomId string) []CardModel {
	game, exists := this.rooms[roomId]
	println(exists)
	if !exists {
		return []CardModel{}
	}

	println(len(game.DeckPile))
	return game.DeckPile
}

func (this *hub) handleRegister(regInfo registerPlayer) {
	roomId := regInfo.player.room
	if roomId == "" {
		exists := true
		for exists {
			guid := uuid.New().String()
			roomId = string(guid[:5])
			_, exists = this.rooms[roomId]
		}

		expansions := []int{}
		for i := 1; i <= regInfo.expansion; i++ {
			expansions = append(expansions, i)
		}

		this.rooms[roomId] = InitGameStateMachine(roomId, this.db, this, expansions)
	} else if _, exists := this.rooms[roomId]; !exists {
		return
	}

	game := this.rooms[roomId]
	if !game.ValidateNewPlayer(*regInfo.player) {
		return
	}

	log.Println("User connecting\n  Username:", regInfo.player.Username, "\n  RoomId:", roomId)
	regInfo.player.SetRoomId(roomId)
	regInfo.player.UpdatePlayerNumber(len(game.players))
	game.players = append(game.players, regInfo.player)
	this.HandleGameUpdate(roomId, MessageType.LobbyStart)
}

func (this *hub) handleUnregister(unregInfo unregisterPlayer) {
	roomId := unregInfo.room
	game := this.rooms[roomId]
	log.Println("Disconnecting Player\n  Player Id:", unregInfo.playerId, "\n  RoomId:", roomId)

	elimCount := 0
	playerIndex := -1
	for i, player := range game.players {
		if player == nil {
			elimCount++
			continue
		}

		if player.Id == unregInfo.playerId {
			player.EliminatePlayer()
			player.disconnected = true
			playerIndex = i
		}

		if player.Eliminated {
			elimCount++
		}

	}

	if len(game.players) == elimCount {
		log.Println("Deleting Game:", roomId)
		delete(this.rooms, game.roomId)
		return
	}

	if game.gameState == GameState.Lobby {
		this.rooms[roomId].players = append(this.rooms[roomId].players[:playerIndex], this.rooms[roomId].players[playerIndex+1:]...)
	}

	this.HandleGameUpdate(roomId, MessageType.GameInfoUpdate)
}

func (this *hub) handlePlayerReady(playerId uuid.UUID, roomId string, request ReadyRequest) {
	log.Println("Player Toggling Ready Status\n  PlayerId:", playerId, "\n  RoomId:", roomId, "\n  NewReadyStatus:", request.IsReady)
	allReady := true
	for _, player := range this.rooms[roomId].players {
		if player.Id == playerId {
			player.SetReadyStatus(request.IsReady)
		}

		allReady = allReady && player.IsReady
	}

	if allReady && len(this.rooms[roomId].players) > 1 {
		this.rooms[roomId].startGame()
		this.HandleGameUpdate(roomId, MessageType.GameStart)
		return
	}

	this.HandleGameUpdate(roomId, MessageType.GameInfoUpdate)
}

func (this *hub) handleArchive(playerId uuid.UUID, roomId string, data ArchiveRequest) {
	log.Println("Archiving Card\n  Player:", playerId, "\n  Archive Card:", data.ArchiveCardId, "\n  Unarchive Card:", data.UnarchiveCardId)
	game := this.rooms[roomId]
	game.ArchiveCard(playerId, data)
	this.HandleGameUpdate(roomId, MessageType.GameInfoUpdate)
}

func (this *hub) handlePlayCard(playerId uuid.UUID, roomId string, data PlayCardRequest) {
	logText := "Playing Card:\n  Player: " + playerId.String() + "\n  Target Player: " + data.TargetPlayerId.String() + "\n  Cards:"
	for _, card := range data.Cards {
		logText += "\n    " + card.String()
	}
	log.Println(logText)

	game := this.rooms[roomId]
	game.PlayCard(playerId, data)
}

func (this *hub) handleDrawCard(playerId uuid.UUID, roomId string, data DrawCardRequest) {
	log.Println("Drawing Card\n  Player:", playerId)
	game := this.rooms[roomId]
	game.DrawCard(playerId, data)
	this.HandleGameUpdate(roomId, MessageType.TurnStart)
}

func (this *hub) handleCardResponse(playerId uuid.UUID, roomId string, data CardResponse) {
	log.Println("Response for Card:\n  Player:", playerId)
	game := this.rooms[roomId]
	game.HandleCardResponse(playerId, data)
}

func (this *hub) handleDiscardCard(playerId uuid.UUID, roomId string, data PlayCardRequest) {
	log.Println("Discarding Card:\n  Player:", playerId.String(), "\n  Card:", data.Cards[0])
	game := this.rooms[roomId]
	game.DiscardCard(playerId, data)
}

func (this *hub) handleUndoCard(playerId uuid.UUID, roomId string, data UndoCardResponse) {
	log.Println("Undoing Card:\n  Player:", playerId.String(), "\n  Card:", data.CardId)
	game := this.rooms[roomId]
	game.UndoCard(playerId, data)
}

func (this *hub) HandleGameUpdate(roomId string, msgType MessageType.MessageType) {
	game := this.rooms[roomId]
	data := GameInfoUpdate{
		MessageType:   msgType,
		Players:       game.players,
		RoomId:        roomId,
		CurrentPlayer: game.players[game.currentTurnData.currentPlayerIndex].Id,
		LastPlayed:    game.LastPlayed,
		Expansions:    game.expansions,
	}

	for _, player := range this.rooms[roomId].players {
		player.UpdateCardCount()
	}

	for _, player := range this.rooms[roomId].players {
		if player.disconnected {
			continue
		}

		data.Hand = player.hand
		data.PlayerId = player.Id
		data.Archive = player.archive
		data.HandMax = player.handMax
		data.ArchiveMax = player.archiveMax
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		player.send <- dataBytes
	}
}

func (this *hub) PrintGame(roomId string) {
	log.Println(this.rooms[roomId])
}

func (this *hub) HandleSendArchitortureUpdate(roomId string, player *Player) {
	data := ArchitortureCardDrawn{
		MessageType: MessageType.Architorture,
		Player:      player,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	for _, p := range this.rooms[roomId].players {
		p.send <- dataBytes
	}
}
