package engine

import (
	"math"
	"math/rand"
	"sort"

	lovelove "hanafuda.moe/lovelove/proto"
)

type playerState struct {
	id               string
	score            int32
	koikoi           bool
	conceded         bool
	position         lovelove.PlayerPosition
	requestedRematch bool
	confirmedTeyaku  bool
}

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

	if len(tableCards) == 0 {
		return make([]*cardState, 0)
	}

	table := make([]*cardState, maxOrder+1)

	for _, card := range tableCards {
		table[card.order] = card
	}

	return table
}

func (gameState *gameState) ToCompleteGameState(playerPosition lovelove.PlayerPosition) (completeGameState *lovelove.CompleteGameState) {
	completeGameState = &lovelove.CompleteGameState{
		Active:    gameState.activePlayer,
		Oya:       gameState.oya,
		MonthHana: monthToHana(gameState.month),
		Month:     gameState.month,
		RedPlayer: &lovelove.PlayerState{
			Hand: &lovelove.HandInformation{},
		},
		WhitePlayer: &lovelove.PlayerState{
			Hand: &lovelove.HandInformation{},
		},
	}

	zones := make(map[CardLocation][]*cardState)

	for _, card := range gameState.cards {
		zone, zoneFound := zones[card.location]
		if !zoneFound {
			zone = make([]*cardState, 0, 12*4)
		}
		zones[card.location] = append(zone, card)
	}

	completeGameState.Teyaku = GetTeyaku(gameState.Hand(playerPosition))

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
			completeGameState.RedPlayer.Collection = cards
		case CardLocation_WhiteCollection:
			completeGameState.WhitePlayer.Collection = cards
		case CardLocation_RedHand:
			completeGameState.RedPlayer.Hand.NumberOfCards = int32(len(zone))
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.RedPlayer.Hand.Cards = cards
			}
		case CardLocation_WhiteHand:
			completeGameState.WhitePlayer.Hand.NumberOfCards = int32(len(zone))
			if playerPosition == lovelove.PlayerPosition_White {
				completeGameState.WhitePlayer.Hand.Cards = cards
			}
		case CardLocation_Drawn:
			completeGameState.DeckFlipCard = cards[0]
		}
	}

	completeGameState.TablePlayOptions = gameState.GetTablePlayOptions(playerPosition)

	completeGameState.RedPlayer.YakuInformation = gameState.GetYakuData(lovelove.PlayerPosition_Red)
	completeGameState.WhitePlayer.YakuInformation = gameState.GetYakuData(lovelove.PlayerPosition_White)

	for _, playerState := range gameState.playerState {
		var player *lovelove.PlayerState

		switch playerState.position {
		case lovelove.PlayerPosition_Red:
			player = completeGameState.RedPlayer
		case lovelove.PlayerPosition_White:
			player = completeGameState.WhitePlayer
		}

		if player == nil {
			continue
		}

		player.Score = playerState.score
		player.Koikoi = playerState.koikoi
		player.Conceded = playerState.conceded
		player.RematchRequested = playerState.requestedRematch
	}

	if gameState.state == GameState_End {
		var gameWinner *playerState
		for _, player := range gameState.playerState {
			if player.position == lovelove.PlayerPosition_UnknownPosition {
				continue
			}

			if player.conceded {
				continue
			}

			if gameWinner == nil {
				gameWinner = player
				continue
			}

			if gameWinner.score == player.score {
				gameWinner = nil
				break
			}

			if gameWinner.score < player.score {
				gameWinner = player
			}
		}

		winningPosition := lovelove.PlayerPosition_UnknownPosition
		if gameWinner != nil {
			winningPosition = gameWinner.position
		}

		completeGameState.GameEnd = &lovelove.GameEnd{
			GameWinner: winningPosition,
		}
		return
	}

	if gameState.state == GameState_ShoubuOpportunity && playerPosition == gameState.activePlayer && playerPosition != lovelove.PlayerPosition_UnknownPosition {
		var player *lovelove.PlayerState

		switch playerPosition {
		case lovelove.PlayerPosition_Red:
			player = completeGameState.RedPlayer
		case lovelove.PlayerPosition_White:
			player = completeGameState.WhitePlayer
		}

		if player == nil {
			return
		}

		completeGameState.ShoubuOpportunity = &lovelove.ShoubuOpportunity{
			Value: gameState.GetShoubuValue(player.YakuInformation, playerPosition),
		}
	}

	return
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

func (gameState *gameState) GetPlayOptionsAcceptedOriginZones(playerPosition lovelove.PlayerPosition) (originZones []lovelove.CardZone) {
	originZones = make([]lovelove.CardZone, 0)
	switch gameState.state {
	case GameState_HandCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		return []lovelove.CardZone{
			GetHandLocation(playerPosition).ToCardZone(),
		}

	case GameState_DeckCardPlay:
		if gameState.activePlayer != playerPosition {
			return
		}

		return []lovelove.CardZone{
			lovelove.CardZone_Drawn,
		}
	case GameState_DeckCardPlayYakuAttained:
		if gameState.activePlayer != playerPosition {
			return
		}

		return []lovelove.CardZone{
			lovelove.CardZone_Drawn,
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
		AcceptedOriginZones: gameState.GetPlayOptionsAcceptedOriginZones(playerPosition),
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

func (gameState *gameState) GetShoubuValue(yakuData []*lovelove.YakuData, playerPosition lovelove.PlayerPosition) (value int32) {
	for _, yaku := range yakuData {
		value += yaku.Value
	}

	if value >= 7 {
		value *= 2
	}

	opponentPosition := getOpponentPosition(playerPosition)
	for _, player := range gameState.playerState {
		if player.position == opponentPosition {
			if player.koikoi {
				value *= 2
			}
			break
		}
	}

	return
}

func (gameState *gameState) ShuffleDeck() {
	deck := gameState.Deck()

	rand.Shuffle(len(deck), func(i, j int) {
		deck[i].order = j
		deck[j].order = i
	})
}

func (gameState *gameState) MoveAllCardsToDeck() {
	order := 0
	for _, card := range gameState.cards {
		card.location = CardLocation_Deck
		card.order = order
		order++
	}
}

func (gameState *gameState) DrawCards(cardsToDraw int, location CardLocation) (cards []*cardState) {
	cards = make([]*cardState, 0)

	if location == CardLocation_Deck {
		return
	}

	deck := gameState.Deck()
	draw := int(math.Min(float64(len(deck)), float64(cardsToDraw)))
	cards = deck[(len(deck) - draw):]

	order := len(gameState.getZoneOrdered(location))

	for _, card := range cards {
		card.location = location
		card.order = order
		order++
	}

	return
}

func (gameState *gameState) Deal() {
	for {
		gameState.MoveAllCardsToDeck()
		gameState.ShuffleDeck()
		table := gameState.DrawCards(8, CardLocation_Table)
		if GetTeyaku(table) != lovelove.TeyakuId_UnknownTeyaku {
			continue
		}

		gameState.DrawCards(8, CardLocation_RedHand)
		gameState.DrawCards(8, CardLocation_WhiteHand)
		break
	}
}

func (card testGameCard) ProperHana() lovelove.Hana {
	month := stringToMonth(card.Month)
	if month == lovelove.Month_UnknownMonth {
		return stringToHana(card.Hana)
	}

	return monthToHana(month)
}

func (gameState *gameState) moveTestCards(cards []testGameCard, zone CardLocation) {
	for order, testCard := range cards {
		if testCard.Variation < 1 || testCard.Variation > 4 {
			continue
		}

		properHana := testCard.ProperHana()
		if properHana == lovelove.Hana_UnknownSeason {
			continue
		}

		card, cardExists := gameState.cards[cardIdFromCardDetails(int32(properHana), int32(testCard.Variation))]
		if !cardExists {
			continue
		}
		card.location = zone
		card.order = order
	}
}

func (gameState *gameState) SetupTestGame(testGame testGame) {
	for {
		for _, card := range gameState.cards {
			card.location = CardLocation_Unknown
			card.order = 0
		}

		gameState.moveTestCards(testGame.Deck, CardLocation_Deck)
		gameState.moveTestCards(testGame.Table, CardLocation_Table)
		gameState.moveTestCards(testGame.RedHand, CardLocation_RedHand)
		gameState.moveTestCards(testGame.WhiteHand, CardLocation_WhiteHand)

		cards := gameState.getZoneOrdered(CardLocation_Unknown)
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})

		deckSize := len(testGame.Deck)

		for order, card := range cards {
			card.location = CardLocation_Deck
			card.order = deckSize + order
		}

		for order, card := range gameState.getZoneOrdered(CardLocation_Unknown) {
			card.location = CardLocation_Deck
			card.order = deckSize + order
		}

		gameState.getZoneOrdered(CardLocation_Deck)

		gameState.DrawCards(8-len(gameState.Table()), CardLocation_Table)
		if GetTeyaku(gameState.Table()) != lovelove.TeyakuId_UnknownTeyaku {
			continue
		}

		gameState.DrawCards(8-len(gameState.Hand(lovelove.PlayerPosition_Red)), CardLocation_RedHand)
		gameState.DrawCards(8-len(gameState.Hand(lovelove.PlayerPosition_White)), CardLocation_WhiteHand)

		deck := gameState.Deck()

		for i, j := 0, len(deck)-1; i < j; i, j = i+1, j-1 {
			deck[i].order = j
			deck[j].order = i
		}

		break
	}
}

type teyakuInformation struct {
	confirmed bool
	id        lovelove.TeyakuId
}

func (gameState *gameState) GetTeyaku() (teyakuMap map[lovelove.PlayerPosition]*teyakuInformation) {
	teyakuMap = make(map[lovelove.PlayerPosition]*teyakuInformation)
	for p := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		if position == lovelove.PlayerPosition_UnknownPosition {
			continue
		}

		confirmed := false
		for _, player := range gameState.playerState {
			if player.position == position {
				confirmed = player.confirmedTeyaku
				break
			}
		}

		teyaku := GetTeyaku(gameState.Hand(position))
		if confirmed || teyaku != lovelove.TeyakuId_UnknownTeyaku {
			teyakuMap[position] = &teyakuInformation{
				confirmed,
				teyaku,
			}
		}
	}
	return
}
