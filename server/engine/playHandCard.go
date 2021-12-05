package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

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

	if game.activePlayer != playerState.position {
		log.Print("Player is not active")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	if game.state != GameState_HandCardPlay {
		log.Print("Game is in wrong state")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	movingCard, movingCardExists := game.cards[request.HandCard.CardId]
	if !movingCardExists {
		log.Print("Card to move is invalid")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	playerHandLocation := CardLocation_RedHand
	if playerState.position == lovelove.PlayerPosition_White {
		playerHandLocation = CardLocation_WhiteHand
	}

	if movingCard.location != playerHandLocation {
		log.Print("Moving card is not in player hand")
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Error,
		}, nil
	}

	if request.TableCard != nil {
		tableCard, tableCardExists := game.cards[request.TableCard.CardId]
		if !tableCardExists {
			log.Print("Card on table doesn't exist")
			return &lovelove.PlayHandCardResponse{
				Status: lovelove.PlayHandCardResponseCode_Error,
			}, nil
		}

		if tableCard.location != CardLocation_Table {
			log.Print("Table card is not on table")
			return &lovelove.PlayHandCardResponse{
				Status: lovelove.PlayHandCardResponseCode_Error,
			}, nil
		}

		if tableCard.card.Hana != movingCard.card.Hana {
			log.Print("Card's suit doesn't match")
			return &lovelove.PlayHandCardResponse{
				Status: lovelove.PlayHandCardResponseCode_Error,
			}, nil
		}

		playerCollectionLocation := CardLocation_RedCollection
		if playerState.position == lovelove.PlayerPosition_White {
			playerCollectionLocation = CardLocation_WhiteCollection
		}

		update := &lovelove.GameStateUpdate{
			Updates: make([]*lovelove.GameStateUpdatePart, 0),
		}

		update.Updates = append(update.Updates, &lovelove.GameStateUpdatePart{
			CardMoveUpdates: []*lovelove.CardMoveUpdate{
				createCardMoveUpdate(movingCard, CardLocation_Table, playerState.position, int32(tableCard.order)),
			},
		})

		movingCard.location = CardLocation_Table

		update.Updates = append(update.Updates, &lovelove.GameStateUpdatePart{
			CardMoveUpdates: []*lovelove.CardMoveUpdate{
				createCardMoveUpdate(movingCard, playerCollectionLocation, playerState.position, 0),
				createCardMoveUpdate(tableCard, playerCollectionLocation, playerState.position, 0),
			},
		})

		tableCard.location = playerCollectionLocation
		movingCard.location = playerCollectionLocation

		deck := game.Deck()
		deckLen := len(deck)
		if deckLen == 0 {
			game.updates <- update
			return &lovelove.PlayHandCardResponse{
				Status: lovelove.PlayHandCardResponseCode_Ok,
			}, nil
		}

		drawnCard := deck[deckLen-1]
		update.Updates = append(update.Updates, &lovelove.GameStateUpdatePart{
			CardMoveUpdates: []*lovelove.CardMoveUpdate{
				createCardMoveUpdate(drawnCard, CardLocation_Drawn, playerState.position, 0),
			},
		})
		drawnCard.location = CardLocation_Drawn

		drawnCardPlayOptions := make([]*cardState, 0)

		for _, card := range game.cards {
			if card.location == CardLocation_Table && card.card.Hana == drawnCard.card.Hana {
				drawnCardPlayOptions = append(drawnCardPlayOptions, card)
			}
		}

		if len(drawnCardPlayOptions) == 0 {
			tablePlayPosition := len(game.Table())
			update.Updates = append(update.Updates, &lovelove.GameStateUpdatePart{
				CardMoveUpdates: []*lovelove.CardMoveUpdate{
					createCardMoveUpdate(drawnCard, CardLocation_Table, playerState.position, int32(tablePlayPosition)),
				},
			})
			drawnCard.location = CardLocation_Table
			drawnCard.order = tablePlayPosition
		} else if len(drawnCardPlayOptions) == 1 {

		}

		game.updates <- update
		return &lovelove.PlayHandCardResponse{
			Status: lovelove.PlayHandCardResponseCode_Ok,
		}, nil
	}

	log.Print("No target")
	return &lovelove.PlayHandCardResponse{
		Status: lovelove.PlayHandCardResponseCode_Error,
	}, nil
}
