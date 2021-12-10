package engine

import (
	"sort"

	lovelove "hanafuda.moe/lovelove/proto"
)

type gameState struct {
	state        GameState
	activePlayer lovelove.PlayerPosition
	oya          lovelove.PlayerPosition
	month        lovelove.Month
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

func (game *gameState) Collection(playerPosition lovelove.PlayerPosition) []*cardState {
	return game.getZoneOrdered(GetCollectionLocation(playerPosition))
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
		Deck:         0,
		OpponentHand: 0,
		Active:       gameState.activePlayer,
		Oya:          gameState.oya,
	}

	for zoneType, zone := range zones {
		sort.SliceStable(zone, func(i, j int) bool {
			return zone[i].order < zone[j].order
		})

		cards := make([]*lovelove.Card, 0)
		for _, card := range zone {
			cards = append(cards, card.card)
		}

		switch zoneType {
		case CardLocation_Deck:
			completeGameState.Deck = int32(len(zone))
		case CardLocation_Table:
			completeGameState.Table = make([]*lovelove.CardMaybe, zone[len(zone)-1].order+1)
			for _, card := range zone {
				completeGameState.Table[card.order] = &lovelove.CardMaybe{
					Card: card.card,
				}
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

	completeGameState.TablePlayOptions = gameState.GetTablePlayOptions(playerPosition)

	completeGameState.YakuInformation = gameState.GetYakuData(playerPosition)
	completeGameState.OpponentYakuInformation = gameState.GetYakuData(getOpponentPosition(playerPosition))

	return completeGameState
}

func (gameState *gameState) GetYakuData(playerPosition lovelove.PlayerPosition) []*lovelove.YakuData {
	yakuPartMap := make(map[lovelove.YakuId]*yakuPart)
	collection := gameState.Collection(playerPosition)
	for _, card := range collection {
		contributions := YakuContribution(card.card, gameState)
		for _, contribution := range contributions {
			yaku, ok := yakuPartMap[contribution.yakuId]
			if !ok {
				yaku = &yakuPart{
					id:         contribution.yakuId,
					cards:      make([]*cardState, 0),
					bonusCards: make([]*cardState, 0),
				}
				yakuPartMap[contribution.yakuId] = yaku
			}
			if contribution.isBonusCard {
				yaku.bonusCards = append(yaku.bonusCards, card)
				continue
			}
			yaku.cards = append(yaku.cards, card)
		}
	}

	yakuData := make([]*lovelove.YakuData, 0)
	yakuCategoryMap := make(map[YakuCategory]*lovelove.YakuData)
	for _, yakuPart := range yakuPartMap {
		if !yakuPart.IsComplete() {
			continue
		}

		category := CategoryFor(yakuPart.id)
		if category == YakuCategory_Other {
			yakuData = append(yakuData, yakuPart.ToYakuData())
			continue
		}

		existingYaku, yakuExists := yakuCategoryMap[category]

		if !yakuExists {
			yakuCategoryMap[category] = yakuPart.ToYakuData()
			continue
		}

		if existingYaku.Value < yakuPart.Value() {
			yakuCategoryMap[category] = yakuPart.ToYakuData()
			continue
		}
	}

	for _, yaku := range yakuCategoryMap {
		yakuData = append(yakuData, yaku)
	}

	return yakuData
}

func (gameState *gameState) GetPlayOptionsOriginZoneUpdate(playerPosition lovelove.PlayerPosition) (originZones []lovelove.PlayerCentricZone) {
	originZones = make([]lovelove.PlayerCentricZone, 0)
	switch gameState.state {
	case GameState_HandCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		return []lovelove.PlayerCentricZone{
			GetHandLocation(playerPosition).ToPlayerCentricZone(playerPosition),
		}

	case GameState_DeckCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		return []lovelove.PlayerCentricZone{
			lovelove.PlayerCentricZone_Drawn,
		}
	}

	return
}

func (gameState *gameState) PlayableCards(playerPosition lovelove.PlayerPosition) []*cardState {
	drawnCard := gameState.DrawnCard()
	playableCards := gameState.Hand(playerPosition)
	if drawnCard != nil {
		playableCards = append(playableCards, drawnCard)
	}
	return playableCards
}

func (gameState *gameState) GetTablePlayOptions(playerPosition lovelove.PlayerPosition) (action *lovelove.ZonePlayOptions) {
	action = &lovelove.ZonePlayOptions{
		AcceptedOriginZones: gameState.GetPlayOptionsOriginZoneUpdate(playerPosition),
		PlayOptions:         make(map[int32]*lovelove.PlayOptions),
		NoTargetPlayOptions: &lovelove.PlayOptions{
			Options: make([]int32, 0),
		},
	}

	playableCards := gameState.PlayableCards(playerPosition)

	tableCards := gameState.Table()

	for _, playable := range playableCards {
		foundMatch := false

		for _, tableCard := range tableCards {
			if tableCard != nil && WillAccept(tableCard.card, playable.card) {
				foundMatch = true
				tableCardOptions, tableCardOptionsExist := action.PlayOptions[tableCard.card.Id]
				if !tableCardOptionsExist {
					tableCardOptions = &lovelove.PlayOptions{
						Options: make([]int32, 0),
					}
					action.PlayOptions[tableCard.card.Id] = tableCardOptions
				}
				tableCardOptions.Options = append(tableCardOptions.Options, playable.card.Id)
			}
		}

		if !foundMatch {
			action.NoTargetPlayOptions.Options = append(action.NoTargetPlayOptions.Options, playable.card.Id)
		}
	}
	return
}

func (game *gameState) BroadcastUpdates(gameUpdates map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart) {
	for _, playerState := range game.playerState {
		updates, updatesExist := gameUpdates[playerState.position]
		if !updatesExist {
			continue
		}

		payload := &lovelove.GameStateUpdate{
			Updates: updates,
		}

		for _, listener := range playerState.listeners {
			listener <- payload
		}
	}
}
