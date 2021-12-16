package engine

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"google.golang.org/protobuf/reflect/protoreflect"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) ConnectToGame(context context.Context, request *lovelove.ConnectToGameRequest) (*lovelove.ConnectToGameResponse, error) {
	log.Print(request.RoomId)

	rpcConnMeta := rpc.GetConnectionMeta(context)
	connMeta := GetConnectionMeta(context)

	if len(connMeta.userId) == 0 {
		//TODO: report no user error
		log.Print("User not logged in ", connMeta)
		return &lovelove.ConnectToGameResponse{}, nil
	}

	gameContext := GetGameContext(context)

	if gameContext == nil {
		//TODO: report no user error
		log.Print("No GameContext! ", connMeta)
		return nil, errors.New("No GameContext!")
	}

	// TODO: room change stop listening to other room
	connMeta.roomId = request.RoomId

	game := gameContext.GameState
	if game == nil {
		cards := make(map[int32]*cardState)
		for hana := range lovelove.Hana_name {
			if hana == 0 {
				continue
			}

			for variation := range lovelove.Variation_name {
				if variation == 0 {
					continue
				}

				id := cardIdFromCardDetails(hana, variation)

				cards[id] = &cardState{
					location: CardLocation_Deck,
					order:    len(cards),
					card: &lovelove.Card{
						Id:        id,
						Hana:      lovelove.Hana(hana),
						Variation: lovelove.Variation(variation),
					},
				}
			}
		}

		oya := lovelove.PlayerPosition(rand.Intn(2) + 1)

		game = &gameState{
			state:        GameState_HandCardPlay,
			activePlayer: oya,
			cards:        cards,
			playerState:  make(map[string]*playerState),
			month:        lovelove.Month_January,
			oya:          oya,
		}

		gameContext.GameState = game

		game.playerState[connMeta.userId] = &playerState{
			id:       connMeta.userId,
			position: lovelove.PlayerPosition(rand.Intn(2) + 1),
		}

		game.Deal()

	} else {
		_, playerExists := game.playerState[connMeta.userId]
		if !playerExists {
			newPlayer := &playerState{
				id:       connMeta.userId,
				position: lovelove.PlayerPosition_UnknownPosition,
			}
			game.playerState[connMeta.userId] = newPlayer

		POSITION:
			for p := range lovelove.PlayerPosition_name {
				position := lovelove.PlayerPosition(p)
				if position == lovelove.PlayerPosition_UnknownPosition {
					continue
				}

				for _, playerState := range game.playerState {
					if playerState.position == position {
						continue POSITION
					}
				}

				newPlayer.position = position
				break
			}
		}
	}

	playerState := game.playerState[connMeta.userId]
	if _, ok := gameContext.listeners[playerState.id]; !ok {
		gameContext.listeners[playerState.id] = make([]chan protoreflect.ProtoMessage, 0)
	}
	gameContext.listeners[playerState.id] = append(gameContext.listeners[playerState.id], rpcConnMeta.Messages)
	rpcConnMeta.Closed.DoOnCompleted(func() {
		listeners, ok := gameContext.listeners[playerState.id]
		if !ok {
			return
		}

		if len(listeners) == 1 {
			delete(gameContext.listeners, playerState.id)
		}

		for i, listener := range listeners {
			if listener == rpcConnMeta.Messages {
				gameContext.listeners[playerState.id] = append(listeners[:i], listeners[i+1:]...)
				return
			}
		}
	})

	playerPosition := playerState.position
	return &lovelove.ConnectToGameResponse{
		Position:  playerPosition,
		GameState: game.ToCompleteGameState(playerPosition),
	}, nil
}
