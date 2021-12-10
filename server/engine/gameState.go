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

	completeGameState.TablePlayOptions = gameState.GetTablePlayOptions(playerPosition)

	completeGameState.YakuInformation = gameState.GetYakuData(playerPosition)
	completeGameState.OpponentYakuInformation = gameState.GetYakuData(getOpponentPosition(playerPosition))

	return completeGameState
}

type yakuPart struct {
	id         lovelove.YakuId
	cards      []*cardState
	bonusCards []*cardState
}

func CategoryFor(yakuId lovelove.YakuId) YakuCategory {
	yakuDefinition, ok := yakuDefinitions[yakuId]
	if !ok {
		return YakuCategory_Other
	}

	return yakuDefinition.category
}

type yakuDefinition struct {
	id        lovelove.YakuId
	category  YakuCategory
	hasBonus  bool
	minCards  int
	baseValue int
}

var yakuDefinitions = map[lovelove.YakuId]*yakuDefinition{
	lovelove.YakuId_Tane: {
		id:        lovelove.YakuId_Tane,
		category:  YakuCategory_Tane,
		baseValue: 1,
		minCards:  5,
		hasBonus:  true,
	},

	lovelove.YakuId_Inoshikachou: {
		id:        lovelove.YakuId_Inoshikachou,
		category:  YakuCategory_Tane,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},

	lovelove.YakuId_AkatanAotanNoChoufuku: {
		id:        lovelove.YakuId_AkatanAotanNoChoufuku,
		category:  YakuCategory_Tanzaku,
		baseValue: 10,
		minCards:  6,
		hasBonus:  true,
	},
	lovelove.YakuId_Akatan: {
		id:        lovelove.YakuId_Akatan,
		category:  YakuCategory_Tanzaku,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},
	lovelove.YakuId_Aotan: {
		id:        lovelove.YakuId_Akatan,
		category:  YakuCategory_Tanzaku,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},
	lovelove.YakuId_Tanzaku: {
		id:        lovelove.YakuId_Tanzaku,
		category:  YakuCategory_Tanzaku,
		baseValue: 1,
		minCards:  5,
		hasBonus:  true,
	},

	lovelove.YakuId_Gokou: {
		id:        lovelove.YakuId_Gokou,
		category:  YakuCategory_Hikari,
		baseValue: 15,
		minCards:  5,
	},
	lovelove.YakuId_Shikou: {
		id:        lovelove.YakuId_Shikou,
		category:  YakuCategory_Hikari,
		baseValue: 8,
		minCards:  4,
	},
	lovelove.YakuId_Ameshikou: {
		id:        lovelove.YakuId_Ameshikou,
		category:  YakuCategory_Hikari,
		baseValue: 7,
		minCards:  4,
	},
	lovelove.YakuId_Sankou: {
		id:        lovelove.YakuId_Sankou,
		category:  YakuCategory_Hikari,
		baseValue: 6,
		minCards:  3,
	},

	lovelove.YakuId_Tsukifuda: {
		id:        lovelove.YakuId_Tsukifuda,
		baseValue: 4,
		minCards:  4,
	},
	lovelove.YakuId_Tsukimizake: {
		id:        lovelove.YakuId_Tsukimizake,
		baseValue: 5,
		minCards:  2,
	},
	lovelove.YakuId_Hanamizake: {
		id:        lovelove.YakuId_Hanamizake,
		baseValue: 5,
		minCards:  2,
	},
	lovelove.YakuId_Kasu: {
		id:        lovelove.YakuId_Kasu,
		baseValue: 1,
		minCards:  10,
		hasBonus:  true,
	},
}

func (yakuPart *yakuPart) IsComplete() bool {
	yakuDefinition, ok := yakuDefinitions[yakuPart.id]
	if !ok {
		return false
	}

	return len(yakuPart.cards) >= yakuDefinition.minCards
}

func (yakuPart *yakuPart) Value() int32 {
	yakuDefinition, ok := yakuDefinitions[yakuPart.id]
	if !ok {
		return 0
	}

	if !yakuDefinition.hasBonus {
		return int32(yakuDefinition.baseValue)
	}

	return int32(yakuDefinition.baseValue + len(yakuPart.cards) + len(yakuPart.bonusCards) - yakuDefinition.minCards)
}

func (yakuPart *yakuPart) AllCardIds() []int32 {
	cards := make([]int32, 0)
	for _, card := range yakuPart.cards {
		cards = append(cards, card.card.Id)
	}

	for _, card := range yakuPart.bonusCards {
		cards = append(cards, card.card.Id)
	}
	return cards
}

func (yakuPart *yakuPart) ToYakuData() *lovelove.YakuData {
	return &lovelove.YakuData{
		Id:    yakuPart.id,
		Cards: yakuPart.AllCardIds(),
		Value: yakuPart.Value(),
	}
}

func cardIsRainman(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Yanagi && card.Variation == lovelove.Variation_Fourth
}

func cardIsSakeCup(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Kiku && card.Variation == lovelove.Variation_Third
}

func cardIsMoon(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Susuki && card.Variation == lovelove.Variation_Fourth
}

func cardIsSakuraCurtain(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Sakura && card.Variation == lovelove.Variation_Fourth
}

func cardIsBoar(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Hagi && card.Variation == lovelove.Variation_Third
}

func cardIsDeer(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Momiji && card.Variation == lovelove.Variation_Third
}

func cardIsButterfly(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Botan && card.Variation == lovelove.Variation_Third
}

func cardIsAotan(card *lovelove.Card) bool {
	return getCardType(card) == CardType_Tanzaku &&
		(card.Hana == lovelove.Hana_Botan ||
			card.Hana == lovelove.Hana_Kiku ||
			card.Hana == lovelove.Hana_Momiji)
}

func cardIsAkatan(card *lovelove.Card) bool {
	return getCardType(card) == CardType_Tanzaku &&
		(card.Hana == lovelove.Hana_Ume ||
			card.Hana == lovelove.Hana_Sakura ||
			card.Hana == lovelove.Hana_Matsu)
}

type yakuContribution struct {
	yakuId      lovelove.YakuId
	isBonusCard bool
}

func (gameState *gameState) YakuContribution(card *lovelove.Card) []*yakuContribution {
	yaku := make([]*yakuContribution, 0)
	if getMonth(card.Hana) == gameState.month {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Tsukifuda})
	}

	cardType := getCardType(card)

	switch cardType {
	case CardType_Hikari:
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Gokou},
			&yakuContribution{yakuId: lovelove.YakuId_Ameshikou},
		)

		if !cardIsRainman(card) {
			yaku = append(
				yaku,
				&yakuContribution{yakuId: lovelove.YakuId_Sankou},
				&yakuContribution{yakuId: lovelove.YakuId_Shikou},
			)
		}
	case CardType_Tane:
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Tane},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Inoshikachou,
				isBonusCard: !cardIsDeer(card) && !cardIsBoar(card) && !cardIsButterfly(card),
			},
		)
	case CardType_Tanzaku:
		isAkatan := cardIsAkatan(card)
		isAotan := cardIsAotan(card)
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Tanzaku},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Akatan,
				isBonusCard: !isAotan,
			},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Aotan,
				isBonusCard: !isAotan,
			},
			&yakuContribution{
				yakuId:      lovelove.YakuId_AkatanAotanNoChoufuku,
				isBonusCard: !isAotan && !isAkatan,
			},
		)
	case CardType_Kasu:
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Kasu})
	}

	if cardIsSakeCup(card) || cardIsMoon(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Tsukimizake})
	}

	if cardIsSakeCup(card) || cardIsSakuraCurtain(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Hanamizake})
	}

	if cardIsSakeCup(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Kasu})
	}

	return yaku
}

func (gameState *gameState) GetYakuData(playerPosition lovelove.PlayerPosition) []*lovelove.YakuData {
	yakuPartMap := make(map[lovelove.YakuId]*yakuPart)
	collection := gameState.Collection(playerPosition)
	for _, card := range collection {
		contributions := gameState.YakuContribution(card.card)
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
