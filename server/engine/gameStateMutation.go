package engine

import (
	"errors"

	lovelove "hanafuda.moe/lovelove/proto"
)

type cardMove struct {
	cardId      int32
	destination CardLocation
	order       int
}

type gameStateChange struct {
	newState               GameState
	activePlayer           lovelove.PlayerPosition
	shoubuOpportunityValue int32
}

type activePlayerChange struct {
	activePlayer lovelove.PlayerPosition
}

type koikoiChange struct {
	koikoiStatus bool
}

type roundEndChange struct {
	winner lovelove.PlayerPosition
}

type gameStateMutation struct {
	cardMoves       []*cardMove
	gameStateChange *gameStateChange
	koikoiChange    map[lovelove.PlayerPosition]*koikoiChange
	roundEndChange  *roundEndChange
}

func PlayDrawnCardMutation(game *gameState, playerPosition lovelove.PlayerPosition) ([]*gameStateMutation, error) {
	drawnCard := game.DrawnCard()
	if drawnCard == nil {
		return nil, errors.New("No drawn card")
	}

	tableCards := game.Table()
	playOptions := make([]*cardState, 0)
	for _, card := range tableCards {
		if card != nil && WillAccept(card.card, drawnCard.card) {
			playOptions = append(playOptions, card)
		}
	}

	if len(playOptions) == 0 {
		return MoveToTable(game, drawnCard.card.Id, CardLocation_Drawn)
	}

	if len(playOptions) == 1 || len(playOptions) == 3 {
		targetCard := playOptions[0]
		playerCollectionLocation := GetCollectionLocation(playerPosition)
		return MoveToPlayOptionMutation(game, drawnCard.card.Id, CardLocation_Drawn, targetCard.card.Id, playerCollectionLocation)
	}

	return []*gameStateMutation{
		{
			gameStateChange: &gameStateChange{
				newState: GameState_DeckCardPlay,
			},
		},
	}, nil
}

func MoveToTable(game *gameState, cardId int32, cardLocation CardLocation) ([]*gameStateMutation, error) {
	movingCard, movingCardExists := game.cards[cardId]
	if !movingCardExists {
		return nil, errors.New("Card to move is invalid")
	}

	if movingCard.location != cardLocation {
		return nil, errors.New("Moving card is not in correct location")
	}

	table := game.Table()
	targetIndex := len(table)
	for i, card := range table {
		if card == nil {
			targetIndex = i
			break
		}
	}

	return []*gameStateMutation{
		{
			cardMoves: []*cardMove{
				{
					cardId:      cardId,
					destination: CardLocation_Table,
					order:       targetIndex,
				},
			},
		},
	}, nil
}

func MoveToPlayOptionMutation(
	game *gameState,
	movingCardId int32,
	movingCardLocation CardLocation,
	destinationCardId int32,
	playerCollectionLocation CardLocation,
) ([]*gameStateMutation, error) {
	movingCard, movingCardExists := game.cards[movingCardId]
	if !movingCardExists {
		return nil, errors.New("Card to move is invalid")
	}

	if movingCard.location != movingCardLocation {
		return nil, errors.New("Moving card is not in correct location")
	}

	destinationCard, destinationCardExists := game.cards[destinationCardId]
	if !destinationCardExists {
		return nil, errors.New("Destination card doesn't exist")
	}

	if destinationCard.location != CardLocation_Table {
		return nil, errors.New("Destination card is not in the correct location")
	}

	if !WillAccept(destinationCard.card, movingCard.card) {
		return nil, errors.New("Destination card can't accept moving card")
	}

	otherMatches := make([]*cardState, 0)

	for _, tableCard := range game.Table() {
		if tableCard == nil || tableCard.card.Id == destinationCardId {
			continue
		}

		if WillAccept(tableCard.card, movingCard.card) {
			otherMatches = append(otherMatches, tableCard)
		}
	}

	collectMutation := &gameStateMutation{
		cardMoves: []*cardMove{
			{
				cardId:      movingCardId,
				destination: playerCollectionLocation,
			},
			{
				cardId:      destinationCardId,
				destination: playerCollectionLocation,
			},
		},
	}

	if len(otherMatches) == 2 {
		for _, card := range otherMatches {
			collectMutation.cardMoves = append(
				collectMutation.cardMoves,
				&cardMove{
					cardId:      card.card.Id,
					destination: playerCollectionLocation,
				},
			)
		}
	}

	return []*gameStateMutation{
		{
			cardMoves: []*cardMove{
				{
					cardId:      movingCardId,
					destination: CardLocation_Table,
					order:       destinationCard.order,
				},
			},
		},
		collectMutation,
	}, nil
}

func DrawCardMutation(game *gameState) ([]*gameStateMutation, error) {
	deck := game.Deck()
	if len(deck) == 0 {
		return nil, errors.New("No cards left to draw")
	}
	drawnCard := deck[len(deck)-1]

	return []*gameStateMutation{
		{
			cardMoves: []*cardMove{
				{
					cardId:      drawnCard.card.Id,
					destination: CardLocation_Drawn,
				},
			},
		},
	}, nil
}

func PlayHandCardMutation(
	game *gameState,
	request *lovelove.PlayHandCardRequest,
	playerPosition lovelove.PlayerPosition,
) ([]*gameStateMutation, error) {
	if game.state != GameState_HandCardPlay {
		return nil, errors.New("Game is in wrong state")
	}

	if game.activePlayer != playerPosition {
		return nil, errors.New("Player is not active")
	}

	playerHandLocation := GetHandLocation(playerPosition)

	if request.TableCard == nil {
		return MoveToTable(game, request.HandCard.CardId, playerHandLocation)
	}

	playerCollectionLocation := GetCollectionLocation(playerPosition)

	return MoveToPlayOptionMutation(game, request.HandCard.CardId, playerHandLocation, request.TableCard.CardId, playerCollectionLocation)
}

func SelectDrawnCardPlayOptionMutation(
	game *gameState,
	request *lovelove.PlayDrawnCardRequest,
	playerPosition lovelove.PlayerPosition,
) ([]*gameStateMutation, error) {
	if request.TableCard == nil {
		return nil, errors.New("No target card")
	}

	if game.state != GameState_DeckCardPlay {
		return nil, errors.New("Game is in wrong state")
	}

	if game.activePlayer != playerPosition {
		return nil, errors.New("Player is not active")
	}

	playerCollectionLocation := GetCollectionLocation(playerPosition)

	mutation, err := MoveToPlayOptionMutation(game, game.DrawnCard().card.Id, CardLocation_Drawn, request.TableCard.CardId, playerCollectionLocation)
	if err != nil {
		return nil, err
	}

	return mutation, nil
}

func ShoubuOpportunityMutation(
	gameState *gameState,
	playerPosition lovelove.PlayerPosition,
	yakuUpdate *yakuUpdate,
) ([]*gameStateMutation, error) {
	return []*gameStateMutation{
		{
			gameStateChange: &gameStateChange{
				newState:               GameState_ShoubuOpportunity,
				shoubuOpportunityValue: gameState.GetShoubuValue(yakuUpdate.completeYakuInfo, playerPosition),
			},
		},
	}, nil
}

func RoundEndMutation(winner lovelove.PlayerPosition) ([]*gameStateMutation, error) {
	return []*gameStateMutation{
		{
			roundEndChange: &roundEndChange{
				winner: winner,
			},
		},
	}, nil
}

func TurnEndMutation(
	game *gameState,
) ([]*gameStateMutation, error) {
	if len(game.Hand(lovelove.PlayerPosition_Red)) == 0 && len(game.Hand(lovelove.PlayerPosition_White)) == 0 {
		return RoundEndMutation(lovelove.PlayerPosition_UnknownPosition)
	}

	return []*gameStateMutation{
		{
			gameStateChange: &gameStateChange{
				newState:     GameState_HandCardPlay,
				activePlayer: getOpponentPosition(game.activePlayer),
			},
		},
	}, nil
}

func CheckForYakuMutation(yakuUpdates *yakuUpdate, playerState *playerState, gameState *gameState) (mutation []*gameStateMutation, err error) {
	mutation = make([]*gameStateMutation, 0)
	_, hasYakuUpdate := yakuUpdates.yakuUpdatesMap[playerState.position]
	if !hasYakuUpdate {
		mutation, err = TurnEndMutation(gameState)
		return
	}

	mutation, err = ShoubuOpportunityMutation(gameState, playerState.position, yakuUpdates)
	return
}
