package engine

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) ConnectToGame(rpcContext context.Context, request *lovelove.ConnectToGameRequest) (response *lovelove.ConnectToGameResponse, rpcError error) {
	rpcError = nil

	log.Print(request.RoomId)

	rpcConnMeta := rpc.GetConnectionMeta(rpcContext)
	connMeta := GetConnectionMeta(rpcContext)

	if len(connMeta.userId) == 0 {
		//TODO: report no user error
		log.Print("User not logged in ", connMeta)
		return &lovelove.ConnectToGameResponse{}, nil
	}

	gameContext := GetGameContext(rpcContext)

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

		if len(game.GetTeyaku()) > 0 {
			game.state = GameState_Teyaku
		}

	} else {
		_, playerExists := game.playerState[connMeta.userId]
		if !playerExists {
			newPlayer := &playerState{
				id: connMeta.userId,
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
	player, playerExisted := gameContext.players[playerState.id]
	if !playerExisted {
		player = &playerMeta{
			position:    playerState.position,
			connections: make([]chan protoreflect.ProtoMessage, 0),
			id:          connMeta.userId,
		}
		gameContext.players[playerState.id] = player
	}

	player.connections = append(player.connections, rpcConnMeta.Messages)

	rpcConnMeta.Closed.DoOnCompleted(func() {
		gameContext.requestQueue <- func() {
			for i, listener := range player.connections {
				if listener == rpcConnMeta.Messages {
					player.connections = append(player.connections[:i], player.connections[i+1:]...)
					break
				}
			}

			if len(player.connections) != 0 {
				return
			}

			disconnectedContext, cancel := context.WithCancel(context.Background())
			player.cancelDisconnect = cancel
			go func() {
				select {
				case <-disconnectedContext.Done():
					return
				case <-time.After(5 * time.Second):
					gameContext.requestQueue <- func() {
						if len(player.connections) != 0 {
							return
						}
						gameContext.ChangeConnectionStatus(connMeta.userId, false)
					}
				}
			}()
		}
	})

	opponentDisconnected := false
	if player.position != lovelove.PlayerPosition_UnknownPosition {
		opponentPosition := getOpponentPosition(player.position)
		for _, player := range gameContext.players {
			if player.position == opponentPosition {
				opponentDisconnected = len(player.connections) == 0
				break
			}
		}
	}

	response = &lovelove.ConnectToGameResponse{
		Position:             playerState.position,
		GameState:            game.ToCompleteGameState(playerState.position),
		OpponentDisconnected: opponentDisconnected,
	}

	if playerState.position == lovelove.PlayerPosition_UnknownPosition || !playerExisted || len(player.connections) != 1 {
		return
	}

	if player.cancelDisconnect != nil {
		player.cancelDisconnect()
	}
	gameContext.ChangeConnectionStatus(connMeta.userId, true)

	return
}
