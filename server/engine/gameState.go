package engine

import (
	"sort"

	lovelove "hanafuda.moe/lovelove/proto"
)

type GameUpdateMap map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart

type gameState struct {
	state        GameState
	id           string
	activePlayer lovelove.PlayerPosition
	oya          lovelove.PlayerPosition
	cards        map[int32]*cardState
	playerState  map[string]*playerState
}

func (game *gameState) getZoneOrdered(cardLocation CardLocation) []*cardState {
	cards := make([]*cardState, 0)
	for _, card := range game.cards {
		if card.location == cardLocation {
			cards = append(cards, card)
		}
	}
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].order < cards[j].order
	})
	return cards
}

func (game *gameState) Deck() []*cardState {
	return game.getZoneOrdered(CardLocation_Deck)
}

func (game *gameState) Hand(playerPosition lovelove.PlayerPosition) []*cardState {
	return game.getZoneOrdered(GetHandLocation(playerPosition))
}

func (game *gameState) DrawnCard() *cardState {
	for _, card := range game.cards {
		if card.location == CardLocation_Drawn {
			return card
		}
	}
	return nil
}

func (game *gameState) Table() []*cardState {
	tableCards := make([]*cardState, 0)
	maxOrder := 0
	for _, card := range game.cards {
		if card.location == CardLocation_Table {
			tableCards = append(tableCards, card)
			if card.order > maxOrder {
				maxOrder = card.order
			}
		}
	}

	table := make([]*cardState, maxOrder+1)

	for _, card := range tableCards {
		table[card.order] = card
	}

	return table
}

func (gameState *gameState) ToCompleteGameState(playerPosition lovelove.PlayerPosition) *lovelove.CompleteGameState {
	zones := make(map[CardLocation][]*cardState)

	for _, card := range gameState.cards {
		zone, zoneFound := zones[card.location]
		if !zoneFound {
			zone = make([]*cardState, 0, 12*4)
		}
		zones[card.location] = append(zone, card)
	}

	completeGameState := &lovelove.CompleteGameState{
		Deck:               0,
		Table:              make([]*lovelove.Card, 0, 12*4),
		Hand:               make([]*lovelove.Card, 0, 12*4),
		Collection:         make([]*lovelove.Card, 0, 12*4),
		OpponentHand:       0,
		OpponentCollection: make([]*lovelove.Card, 0, 12*4),
		Active:             gameState.activePlayer,
		Oya:                gameState.oya,
	}

	for zoneType, zone := range zones {
		sort.SliceStable(zone, func(i, j int) bool {
			return zone[i].order < zone[j].order
		})

		cards := make([]*lovelove.Card, 0, 12*4)
		for _, card := range zone {
			cards = append(cards, card.card)
		}

		switch zoneType {
		case CardLocation_Deck:
			completeGameState.Deck = int32(len(zone))
		case CardLocation_Table:
			completeGameState.Table = make([]*lovelove.Card, zone[len(zone)-1].order+1)
			for _, card := range zone {
				completeGameState.Table[card.order] = card.card
			}
		case CardLocation_RedCollection:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.Collection = cards
			} else {
				completeGameState.OpponentCollection = cards
			}
		case CardLocation_WhiteCollection:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.OpponentCollection = cards
			} else {
				completeGameState.Collection = cards
			}
		case CardLocation_RedHand:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.Hand = cards
			} else {
				completeGameState.OpponentHand = int32(len(zone))
			}
		case CardLocation_WhiteHand:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.OpponentHand = int32(len(zone))
			} else {
				completeGameState.Hand = cards
			}
		case CardLocation_Drawn:
			completeGameState.DeckFlipCard = cards[0]
		}
	}

	completeGameState.Action = gameState.GetActionForPosition(playerPosition)

	return completeGameState
}

func (gameState *gameState) GetActionForPosition(playerPosition lovelove.PlayerPosition) (action *lovelove.PlayerAction) {
	switch gameState.state {
	case GameState_HandCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		action = &lovelove.PlayerAction{
			Type:        lovelove.PlayerActionType_HandCardPlayOpportunity,
			PlayOptions: make(map[int32]*lovelove.PlayOptions),
		}

		for _, tableCard := range gameState.Table() {
			playOptions := make([]int32, 0)
			for _, handCard := range gameState.Hand(playerPosition) {
				if tableCard != nil && WillAccept(tableCard.card, handCard.card) {
					playOptions = append(playOptions, handCard.card.Id)
				}
			}

			if len(playOptions) > 0 {
				action.PlayOptions[tableCard.card.Id] = &lovelove.PlayOptions{
					Options: playOptions,
				}
			}
		}
	case GameState_DeckCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		drawnCard := gameState.DrawnCard()
		if drawnCard == nil {
			return
		}

		action = &lovelove.PlayerAction{
			Type:        lovelove.PlayerActionType_HandCardPlayOpportunity,
			PlayOptions: make(map[int32]*lovelove.PlayOptions),
		}

		for _, tableCard := range gameState.Table() {
			playOptions := make([]int32, 0)
			if tableCard != nil && WillAccept(tableCard.card, drawnCard.card) {
				playOptions = append(playOptions, drawnCard.card.Id)
			}

			if len(playOptions) > 0 {
				action.PlayOptions[tableCard.card.Id] = &lovelove.PlayOptions{
					Options: playOptions,
				}
			}
		}
	}
	return
}

func (game *gameState) applyCardMoves(cardMoves []*cardMove) (updatesMap map[lovelove.PlayerPosition][]*lovelove.CardMoveUpdate) {
	updatesMap = make(map[lovelove.PlayerPosition][]*lovelove.CardMoveUpdate)

	if cardMoves == nil {
		return
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		updatesMap[position] = make([]*lovelove.CardMoveUpdate, 0)
	}

	for _, move := range cardMoves {
		movingCard := game.cards[move.cardId]

		for p, _ := range lovelove.PlayerPosition_name {
			position := lovelove.PlayerPosition(p)

			if !LocationIsVisible(move.destination, position) && !LocationIsVisible(movingCard.location, position) {
				continue
			}

			cardMove := &lovelove.CardMoveUpdate{
				MovedCard: movingCard.card,
				OriginSlot: &lovelove.CardSlot{
					Zone:  movingCard.location.ToPlayerCentricZone(position),
					Index: int32(movingCard.order),
				},
				DestinationSlot: &lovelove.CardSlot{
					Zone:  move.destination.ToPlayerCentricZone(position),
					Index: int32(move.order),
				},
			}

			updatesMap[position] = append(updatesMap[position], cardMove)
		}

		movingCard.location = move.destination
		movingCard.order = move.order
	}

	return updatesMap
}

func (game *gameState) applyGameStateChange(gameStateChange *gameStateChange) (updatesMap map[lovelove.PlayerPosition]*lovelove.ActionUpdate) {
	updatesMap = make(map[lovelove.PlayerPosition]*lovelove.ActionUpdate)

	if gameStateChange == nil {
		return
	}

	game.state = gameStateChange.newState

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		actionUpdate := game.GetActionForPosition(position)
		if actionUpdate != nil {
			updatesMap[position] = &lovelove.ActionUpdate{
				Action: game.GetActionForPosition(position),
			}
		}
	}

	return
}

func (game *gameState) Apply(mutations []*gameStateMutation) GameUpdateMap {
	updatesMap := make(map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart)

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		updatesMap[position] = make([]*lovelove.GameStateUpdatePart, 0)
	}

	for _, mutation := range mutations {
		cardUpdatesMap := game.applyCardMoves(mutation.cardMoves)
		actionUpdatesMap := game.applyGameStateChange(mutation.gameStateChange)

		for p, _ := range lovelove.PlayerPosition_name {
			position := lovelove.PlayerPosition(p)
			updatePart := &lovelove.GameStateUpdatePart{}
			updatesMap[position] = append(updatesMap[position], updatePart)

			if cardUpdate, ok := cardUpdatesMap[position]; ok {
				updatePart.CardMoveUpdates = cardUpdate
			}

			if actionUpdate, ok := actionUpdatesMap[position]; ok {
				updatePart.ActionUpdate = actionUpdate
			}
		}

	}

	return updatesMap
}

func (game *gameState) SendUpdates(gameUpdates []GameUpdateMap) {
	payloadsForPosition := make(map[lovelove.PlayerPosition]*lovelove.GameStateUpdate)
	for _, gameUpdateMap := range gameUpdates {
		for position, gameUpdate := range gameUpdateMap {
			payload, payloadExists := payloadsForPosition[position]
			if !payloadExists {
				payload = &lovelove.GameStateUpdate{
					Updates: make([]*lovelove.GameStateUpdatePart, 0),
				}
				payloadsForPosition[position] = payload
			}

			payload.Updates = append(payload.Updates, gameUpdate...)
		}
	}

	for _, playerState := range game.playerState {
		payload, payloadExists := payloadsForPosition[playerState.position]
		if !payloadExists {
			continue
		}

		for _, listener := range playerState.listeners {
			listener <- payload
		}
	}
}
