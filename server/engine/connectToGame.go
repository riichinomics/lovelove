package engine

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"google.golang.org/protobuf/proto"
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
		deck := make([]*lovelove.Card, 12*4)

		for hana := range lovelove.Hana_name {
			if hana == 0 {
				continue
			}

			for variation := range lovelove.Variation_name {
				if variation == 0 {
					continue
				}

				id := cardIdFromCardDetails(hana, variation)

				deck[id] = &lovelove.Card{
					Id:        id,
					Hana:      lovelove.Hana(hana),
					Variation: lovelove.Variation(variation),
				}
			}
		}

		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})

		oya := lovelove.PlayerPosition(rand.Intn(1) + 1)

		game = &gameState{
			state:        GameState_HandCardPlay,
			activePlayer: oya,
			cards:        make(map[int32]*cardState),
			playerState:  make(map[string]*playerState),
			oya:          oya,
		}

		gameContext.GameState = game

		game.playerState[connMeta.userId] = &playerState{
			id:        connMeta.userId,
			position:  lovelove.PlayerPosition(rand.Intn(1) + 1),
			listeners: make([]chan proto.Message, 0),
		}

		moveCards(game.cards, deck[0:8], CardLocation_Table)
		moveCards(game.cards, deck[8:16], CardLocation_RedHand)
		moveCards(game.cards, deck[16:24], CardLocation_WhiteHand)
		moveCards(game.cards, deck[24:], CardLocation_Deck)
	} else {
		_, playerExists := game.playerState[connMeta.userId]
		if !playerExists {
			newPlayer := &playerState{
				id:        connMeta.userId,
				position:  lovelove.PlayerPosition_UnknownPosition,
				listeners: make([]chan proto.Message, 0),
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
	playerState.listeners = append(playerState.listeners, rpcConnMeta.Messages)
	rpcConnMeta.Closed.DoOnCompleted(func() {
		for i, listener := range playerState.listeners {
			if listener == rpcConnMeta.Messages {
				playerState.listeners = append(playerState.listeners[:i], playerState.listeners[i+1:]...)
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
