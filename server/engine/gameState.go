package engine

import (
	"sort"

	"google.golang.org/protobuf/proto"
	lovelove "hanafuda.moe/lovelove/proto"
)

type gameState struct {
	updates   chan *lovelove.GameStateUpdate
	listeners []chan proto.Message

	state        GameState
	id           string
	activePlayer lovelove.PlayerPosition
	oya          lovelove.PlayerPosition
	cards        map[int32]*cardState
	playerState  map[string]*playerState
}

func (game *gameState) Deck() []*cardState {
	deck := make([]*cardState, 0)
	for _, card := range game.cards {
		if card.location == CardLocation_Deck {
			deck = append(deck, card)
		}
	}
	sort.SliceStable(deck, func(i, j int) bool {
		return deck[i].order < deck[j].order
	})
	return deck
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

	if gameState.activePlayer == playerPosition && gameState.state == GameState_HandCardPlay {
		completeGameState.Action = &lovelove.PlayerAction{
			Type:        lovelove.PlayerActionType_HandCardPlayOpportunity,
			PlayOptions: make(map[int32]*lovelove.PlayOptions),
		}
		for _, tableCard := range completeGameState.Table {
			playOptions := make([]int32, 0)
			for _, handCard := range completeGameState.Hand {
				if tableCard != nil && WillAccept(tableCard, handCard) {
					playOptions = append(playOptions, handCard.Id)
				}
			}

			if len(playOptions) > 0 {
				completeGameState.Action.PlayOptions[tableCard.Id] = &lovelove.PlayOptions{
					Options: playOptions,
				}
			}
		}
	}

	return completeGameState
}
