package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) PlayHandCard(context context.Context, request *lovelove.PlayHandCardRequest) (response *lovelove.PlayHandCardResponse, rpcError error) {
	response = &lovelove.PlayHandCardResponse{
		Status: lovelove.PlayHandCardResponseCode_Error,
	}
	rpcError = nil

	log.Print("PlayHandCard")

	if request.HandCard == nil {
		log.Print("No hand card")
		return
	}

	// TODO: deal with missing connection problem?
	rpcConnMeta := rpc.GetConnectionMeta(context)
	connMeta, connMetaFound := server.connectionMeta[rpcConnMeta.ConnId]

	if !connMetaFound || len(connMeta.userId) == 0 {
		log.Print("Player not identified")
		return
	}

	if len(connMeta.roomId) == 0 {
		log.Print("User not in room")
		return
	}

	game, gameFound := server.games[connMeta.roomId]

	if !gameFound {
		log.Print("Not connected to room")
		return
	}

	playerState, playerStateFound := game.playerState[connMeta.userId]

	if !playerStateFound {
		log.Print("Player not in game")
		return
	}

	mutation, err := PlayHandCardMutation(game, request, playerState.position)
	if err != nil {
		log.Print(err.Error())
		return
	}

	response = &lovelove.PlayHandCardResponse{
		Status: lovelove.PlayHandCardResponseCode_Ok,
	}

	mutationContext := NewGameMutationContext(game)
	defer mutationContext.BroadcastUpdates()

	mutationContext.Apply(mutation)

	mutation, err = DrawCardMutation(game)

	if err != nil {
		//TODO: report error
		return
	}

	mutationContext.Apply(mutation)

	mutation, err = PlayDrawnCardMutation(game, playerState.position)

	if err != nil {
		return
	}

	mutationContext.Apply(mutation)

	if mutation[0].gameStateChange != nil {

	}

	return
}
