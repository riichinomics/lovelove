package engine

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) ConnectToGame(rpcContext context.Context, request *lovelove.ConnectToGameRequest) (response *lovelove.ConnectToGameResponse, rpcError error) {
	rpcError = nil
	response = &lovelove.ConnectToGameResponse{
		Status: lovelove.ConnectToGameResponseCode_ConnectToGameError,
	}

	log.Print("Connecting to room: ", request.RoomId)

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

	game := gameContext.GameState
	if game == nil {
		response = &lovelove.ConnectToGameResponse{
			Status: lovelove.ConnectToGameResponseCode_ConnectToGameWaiting,
		}

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

		testGame, testGameExists := server.testGames[gameContext.id]

		oya := lovelove.PlayerPosition_Red
		if !testGameExists {
			oya = lovelove.PlayerPosition(rand.Intn(2) + 1)
		}

		game = &gameState{
			activePlayer: oya,
			cards:        cards,
			playerState:  make(map[string]*playerState),
			month:        lovelove.Month_January,
			oya:          oya,
		}

		gameContext.GameState = game

		game.playerState[connMeta.userId] = &playerState{
			id:       connMeta.userId,
			position: lovelove.PlayerPosition_Red,
		}

		if testGameExists {
			log.Print("seting up test game ", gameContext.id)
			game.SetupTestGame(testGame)
		} else {
			log.Print("Making New Game")
			game.Deal()
		}
	} else {
		log.Print("Connecting to existing game")
		existingPlayer, playerExists := game.playerState[connMeta.userId]
		if !playerExists {
			if len(game.playerState) >= 2 {
				log.Print("Game full")
				response = &lovelove.ConnectToGameResponse{
					Status: lovelove.ConnectToGameResponseCode_ConnectToGameFull,
				}
				return
			}

			log.Print("Join Waiting Game")

			game.state = GameState_HandCardPlay

			if len(game.GetTeyaku()) > 0 {
				game.state = GameState_Teyaku
			}

			response = &lovelove.ConnectToGameResponse{
				Status: lovelove.ConnectToGameResponseCode_ConnectToGameOk,
			}

			game.playerState[connMeta.userId] = &playerState{
				id:       connMeta.userId,
				position: lovelove.PlayerPosition_White,
			}

			defer gameContext.BroadcastGameStart(lovelove.PlayerPosition_UnknownPosition)
		} else {
			if len(game.playerState) < 2 {
				log.Print("First Player Reconnect")
				response = &lovelove.ConnectToGameResponse{
					Status: lovelove.ConnectToGameResponseCode_ConnectToGameWaiting,
				}
			} else {
				log.Print("Reconnect to game")
				response = &lovelove.ConnectToGameResponse{
					Status: lovelove.ConnectToGameResponseCode_ConnectToGameOk,
				}
				defer gameContext.BroadcastGameStart(existingPlayer.position)
			}
		}
	}

	playerState := game.playerState[connMeta.userId]
	response.PlayerPosition = playerState.position
	response.OpponentDisconnected = gameContext.GetOpponentConnectionStatus(playerState.position)
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
	gameContext.PlayerConnected()

	if connMeta.roomChangedNotify != nil {
		connMeta.roomChangedNotify()
	}

	connMeta.roomId = request.RoomId

	roomChangedContext, roomChangedContextNotify := context.WithCancel(context.Background())
	connMeta.roomChangedNotify = roomChangedContextNotify

	rpcConnMeta.Closed.DoOnCompleted(func() {
		gameContext.requestQueue <- func() {
			gameContext.PlayerLeftRoom(player, rpcConnMeta.Messages)
		}
	}, rxgo.WithContext(roomChangedContext))

	if playerState.position == lovelove.PlayerPosition_UnknownPosition || !playerExisted || len(player.connections) != 1 {
		return
	}

	if player.cancelDisconnect != nil {
		player.cancelDisconnect()
	}
	gameContext.ChangeConnectionStatus(connMeta.userId, true)

	return
}
