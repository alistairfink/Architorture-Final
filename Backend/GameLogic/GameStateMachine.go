package GameLogic

import (
	"Architorture-Backend/DataLayer"
	"Architorture-Backend/GameLogic/GameState"
	"Architorture-Backend/GameLogic/MessageType"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type gameStateMachine struct {
	expansions      []int
	hub             *hub
	gameState       GameState.GameState
	roomId          string
	players         []*Player
	DeckPile        []CardModel
	DiscardPile     []CardModel
	LastPlayed      CardModel
	db              *DataLayer.DatabaseConnection
	forward         bool
	currentTurnData CurrentTurnModel
	undoAction      UndoAction
}

type CardModel struct {
	Id              uuid.UUID `json:"id"`
	CardTypeId      int       `json:"cardTypeId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PlayImmediately bool      `json:"playImmediately"`
	Archivable      bool      `json:"archivable"`
	DbId            int       `json:"dbId"`
	ExpansionId     int       `json:"expansionId"`
	CardNumber      int       `json:"cardNumber"`
}

type CurrentTurnModel struct {
	currentPlayerIndex                   int
	archiveUsed                          bool
	actionCancelled                      bool
	cardResponseId                       int
	cardresponseUserId                   uuid.UUID
	shouldEndTurnAfterArchitortureAction bool
	actionCarriedOut                     bool
}

type UndoAction struct {
	actionUndone    bool
	cardId          int
	source          *Player
	target          *Player
	cardIds         []uuid.UUID
	numCards        int
	targetCardIndex int
	cardRequestName string
}

func InitGameStateMachine(roomId string, db *DataLayer.DatabaseConnection, hub *hub, expansions []int) *gameStateMachine {
	return &gameStateMachine{
		gameState:  GameState.Lobby,
		roomId:     roomId,
		players:    []*Player{},
		db:         db,
		forward:    true,
		hub:        hub,
		expansions: expansions,
	}
}

func (this *gameStateMachine) startGame() {
	log.Println("Making new game in room:", this.roomId)
	dbCards := this.db.GetCardsWithoutArchitortureOrSaveCards(this.expansions)
	cards := []CardModel{}
	for _, dbCard := range dbCards {
		for i := 0; i < dbCard.Quantity; i++ {
			cards = append(cards, CardModel{
				Id:              uuid.New(),
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

	this.DeckPile = cards
	this.shuffleDeckPile()

	for i := 0; i < 5; i++ {
		for j := 0; j < len(this.players); j++ {
			card := this.removeTopCard()
			for card.PlayImmediately {
				this.DeckPile = append(this.DeckPile, card)
				card = this.removeTopCard()
			}

			this.players[j].AddCardToHand(card)
		}
	}

	saveDBCards := this.db.GetSaveCards(this.expansions)
	saveCards := []CardModel{}
	for _, dbCard := range saveDBCards {
		for i := 0; i < dbCard.Quantity; i++ {
			saveCards = append(saveCards, CardModel{
				Id:              uuid.New(),
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
	this.shuffleCards(saveCards)

	for i := 0; i < len(this.players); i++ {
		this.players[i].AddCardToHand(saveCards[i])
	}

	this.DeckPile = append(this.DeckPile, saveCards[len(this.players)])

	architortureDBCards := this.db.GetArchitortureCards(this.expansions)
	architortureCards := []CardModel{}
	for _, dbCard := range architortureDBCards {
		for i := 0; i < dbCard.Quantity; i++ {
			architortureCards = append(architortureCards, CardModel{
				Id:              uuid.New(),
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

	this.shuffleCards(architortureCards)
	for i := 0; i < len(this.players)-1; i++ {
		this.DeckPile = append(this.DeckPile, architortureCards[i])
	}

	this.shuffleDeckPile()
	this.forward = true
	this.initializeNewTurnObject(0)
	this.gameState = GameState.InPlay
}

func (this *gameStateMachine) DrawCard(playerId uuid.UUID, data DrawCardRequest) {
	if this.currentTurnData.cardResponseId != -1 {
		return
	}

	if playerId != this.players[this.currentTurnData.currentPlayerIndex].Id {
		return
	}

	player := this.findPlayer(playerId)
	if player == nil || len(player.hand)+player.drawCount > player.handMax {
		return
	}

	actionRequired := false
	for i := 0; i < player.drawCount; i++ {
		topCard := this.removeTopCard()
		switch topCard.DbId {
		case 5, 6:
			this.handleArchitectureDumpCard(player, topCard)
			actionRequired = true
		case 7, 8:
			this.handleArchitectureMemoryLoss(player, topCard)
			actionRequired = true
		case 3, 4:
			this.currentTurnData.shouldEndTurnAfterArchitortureAction = true
			player.AddCardToHand(topCard)
			actionRequired = true
		default:
			player.AddCardToHand(topCard)
		}
	}

	player.drawCount = 1
	if !actionRequired {
		this.endTurn()
	}
}

func (this *gameStateMachine) ArchiveCard(playerId uuid.UUID, data ArchiveRequest) {
	if this.currentTurnData.cardResponseId != -1 ||
		playerId != this.players[this.currentTurnData.currentPlayerIndex].Id ||
		this.currentTurnData.archiveUsed {
		return
	}

	player := this.findPlayer(playerId)
	if player == nil {
		return
	}

	if !player.ArchiveCard(data.ArchiveCardId, data.UnarchiveCardId) {
		return
	}

	this.currentTurnData.archiveUsed = true
}

func (this *gameStateMachine) PlayCard(playerId uuid.UUID, data PlayCardRequest) {
	player := this.findPlayer(playerId)
	if player == nil {
		return
	}

	if player.Id != this.players[this.currentTurnData.currentPlayerIndex].Id || !player.HasCards(data.Cards) {
		return
	}

	if this.currentTurnData.cardResponseId != -1 {
		return
	}

	this.MapPlayCardLogic(player, data)
}

func (this *gameStateMachine) HandleCardResponse(playerId uuid.UUID, data CardResponse) {
	if this.currentTurnData.cardResponseId != data.ActionCardId {
		return
	}

	player := this.findPlayer(playerId)
	if player == nil {
		return
	}

	if player.Id != this.currentTurnData.cardresponseUserId {
		return
	}

	this.MapCardResponseLogic(player, data)
}

func (this *gameStateMachine) GetAvailableCardNames() []string {
	return this.db.GetAvailableCardNames(this.expansions)
}

func (this *gameStateMachine) ValidateNewPlayer(newPlayer Player) bool {
	if this.gameState != GameState.Lobby {
		return false
	}

	if len(this.players) == 5 {
		return false
	}

	return true
}

func (this *gameStateMachine) DiscardCard(playerId uuid.UUID, data PlayCardRequest) {
	if playerId != this.players[this.currentTurnData.currentPlayerIndex].Id {
		return
	}

	player := this.findPlayer(playerId)
	if player == nil || len(data.Cards) != 1 {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.DeckPile = append(this.DeckPile, card)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) UndoCard(playerId uuid.UUID, data UndoCardResponse) {
	player := this.findPlayer(playerId)
	if player == nil {
		return
	}

	this.handleUndo(player, data)
}

func (this *gameStateMachine) shuffleDeckPile() {
	this.shuffleCards(this.DeckPile)
}

func (this *gameStateMachine) shuffleCards(cards []CardModel) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(cards) > 0 {
		n := len(cards)
		randIndex := r.Intn(n)
		cards[n-1], cards[randIndex] = cards[randIndex], cards[n-1]
		cards = cards[:n-1]
	}
}

func (this *gameStateMachine) removeTopCard() CardModel {
	if len(this.DeckPile) < 4 {
		cards := this.DiscardPile
		this.DeckPile = append(this.DeckPile, cards...)
		this.DiscardPile = []CardModel{}
	}

	topCard := this.DeckPile[0]
	this.DeckPile = this.DeckPile[1:]

	return topCard
}

func (this *gameStateMachine) initializeNewTurnObject(playerTurnIndex int) {
	this.currentTurnData = CurrentTurnModel{
		currentPlayerIndex:                   playerTurnIndex,
		archiveUsed:                          false,
		actionCancelled:                      false,
		cardResponseId:                       -1,
		cardresponseUserId:                   uuid.Nil,
		shouldEndTurnAfterArchitortureAction: false,
		actionCarriedOut:                     false,
	}
}

func (this *gameStateMachine) endTurn() {
	allEliminated := true
	for _, p := range this.players {
		allEliminated = allEliminated && p.Eliminated
	}

	if allEliminated {
		return
	}

	nextTurnIndex := this.getNextPlayerIndex(this.currentTurnData.currentPlayerIndex)
	for this.players[nextTurnIndex].Eliminated {
		nextTurnIndex = this.getNextPlayerIndex(nextTurnIndex)
	}

	for this.players[nextTurnIndex].skipCount > 0 {
		this.players[nextTurnIndex].DecrementSkipCounter()
		nextTurnIndex = this.getNextPlayerIndex(nextTurnIndex)
	}

	this.initializeNewTurnObject(nextTurnIndex)
}

func (this *gameStateMachine) getNextPlayerIndex(currentIndex int) int {
	if this.forward {
		return (currentIndex + 1) % len(this.players)
	}

	nextTurn := currentIndex - 1
	if nextTurn < 0 {
		nextTurn = len(this.players) - 1
	}

	return nextTurn
}

func (this *gameStateMachine) findPlayer(playerId uuid.UUID) *Player {
	var player *Player = nil
	for _, p := range this.players {
		if p.Id == playerId {
			player = p
			break
		}
	}

	return player
}

func (this *gameStateMachine) discardCard(card CardModel, insertIntoLastPlayed bool) {
	this.DiscardPile = append(this.DiscardPile, card)
	if insertIntoLastPlayed {
		this.LastPlayed = card
	}
}

func (this *gameStateMachine) createNewUndoAction(cardId int) {
	this.undoAction = UndoAction{
		actionUndone:    false,
		cardId:          cardId,
		source:          nil,
		target:          nil,
		cardIds:         []uuid.UUID{},
		numCards:        -1,
		targetCardIndex: -1,
		cardRequestName: "",
	}
}

func (this *gameStateMachine) insertRandomly(card CardModel) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	index := r.Intn(len(this.DeckPile))
	this.DeckPile = append(this.DeckPile[:index], append([]CardModel{card}, this.DeckPile[index:]...)...)
}
