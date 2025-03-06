package GameLogic

import (
	"Architorture-Backend/DataLayer/DataModels/CardTypeEnum"
	"Architorture-Backend/GameLogic/MessageType"
	"Architorture-Backend/GameLogic/UndoRequestState"
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

func (this *gameStateMachine) MapPlayCardLogic(player *Player, data PlayCardRequest) {
	this.currentTurnData.actionCarriedOut = false
	switch player.GetCardDBType(data.Cards[0]) {
	case 9, 10:
		this.handleResourceShuffle(player, data)
	case 11, 12, 13:
		this.handleMemoryExpansionCard(player, data)
	case 14, 15, 16:
		this.handleUSBExpansionCard(player, data)
	case 17, 18:
		this.handleFoundAUsbCard(player, data)
	case 19:
		this.handleVendettaCard(player, data, false)
	case 20:
		this.handlePreviewCard(player, data, 2, false)
	case 21:
		this.handlePreviewCard(player, data, 4, false)
	case 22:
		this.handleThankUNextCard(player, data)
	case 23:
		this.handleShuffleCard(player, data)
	case 24:
		this.handleReverseCard(player, data)
	case 25, 26:
		this.handlePreviewCard(player, data, 2, true)
	case 31, 32:
		this.handleVendettaCard(player, data, true)
	case 33:
		this.handleCoopCard(player, data)
	case 36:
		this.handleSharePreviewCard(player, data)
	case 38:
		this.handleAssistanceCard(player, data)
	case 29, 30, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		this.handleKnowledgeCards(player, data)
	case -1:
		log.Println("Error playing card.")
		return
	}
}

func (this *gameStateMachine) handleReverseCard(player *Player, data PlayCardRequest) {
	card := player.RemoveCard(data.Cards[0])
	this.forward = !this.forward
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.endTurn()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleCoopCard(player *Player, data PlayCardRequest) {
	card := player.RemoveCard(data.Cards[0])
	player.SetSkipCount(2)
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.endTurn()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleUSBExpansionCard(player *Player, data PlayCardRequest) {
	player.IncrementUSBMemory(1)
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleMemoryExpansionCard(player *Player, data PlayCardRequest) {
	player.IncrementHandSize(1)
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleThankUNextCard(player *Player, data PlayCardRequest) {
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.endTurn()
	this.createNewUndoAction(-1)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handlePreviewCard(player *Player, data PlayCardRequest, cardCount int, requiresResponse bool) {
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.previewLogic(player, cardCount, card.DbId, requiresResponse)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) previewLogic(player *Player, cardCount int, cardDbId int, requiresResponse bool) {
	top2Cards := this.DeckPile[:cardCount]
	this.currentTurnData.actionCarriedOut = true
	if requiresResponse {
		this.currentTurnData.cardResponseId = cardDbId
		this.currentTurnData.cardresponseUserId = player.Id
	}

	previewCardObject := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         cardDbId,
		Cards:            top2Cards,
		RequiresResponse: requiresResponse,
	}

	dataBytes, err := json.Marshal(previewCardObject)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	player.send <- dataBytes
}

func (this *gameStateMachine) handleShuffleCard(player *Player, data PlayCardRequest) {
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.shuffleDeckPile()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleResourceShuffle(player *Player, data PlayCardRequest) {
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.resourceShuffleLogic()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) resourceShuffleLogic() {
	cards := []CardModel{}
	for _, p := range this.players {
		cards = append(cards, p.ClearHand()...)
	}

	playerIndex := 0
	this.shuffleCards(cards)
	for _, card := range cards {
		this.players[playerIndex].AddCardToHand(card)
		playerIndex = (playerIndex + 1) % len(this.players)
	}
}

func (this *gameStateMachine) handleVendettaCard(player *Player, data PlayCardRequest, useTarget bool) {
	if !player.HasCards(data.CardsToGive) {
		return
	}

	var target *Player = nil
	if !useTarget {
		target = this.players[this.getNextPlayerIndex(this.currentTurnData.currentPlayerIndex)]
	} else {
		for _, p := range this.players {
			if p.Id == data.TargetPlayerId {
				target = p
				break
			}
		}
	}

	if target == nil {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	if target.HasUndoCard() {
		this.createNewUndoAction(card.DbId)
		this.undoAction.source = player
		this.undoAction.target = target
		this.undoAction.cardIds = data.CardsToGive
		data := UndoCardRequest{
			MessageType: MessageType.Undo,
			Cards:       []CardModel{card},
			Undone:      false,
			Stage:       UndoRequestState.Initial,
		}

		this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		target.send <- dataBytes
		return
	}

	this.vendettaLogic(player, target, data.CardsToGive[0])
	this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
}

func (this *gameStateMachine) vendettaLogic(player *Player, target *Player, cardToGive uuid.UUID) {
	card := player.RemoveCard(cardToGive)
	target.AddCardToHand(card)
	target.drawCount = 2
	this.endTurn()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleFoundAUsbCard(player *Player, data PlayCardRequest) {
	var target *Player = nil
	for _, p := range this.players {
		if p.Id == data.TargetPlayerId {
			target = p
			break
		}
	}

	if target == nil {
		return
	}

	if !target.CanDecrementUSBMemory(1) {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	if target.HasUndoCard() {
		this.createNewUndoAction(card.DbId)
		this.undoAction.source = player
		this.undoAction.target = target
		data := UndoCardRequest{
			MessageType: MessageType.Undo,
			Cards:       []CardModel{card},
			Undone:      false,
			Stage:       UndoRequestState.Initial,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		target.send <- dataBytes
		return
	}

	this.foundAUsbLogic(player, target)
	this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
}

func (this *gameStateMachine) foundAUsbLogic(player *Player, target *Player) {
	target.DecrementUSBMemory(1)
	player.IncrementUSBMemory(1)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleSharePreviewCard(player *Player, data PlayCardRequest) {
	var target *Player = nil
	for _, p := range this.players {
		if p.Id == data.TargetPlayerId {
			target = p
			break
		}
	}

	if target == nil {
		return
	}

	this.currentTurnData.actionCancelled = false
	this.currentTurnData.actionCarriedOut = false
	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	this.createNewUndoAction(-1)
	this.sharePreviewLogic(player, target, card.DbId)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) sharePreviewLogic(player *Player, target *Player, cardDbId int) {
	top3Cards := this.DeckPile[:3]
	this.currentTurnData.actionCarriedOut = true
	previewCardObject := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         cardDbId,
		Cards:            top3Cards,
		RequiresResponse: true,
	}

	dataBytes, err := json.Marshal(previewCardObject)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	player.send <- dataBytes

	previewCardObject.RequiresResponse = false
	dataBytes, err = json.Marshal(previewCardObject)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	target.send <- dataBytes
}

func (this *gameStateMachine) handleAssistanceCard(player *Player, data PlayCardRequest) {
	var target *Player = nil
	for _, p := range this.players {
		if p.Id == data.TargetPlayerId {
			target = p
			break
		}
	}

	if target == nil {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, true)
	if target.HasUndoCard() {
		this.createNewUndoAction(card.DbId)
		this.undoAction.source = player
		this.undoAction.target = target
		data := UndoCardRequest{
			MessageType: MessageType.Undo,
			Cards:       []CardModel{card},
			Undone:      false,
			Stage:       UndoRequestState.Initial,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		target.send <- dataBytes
		return
	}

	this.assistanceLogic(target, card.DbId)
}

func (this *gameStateMachine) assistanceLogic(target *Player, cardDbId int) {
	requestCardSelectionObject := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         cardDbId,
		RequiresResponse: true,
		Cards:            []CardModel{},
	}

	this.currentTurnData.cardResponseId = cardDbId
	this.currentTurnData.cardresponseUserId = target.Id

	dataBytes, err := json.Marshal(requestCardSelectionObject)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
	target.send <- dataBytes
}

func (this *gameStateMachine) handleKnowledgeCards(player *Player, data PlayCardRequest) {
	if len(data.Cards) != 2 && len(data.Cards) != 3 {
		return
	}

	var target *Player = nil
	for _, p := range this.players {
		if p.Id == data.TargetPlayerId {
			target = p
			break
		}
	}

	if target == nil {
		return
	}

	if len(data.Cards) == 2 {
		if !target.ValidateIndex(data.TargetCardIndex) {
			return
		}
	} else {

	}

	previousDbId := -1
	cards := []CardModel{}
	for _, cardId := range data.Cards {
		if player.GetCardDBType(cardId) == -1 {
			return
		}

		card := player.GetCard(cardId)
		if previousDbId == 29 || previousDbId == 30 || previousDbId == -1 {
			previousDbId = card.DbId
		} else {
			switch card.DbId {
			case 39, 40, 41, 42:
				if previousDbId != 39 && previousDbId != 40 && previousDbId != 41 && previousDbId != 42 {
					return
				}
			case 43, 44, 45:
				if previousDbId != 43 && previousDbId != 44 && previousDbId != 45 {
					return
				}
			case 46, 47, 48:
				if previousDbId != 46 && previousDbId != 47 && previousDbId != 48 {
					return
				}
			case 49, 50, 51:
				if previousDbId != 49 && previousDbId != 50 && previousDbId != 51 {
					return
				}
			case 52, 53, 54:
				if previousDbId != 52 && previousDbId != 53 && previousDbId != 54 {
					return
				}
			case 55, 56, 57:
				if previousDbId != 55 && previousDbId != 56 && previousDbId != 57 {
					return
				}
			}

			previousDbId = card.DbId
		}

		cards = append(cards, card)
	}

	for _, cardId := range data.Cards {
		player.RemoveCard(cardId)
	}

	for _, card := range cards {
		this.discardCard(card, true)
	}

	if target.HasUndoCard() {
		this.createNewUndoAction(cards[0].DbId)
		this.undoAction.source = player
		this.undoAction.target = target
		this.undoAction.numCards = len(cards)
		this.undoAction.targetCardIndex = data.TargetCardIndex
		this.undoAction.cardRequestName = data.CardrequestName
		data := UndoCardRequest{
			MessageType: MessageType.Undo,
			Cards:       cards,
			Undone:      false,
			Stage:       UndoRequestState.Initial,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		target.send <- dataBytes
		return
	}

	this.knowledgeLogic(player, target, len(cards), data.TargetCardIndex, data.CardrequestName)
}

func (this *gameStateMachine) knowledgeLogic(player *Player, target *Player, numberOfCards int, targetCardIndex int, cardRequestName string) {
	if numberOfCards == 2 {
		cardToTransfer := target.RemoveCardByIndex(targetCardIndex)
		player.AddCardToHand(cardToTransfer)
	} else {
		cardToTranfer, hasCard := target.GetCardByName(cardRequestName)
		if hasCard {
			player.AddCardToHand(cardToTranfer)
		}
	}

	this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
}

func (this *gameStateMachine) handleArchitectureDumpCard(player *Player, topCard CardModel) {
	this.insertRandomly(topCard)
	discardRequest := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         topCard.DbId,
		RequiresResponse: true,
	}

	dataBytes, err := json.Marshal(discardRequest)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	player.send <- dataBytes
	this.currentTurnData.cardResponseId = topCard.DbId
	this.currentTurnData.cardresponseUserId = player.Id
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleArchitectureMemoryLoss(player *Player, topCard CardModel) {
	this.insertRandomly(topCard)
	discardRequest := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         topCard.DbId,
		RequiresResponse: true,
	}

	playerNotification := GameInfoUpdate{
		MessageType: MessageType.ArchitectureMemoryLoss,
	}
	playerDataBytes, err := json.Marshal(playerNotification)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	player.send <- playerDataBytes

	dataBytes, err := json.Marshal(discardRequest)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	index := this.getNextPlayerIndex(this.currentTurnData.currentPlayerIndex)
	nextPlayer := this.players[index]
	nextPlayer.send <- dataBytes
	this.currentTurnData.cardresponseUserId = nextPlayer.Id
	this.currentTurnData.cardResponseId = topCard.DbId
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) MapCardResponseLogic(player *Player, data CardResponse) {
	switch data.ActionCardId {
	case 25, 26:
		this.handleChangeOfPlansCardResponse(player, data, 2)
	case 5, 6:
		this.handleArchitectureDumpCardResponse(player, data)
	case 38:
		this.handleAssistanceResponse(player, data)
	case 7, 8:
		this.handleArchitectureMemoryLossResponse(player, data)
	case 3, 4:
		this.handleSaveCardSelect(player, data)
	}
}

func (this *gameStateMachine) handleChangeOfPlansCardResponse(player *Player, data CardResponse, cardCount int) {
	if len(data.Cards) != cardCount {
		return
	}

	if this.currentTurnData.actionCancelled {
		return
	}

	cards := map[uuid.UUID]CardModel{}
	for i := 0; i < len(data.Cards); i++ {
		cards[this.DeckPile[i].Id] = this.DeckPile[i]
	}

	for i := 0; i < len(data.Cards); i++ {
		if _, exists := cards[data.Cards[i]]; !exists {
			return
		}
	}

	for i := 0; i < len(data.Cards); i++ {
		this.DeckPile[i] = cards[data.Cards[i]]
	}

	this.currentTurnData.cardResponseId = -1
	this.currentTurnData.cardresponseUserId = uuid.Nil
}

func (this *gameStateMachine) handleArchitectureDumpCardResponse(player *Player, data CardResponse) {
	dbId := player.GetCardDBType(data.Cards[0])
	if dbId == -1 {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.discardCard(card, false)
	this.endTurn()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleAssistanceResponse(player *Player, data CardResponse) {
	if this.currentTurnData.actionCancelled {
		return
	}

	if len(data.Cards) != 1 {
		return
	}

	dbId := player.GetCardDBType(data.Cards[0])
	if dbId == -1 {
		return
	}

	card := player.RemoveCard(data.Cards[0])
	this.players[this.currentTurnData.currentPlayerIndex].AddCardToHand(card)
	this.currentTurnData.cardResponseId = -1
	this.currentTurnData.cardresponseUserId = uuid.Nil
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
	this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
}

func (this *gameStateMachine) handleArchitectureMemoryLossResponse(player *Player, data CardResponse) {
	currentPlayer := this.players[this.currentTurnData.currentPlayerIndex]
	if !currentPlayer.ValidateIndex(data.CardIndex) {
		return
	}

	card := currentPlayer.RemoveCardByIndex(data.CardIndex)
	this.discardCard(card, false)
	this.endTurn()
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) HandleArchitortureCard(player *Player, architortureCard CardModel) bool {
	canSave := false
	for _, card := range player.hand {
		if card.CardTypeId == CardTypeEnum.Save {
			canSave = true
		}
	}

	if !canSave {
		return false
	}

	this.currentTurnData.cardResponseId = architortureCard.DbId
	this.currentTurnData.cardresponseUserId = player.Id
	this.hub.HandleSendArchitortureUpdate(this.roomId, player)
	data := CardActionData{
		MessageType:      MessageType.CardAction,
		CardDbId:         architortureCard.DbId,
		Cards:            []CardModel{architortureCard},
		RequiresResponse: true,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error:", err)
		return true
	}

	player.send <- dataBytes
	return true
}

func (this *gameStateMachine) handleSaveCardSelect(player *Player, data CardResponse) {
	saveCard := player.GetCard(data.Cards[0])
	architortureCard := player.GetCard(data.Cards[1])
	switch saveCard.DbId {
	case 35:
		this.handleTimeExntensionCard(player, saveCard, architortureCard)
	case 1, 2:
		this.handleEurekaCard(player, saveCard, architortureCard)
	case 37:
		this.handleNumbCard(player, saveCard, architortureCard, data.Target)
	}
}

func (this *gameStateMachine) handleEurekaCard(player *Player, saveCard CardModel, architortureCard CardModel) {
	player.RemoveCard(saveCard.Id)
	player.RemoveCard(architortureCard.Id)
	this.insertRandomly(architortureCard)
	this.discardCard(saveCard, true)
	this.currentTurnData.cardResponseId = -1
	this.currentTurnData.cardresponseUserId = uuid.Nil
	if this.currentTurnData.shouldEndTurnAfterArchitortureAction {
		this.endTurn()
	}

	this.hub.HandleSendArchitortureUpdate(this.roomId, nil)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleTimeExntensionCard(player *Player, saveCard CardModel, architortureCard CardModel) {
	architortureCount := 0
	timeExtensionCount := 0
	for _, c := range player.hand {
		if c.DbId == 3 || c.DbId == 4 {
			architortureCount++
		} else if c.DbId == 35 {
			timeExtensionCount++
		}
	}

	if timeExtensionCount < architortureCount {
		player.EliminatePlayer()
	}

	player.AddTimeExtensionUsage(saveCard.Id, architortureCard.Id)
	this.currentTurnData.cardResponseId = -1
	this.currentTurnData.cardresponseUserId = uuid.Nil
	if this.currentTurnData.shouldEndTurnAfterArchitortureAction {
		this.endTurn()
	}

	this.hub.HandleSendArchitortureUpdate(this.roomId, nil)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
}

func (this *gameStateMachine) handleNumbCard(player *Player, saveCard CardModel, architortureCard CardModel, targetId uuid.UUID) {
	var target *Player
	for _, p := range this.players {
		if p.Id == targetId {
			target = p
			break
		}
	}

	if target == nil {
		player.EliminatePlayer()
	}

	if this.currentTurnData.shouldEndTurnAfterArchitortureAction {
		this.endTurn()
	}

	player.RemoveCard(saveCard.Id)
	this.discardCard(saveCard, true)
	player.RemoveCard(architortureCard.Id)
	this.hub.HandleSendArchitortureUpdate(this.roomId, nil)
	this.hub.HandleGameUpdate(this.roomId, MessageType.GameInfoUpdate)
	target.AddCardToHand(architortureCard)
}

func (this *gameStateMachine) handleUndo(player *Player, data UndoCardResponse) {
	if this.undoAction.cardId == -1 {
		return
	}

	if !data.Use {
		if !this.undoAction.actionUndone {
			switch this.undoAction.cardId {
			case 17, 18:
				this.foundAUsbLogic(this.undoAction.source, this.undoAction.target)
			case 19, 31, 32:
				this.vendettaLogic(this.undoAction.source, this.undoAction.target, this.undoAction.cardIds[0])
			case 38:
				this.assistanceLogic(this.undoAction.target, this.undoAction.cardId)
			case 29, 30, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
				this.knowledgeLogic(this.undoAction.source, this.undoAction.target, this.undoAction.numCards, this.undoAction.targetCardIndex, this.undoAction.cardRequestName)
			}
		}

		this.createNewUndoAction(-1)
		this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
		return
	} else if !player.HasCards([]uuid.UUID{data.CardId}) {
		return
	}

	card := player.RemoveCard(data.CardId)
	this.discardCard(card, true)
	if !this.undoAction.actionUndone {
		this.undoAction.actionUndone = true
		if this.undoAction.source.HasNotInThisLifetimeCard() {
			data := UndoCardRequest{
				MessageType: MessageType.Undo,
				Cards:       []CardModel{card},
				Undone:      true,
				Stage:       UndoRequestState.CanUndoUndo,
			}

			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			this.undoAction.source.send <- dataBytes
		} else {
			this.createNewUndoAction(-1)
		}
	} else {
		switch this.undoAction.cardId {
		case 17, 18:
			this.foundAUsbLogic(this.undoAction.source, this.undoAction.target)
		case 19, 31, 32:
			this.vendettaLogic(this.undoAction.source, this.undoAction.target, this.undoAction.cardIds[0])
		case 38:
			this.assistanceLogic(this.undoAction.target, this.undoAction.cardId)
		case 29, 30, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			this.knowledgeLogic(this.undoAction.source, this.undoAction.target, this.undoAction.numCards, this.undoAction.targetCardIndex, this.undoAction.cardRequestName)
		}

		this.createNewUndoAction(-1)
	}

	this.hub.HandleGameUpdate(this.roomId, MessageType.EndLoading)
}
