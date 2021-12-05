package engine

import (
	"context"
	"errors"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (game *gameState) Apply(mutations []*gameStateMutation, playerPosition lovelove.PlayerPosition) []*lovelove.GameStateUpdatePart {
	updates := make([]*lovelove.GameStateUpdatePart, 0)

	for _, mutation := range mutations {
		updatePart := &lovelove.GameStateUpdatePart{}
		updates = append(updates, updatePart)
		if mutation.cardMoves != nil {
			updatePart.CardMoveUpdates = make([]*lovelove.CardMoveUpdate, 0)
			for _, move := range mutation.cardMoves {
				movingCard := game.cards[move.cardId]
				updatePart.CardMoveUpdates = append(updatePart.CardMoveUpdates, &lovelove.CardMoveUpdate{
					MovedCard: movingCard.card,
					OriginSlot: &lovelove.CardSlot{
						Zone:  movingCard.location.ToPlayerCentricZone(playerPosition),
						Index: int32(movingCard.order),
					},
					DestinationSlot: &lovelove.CardSlot{
						Zone:  move.destination.ToPlayerCentricZone(playerPosition),
						Index: int32(move.order),
					},
				})

				movingCard.location = move.destination
				movingCard.order = move.order
			}
		}
	}

	return updates
}

func (server loveLoveRpcServer) PlayHandCard(context context.Context, request *lovelove.PlayHandCardRequest) (*lovelove.PlayHandCardResponse, error) {
	log.Print("PlayHandCard")

	if request.HandCard == nil {
		log.Print("No hand card")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	// TODO: deal with missing connection problem?
	rpcConnMeta := rpc.GetConnectionMeta(context)
	connMeta, connMetaFound := server.connectionMeta[rpcConnMeta.ConnId]

	if !connMetaFound || len(connMeta.userId) == 0 {
		log.Print("Player not identified")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	userMeta, userMetaFound := server.userMeta[connMeta.userId]
	if !userMetaFound || len(userMeta.roomId) == 0 {
		log.Print("User not in room")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	game, gameFound := server.games[userMeta.roomId]

	if !gameFound {
		log.Print("Not connected to room")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	playerState, playerStateFound := game.playerState[connMeta.userId]

	if !playerStateFound {
		log.Print("Player not in game")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	mutation, err := PlayHandCard(game, request, playerState)
	if err != nil {
		log.Print(err.Error())
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	updates := game.Apply(mutation, playerState.position)

	update := &lovelove.GameStateUpdate{
		Updates: updates,
	}

	game.updates <- update
	return &lovelove.PlayHandCardResponse{
		Status: lovelove.PlayHandCardResponseCode_Ok,
	}, nil
}

type cardMove struct {
	cardId      int32
	destination CardLocation
	order       int
}

type gameStateMutation struct {
	cardMoves []*cardMove
}

func MoveToTable(game *gameState, cardId int32, cardLocation CardLocation) ([]*gameStateMutation, error) {
	movingCard, movingCardExists := game.cards[cardId]
	if !movingCardExists {
		return nil, errors.New("Card to move is invalid")
	}

	if movingCard.location != cardLocation {
		return nil, errors.New("Moving card is not in correct location")
	}

	return []*gameStateMutation{
		&gameStateMutation{
			cardMoves: []*cardMove{
				&cardMove{
					cardId:      movingCard.card.Id,
					destination: CardLocation_Table,
				},
			},
		},
	}, nil
}

func MoveToPlayOption(
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

	return []*gameStateMutation{
		&gameStateMutation{
			cardMoves: []*cardMove{
				&cardMove{
					cardId:      movingCardId,
					destination: CardLocation_Table,
					order:       destinationCard.order,
				},
			},
		},
		&gameStateMutation{
			cardMoves: []*cardMove{
				&cardMove{
					cardId:      movingCardId,
					destination: playerCollectionLocation,
				},
				&cardMove{
					cardId:      destinationCardId,
					destination: playerCollectionLocation,
				},
			},
		},
	}, nil
}

func GetHandLocation(playerPosition lovelove.PlayerPosition) CardLocation {
	if playerPosition == lovelove.PlayerPosition_White {
		return CardLocation_WhiteHand
	}
	return CardLocation_RedHand
}

func GetCollectionLocation(playerPosition lovelove.PlayerPosition) CardLocation {
	if playerPosition == lovelove.PlayerPosition_White {
		return CardLocation_WhiteCollection
	}
	return CardLocation_RedCollection
}

func PlayHandCard(
	game *gameState,
	request *lovelove.PlayHandCardRequest,
	playerState *playerState,
) ([]*gameStateMutation, error) {
	if game.state != GameState_HandCardPlay {
		return nil, errors.New("Game is in wrong state")
	}

	if game.activePlayer != playerState.position {
		return nil, errors.New("Player is not active")
	}

	playerHandLocation := GetHandLocation(playerState.position)

	if request.TableCard == nil {
		mutation, err := MoveToTable(game, request.HandCard.CardId, playerHandLocation)

		if err != nil {
			return nil, err
		}

		return mutation, nil
	}

	playerCollectionLocation := GetCollectionLocation(playerState.position)

	mutation, err := MoveToPlayOption(game, request.HandCard.CardId, playerHandLocation, request.TableCard.CardId, playerCollectionLocation)

	if err != nil {
		return nil, err
	}

	return mutation, nil
}
