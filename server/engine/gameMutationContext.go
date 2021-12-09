package engine

import (
	lovelove "hanafuda.moe/lovelove/proto"
)

type cardMoveMeta struct {
	cardState        *cardState
	originalLocation CardLocation
}

type gameMutationContext struct {
	gameState   *gameState
	movingCards map[int32]*cardMoveMeta
	updatesMap  map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart
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

func (gameMutationContext *gameMutationContext) finalisePlayOptions() (playOptionsUpdateMap map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate) {
	playOptionsUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate)

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		playOptionsUpdateMap[position] = &lovelove.PlayOptionsUpdate{
			DefunctOptions: make([]*lovelove.PlayOption, 0),
			NewOptions:     make([]*lovelove.PlayOption, 0),
		}
	}

	tableCards := gameMutationContext.gameState.Table()

	for _, movingCard := range gameMutationContext.movingCards {
		if movingCard.cardState.location != CardLocation_Drawn {
			continue
		}

		foundMatch := false

		for _, tableCard := range tableCards {
			if tableCard == nil || !WillAccept(tableCard.card, movingCard.cardState.card) {
				continue
			}

			foundMatch = true

			for _, playOptions := range playOptionsUpdateMap {
				playOptions.NewOptions = append(
					playOptions.NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
						TargetCardId: &lovelove.CardId{CardId: tableCard.card.Id},
					},
				)
			}
		}

		if !foundMatch {
			for _, playOptions := range playOptionsUpdateMap {
				playOptions.NewOptions = append(
					playOptions.NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
					},
				)
			}
		}

		break
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		handCards := gameMutationContext.gameState.Hand(position)

		for _, movingCard := range gameMutationContext.movingCards {
			if movingCard.cardState.location != CardLocation_Table {
				continue
			}

			for _, tableCard := range tableCards {
				if tableCard == nil {
					continue
				}

				if movingCard.cardState.card.Id == tableCard.card.Id {
					continue
				}

				if !WillAccept(tableCard.card, movingCard.cardState.card) {
					continue
				}

				playOptionsUpdateMap[position].DefunctOptions = append(
					playOptionsUpdateMap[position].DefunctOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
						TargetCardId: &lovelove.CardId{CardId: tableCard.card.Id},
					},
				)
			}

			for _, handCard := range handCards {
				if !WillAccept(movingCard.cardState.card, handCard.card) {
					continue
				}

				wasUnmatched := true

				for _, tableCard := range tableCards {
					if tableCard == nil {
						continue
					}

					if movingCard.cardState.card.Id == tableCard.card.Id {
						continue
					}

					if WillAccept(tableCard.card, handCard.card) {
						wasUnmatched = false
						break
					}
				}

				if wasUnmatched {
					playOptionsUpdateMap[position].DefunctOptions = append(
						playOptionsUpdateMap[position].DefunctOptions,
						&lovelove.PlayOption{
							OriginCardId: &lovelove.CardId{CardId: handCard.card.Id},
						},
					)
				}

				playOptionsUpdateMap[position].NewOptions = append(
					playOptionsUpdateMap[position].NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: handCard.card.Id},
						TargetCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
					},
				)
			}
		}

		for _, movingCard := range gameMutationContext.movingCards {
			if movingCard.originalLocation != CardLocation_Table {
				continue
			}

			playOptionsUpdateMap[position].DefunctOptions = append(
				playOptionsUpdateMap[position].DefunctOptions,
				&lovelove.PlayOption{
					TargetCardId: &lovelove.CardId{
						CardId: movingCard.cardState.card.Id,
					},
				},
			)

			cardsForValidation := make([]*cardState, 0)
			for _, handCard := range handCards {
				if WillAccept(movingCard.cardState.card, handCard.card) {
					cardsForValidation = append(cardsForValidation, handCard)
				}
			}

		CARD_REVALIDATION:
			for _, card := range cardsForValidation {
				for _, tableCard := range tableCards {
					if tableCard != nil && WillAccept(tableCard.card, card.card) {
						continue CARD_REVALIDATION
					}
				}

				playOptionsUpdateMap[position].NewOptions = append(
					playOptionsUpdateMap[position].NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{
							CardId: card.card.Id,
						},
					},
				)
			}
		}
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
				OriginZones: gameMutationContext.gameState.GetPlayOptionsOriginZoneUpdate(position),
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

func (gameMutationContext *gameMutationContext) BroadcastUpdates() {
	playOptionsUpdateMap := gameMutationContext.finalisePlayOptions()
	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)

		updatePart := &lovelove.GameStateUpdatePart{
			PlayOptionsUpdate: playOptionsUpdateMap[position],
		}

		gameMutationContext.updatesMap[position] = append(gameMutationContext.updatesMap[position], updatePart)
	}

	gameMutationContext.gameState.BroadcastUpdates(gameMutationContext.updatesMap)
}
