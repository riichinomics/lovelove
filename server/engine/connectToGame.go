package engine

import (
	"context"
	"log"
	"math/rand"

	"google.golang.org/protobuf/proto"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) ConnectToGame(context context.Context, request *lovelove.ConnectToGameRequest) (*lovelove.ConnectToGameResponse, error) {
	log.Print(request.RoomId)

	game, gameFound := server.games[request.RoomId]

	// TODO: deal with missing connection problem?
	rpcConnMeta := rpc.GetConnectionMeta(context)
	connMeta := server.connectionMeta[rpcConnMeta.ConnId]
	userMetaData, userFound := server.userMeta[connMeta.userId]
	if !userFound {
		userMetaData = &userMeta{}
		server.userMeta[connMeta.userId] = userMetaData
	}

	// TODO: room change stop listening to other room
	userMetaData.roomId = request.RoomId

	if len(connMeta.userId) == 0 {
		log.Print("Player not identified")
		return &lovelove.ConnectToGameResponse{}, nil
	}

	if !gameFound {
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
			updates:      make(chan *lovelove.GameStateUpdate),
			listeners:    make([]chan proto.Message, 0),
			state:        GameState_HandCardPlay,
			id:           request.RoomId,
			activePlayer: oya,
			cards:        make(map[int32]*cardState),
			playerState:  make(map[string]*playerState),
			oya:          oya,
		}

		go func() {
			for update := range game.updates {
				log.Print("Update detencted", update)
				for _, listener := range game.listeners {
					log.Print("Sending Update to listener", listener)
					listener <- update
				}
			}
		}()

		game.playerState[connMeta.userId] = &playerState{
			id:       connMeta.userId,
			position: lovelove.PlayerPosition(rand.Intn(1) + 1),
		}

		moveCards(game.cards, deck[0:8], CardLocation_Table)
		moveCards(game.cards, deck[8:16], CardLocation_RedHand)
		moveCards(game.cards, deck[16:24], CardLocation_WhiteHand)
		moveCards(game.cards, deck[24:], CardLocation_Deck)

		server.games[game.id] = game
	} else {
		_, playerExists := game.playerState[connMeta.userId]
		if !playerExists && len(game.playerState) < 2 {
			newPlayerPosition := lovelove.PlayerPosition_Red
			for _, playerState := range game.playerState {
				if playerState.position == lovelove.PlayerPosition_Red {
					newPlayerPosition = lovelove.PlayerPosition_White
				}
			}

			game.playerState[connMeta.userId] = &playerState{
				id:       connMeta.userId,
				position: newPlayerPosition,
			}
		}
	}

	playerPosition := game.playerState[connMeta.userId].position

	game.listeners = append(game.listeners, rpcConnMeta.Messages)
	rpcConnMeta.Closed.DoOnCompleted(func() {
		for i, listener := range game.listeners {
			if listener == rpcConnMeta.Messages {
				game.listeners = append(game.listeners[:i], game.listeners[i+1:]...)
				return
			}
		}
	})

	return &lovelove.ConnectToGameResponse{
		Position:  playerPosition,
		GameState: game.ToCompleteGameState(playerPosition),
	}, nil
}
