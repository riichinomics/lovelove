package engine

import (
	lovelove "hanafuda.moe/lovelove/proto"
)

type cardMoveMeta struct {
	cardState        *cardState
	originalLocation CardLocation
}

type MovedCardsMeta map[int32]*cardMoveMeta

type gameMutationContext struct {
	gameState   *gameState
	movingCards MovedCardsMeta
	updatesMap  map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart
}

func (context *gameMutationContext) MovedCards() MovedCardsMeta {
	return context.movingCards
}

// Updates to game state, different positions see different things because of hidden zones
// The position key determines who the update is for
func (context *gameMutationContext) UpdatesByPosition() map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart {
	return context.updatesMap
}

func NewGameMutationContext(gameState *gameState) (context *gameMutationContext) {
	context = &gameMutationContext{
		gameState:   gameState,
		movingCards: make(map[int32]*cardMoveMeta),
		updatesMap:  make(map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart),
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		context.updatesMap[position] = make([]*lovelove.GameStateUpdatePart, 0)
	}

	return
}

func (gameMutationContext *gameMutationContext) applyCardMoves(cardMoves []*cardMove) (
	cardMoveUpdateMap map[lovelove.PlayerPosition][]*lovelove.CardMoveUpdate,
) {
	cardMoveUpdateMap = make(map[lovelove.PlayerPosition][]*lovelove.CardMoveUpdate)

	if cardMoves == nil {
		return
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		cardMoveUpdateMap[position] = make([]*lovelove.CardMoveUpdate, 0)
	}

	for _, move := range cardMoves {
		movingCard := gameMutationContext.gameState.cards[move.cardId]

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

			cardMoveUpdateMap[position] = append(cardMoveUpdateMap[position], cardMove)
		}

		if _, cardAlreadyMoved := gameMutationContext.movingCards[movingCard.card.Id]; !cardAlreadyMoved {
			gameMutationContext.movingCards[movingCard.card.Id] = &cardMoveMeta{
				cardState:        movingCard,
				originalLocation: movingCard.location,
			}
		}

		movingCard.location = move.destination
		movingCard.order = move.order
	}
	return
}

func (gameMutationContext *gameMutationContext) applyGameStateChange(gameStateChange *gameStateChange) (updatesMap map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate) {
	updatesMap = make(map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate)

	if gameStateChange == nil {
		return
	}

	gameMutationContext.gameState.state = gameStateChange.newState

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		actionUpdate := gameMutationContext.gameState.GetTablePlayOptions(position)
		if actionUpdate != nil {
			updatesMap[position] = &lovelove.PlayOptionsUpdate{
				UpdatedAcceptedOriginZones: &lovelove.PlayOptionsZoneUpdate{
					Zones: gameMutationContext.gameState.GetPlayOptionsAcceptedOriginZones(position),
				},
			}
		}
	}

	return
}

func (gameMutationContext *gameMutationContext) Apply(mutations []*gameStateMutation) {
	for _, mutation := range mutations {
		cardUpdatesMap := gameMutationContext.applyCardMoves(mutation.cardMoves)
		actionUpdatesMap := gameMutationContext.applyGameStateChange(mutation.gameStateChange)

		for p, _ := range lovelove.PlayerPosition_name {
			position := lovelove.PlayerPosition(p)
			updatePart := &lovelove.GameStateUpdatePart{}
			gameMutationContext.updatesMap[position] = append(gameMutationContext.updatesMap[position], updatePart)

			if cardUpdate, ok := cardUpdatesMap[position]; ok {
				updatePart.CardMoveUpdates = cardUpdate
			}

			if actionUpdate, ok := actionUpdatesMap[position]; ok {
				updatePart.PlayOptionsUpdate = actionUpdate
			}
		}
	}
}
