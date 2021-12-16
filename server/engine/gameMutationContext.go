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
}

func (context *gameMutationContext) MovedCards() MovedCardsMeta {
	return context.movingCards
}

func NewGameMutationContext(gameState *gameState) (context *gameMutationContext) {
	context = &gameMutationContext{
		gameState:   gameState,
		movingCards: make(map[int32]*cardMoveMeta),
	}

	return
}

type GameUpdateMap map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart

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

func (gameMutationContext *gameMutationContext) applyGameStateChange(gameStateChange *gameStateChange) (
	actionUpdateMap map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate,
	shoubuOpportunityUpdateMap map[lovelove.PlayerPosition]*lovelove.ShoubuOpportunityUpdate,
	activePlayerUpdates map[lovelove.PlayerPosition]*lovelove.ActivePlayerUpdate,
) {
	actionUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate)
	shoubuOpportunityUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.ShoubuOpportunityUpdate)
	activePlayerUpdates = make(map[lovelove.PlayerPosition]*lovelove.ActivePlayerUpdate)

	if gameStateChange == nil {
		return
	}

	previousGameState := gameMutationContext.gameState.state
	gameMutationContext.gameState.state = gameStateChange.newState

	if gameStateChange.activePlayer != lovelove.PlayerPosition_UnknownPosition {
		gameMutationContext.gameState.activePlayer = gameStateChange.activePlayer
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		actionUpdate := gameMutationContext.gameState.GetTablePlayOptions(position)
		if actionUpdate != nil {
			actionUpdateMap[position] = &lovelove.PlayOptionsUpdate{
				UpdatedAcceptedOriginZones: &lovelove.PlayOptionsZoneUpdate{
					Zones: gameMutationContext.gameState.GetPlayOptionsAcceptedOriginZones(position),
				},
			}
		}

		if position == gameMutationContext.gameState.activePlayer {
			if gameStateChange.newState == GameState_ShoubuOpportunity {
				shoubuOpportunityUpdateMap[position] = &lovelove.ShoubuOpportunityUpdate{
					Available: true,
					Value:     gameStateChange.shoubuOpportunityValue,
				}
			} else if previousGameState == GameState_ShoubuOpportunity {
				shoubuOpportunityUpdateMap[position] = &lovelove.ShoubuOpportunityUpdate{
					Available: false,
				}
			}
		}

		activePlayerUpdates[position] = &lovelove.ActivePlayerUpdate{
			Position: gameStateChange.activePlayer,
		}
	}

	return
}

func (gameMutationContext *gameMutationContext) applyKoikoiChanges(koikoiChanges map[lovelove.PlayerPosition]*koikoiChange) (
	koikoiUpdates map[lovelove.PlayerPosition]*lovelove.KoikoiUpdate,
) {
	koikoiUpdates = make(map[lovelove.PlayerPosition]*lovelove.KoikoiUpdate)
	if koikoiChanges == nil {
		return
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		opponentPosition := getOpponentPosition(position)

		for _, player := range gameMutationContext.gameState.playerState {
			koikoiChange, koikoiChangeExists := koikoiChanges[player.position]
			if !koikoiChangeExists {
				continue
			}

			player.koikoi = koikoiChange.koikoiStatus

			if koikoiChange.koikoiStatus && !player.koikoi {
				existingUpdate, updateExists := koikoiUpdates[position]
				if !updateExists {
					existingUpdate = &lovelove.KoikoiUpdate{}
					koikoiUpdates[position] = existingUpdate
				}

				if player.position == position {
					existingUpdate.Self = true
					continue
				}

				if player.position == opponentPosition {
					existingUpdate.Opponent = true
					continue
				}
			}

			player.koikoi = koikoiChange.koikoiStatus
		}
	}

	return
}

func (gameMutationContext *gameMutationContext) applyRoundEndChange(roundEndChange *roundEndChange) (
	roundEndUpdates map[lovelove.PlayerPosition]*lovelove.RoundEndResult,
) {
	roundEndUpdates = make(map[lovelove.PlayerPosition]*lovelove.RoundEndResult)
	if roundEndChange == nil {
		return
	}

	shoubuValue := int32(0)

	if roundEndChange.winner != lovelove.PlayerPosition_UnknownPosition {
		yakuInfo := gameMutationContext.gameState.GetYakuData(roundEndChange.winner)
		shoubuValue = gameMutationContext.gameState.GetShoubuValue(yakuInfo, roundEndChange.winner)

		gameMutationContext.gameState.oya = roundEndChange.winner

		for _, player := range gameMutationContext.gameState.playerState {
			player.koikoi = false
			if player.position != roundEndChange.winner {
				continue
			}

			player.score += shoubuValue
		}
	}

	gameMutationContext.gameState.month = lovelove.Month(int(gameMutationContext.gameState.month+1) % len(lovelove.Month_name))

	gameMutationContext.gameState.activePlayer = gameMutationContext.gameState.oya
	gameMutationContext.gameState.state = GameState_HandCardPlay

	gameMutationContext.gameState.Deal()

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		roundEndUpdates[position] = &lovelove.RoundEndResult{
			Winner:    roundEndChange.winner,
			Winnings:  shoubuValue,
			NextRound: gameMutationContext.gameState.ToCompleteGameState(position),
		}
	}

	return
}

// Updates to game state, different positions see different things because of hidden zones
// The position key determines who the update is for
func (gameMutationContext *gameMutationContext) Apply(mutations []*gameStateMutation) (updatesMap GameUpdateMap) {
	updatesMap = make(GameUpdateMap)
	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		updatesMap[position] = make([]*lovelove.GameStateUpdatePart, 0)
	}

	for _, mutation := range mutations {
		cardUpdatesMap := gameMutationContext.applyCardMoves(mutation.cardMoves)
		actionUpdatesMap, shoubuOpportunityUpdateMap, activePlayerUpdates := gameMutationContext.applyGameStateChange(mutation.gameStateChange)
		koikoiUpdates := gameMutationContext.applyKoikoiChanges(mutation.koikoiChange)
		roundEndUpdatesMap := gameMutationContext.applyRoundEndChange(mutation.roundEndChange)

		for p, _ := range lovelove.PlayerPosition_name {
			position := lovelove.PlayerPosition(p)
			updatePart := &lovelove.GameStateUpdatePart{}
			updatesMap[position] = append(updatesMap[position], updatePart)

			if cardUpdate, ok := cardUpdatesMap[position]; ok {
				updatePart.CardMoveUpdates = cardUpdate
			}

			if actionUpdate, ok := actionUpdatesMap[position]; ok {
				updatePart.PlayOptionsUpdate = actionUpdate
			}

			if shoubuOpportunityUpdate, ok := shoubuOpportunityUpdateMap[position]; ok {
				updatePart.ShoubuOpportunityUpdate = shoubuOpportunityUpdate
			}

			if koikoiUpdate, ok := koikoiUpdates[position]; ok {
				updatePart.KoikoiUpdate = koikoiUpdate
			}

			if activePlayerUpdate, ok := activePlayerUpdates[position]; ok {
				updatePart.ActivePlayerUpdate = activePlayerUpdate
			}

			if roundEndUpdates, ok := roundEndUpdatesMap[position]; ok {
				updatePart.RoundEndResult = roundEndUpdates
			}
		}
	}

	return
}

// Play options are private per player, the position key indicates which player should receive the update
func (context *gameMutationContext) buildPlayOptionsUpdate() (playOptionsUpdateMap map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate) {
	playOptionsUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate)

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		playOptionsUpdateMap[position] = &lovelove.PlayOptionsUpdate{
			DefunctOptions: make([]*lovelove.PlayOption, 0),
			NewOptions:     make([]*lovelove.PlayOption, 0),
		}
	}

	tableCards := context.gameState.Table()

	for _, movingCard := range context.movingCards {
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
		handCards := context.gameState.Hand(position)

		for _, movingCard := range context.movingCards {
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

		for _, movingCard := range context.movingCards {
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

func (context *gameMutationContext) BuildPlayOptions() (updatesMap GameUpdateMap) {
	updatesMap = make(GameUpdateMap)
	playOptionsUpdateMap := context.buildPlayOptionsUpdate()

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		updatePart := &lovelove.GameStateUpdatePart{
			PlayOptionsUpdate: playOptionsUpdateMap[position],
		}

		updatesMap[position] = append(updatesMap[position], updatePart)
	}
	return
}
