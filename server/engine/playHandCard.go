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

	mutation, err := PlayHandCard(game, request, playerState.position)
	if err != nil {
		log.Print(err.Error())
		return
	}

	response = &lovelove.PlayHandCardResponse{
		Status: lovelove.PlayHandCardResponseCode_Ok,
	}

	updates := make([]GameUpdateMap, 0)
	defer func() {
		game.SendUpdates(updates)
	}()

	updates = append(updates, game.Apply(mutation))

	mutation, err = DrawCard(game)

	if err != nil {
		//TODO: report error
		return
	}

	updates = append(updates, game.Apply(mutation))

	mutation, err = PlayDrawnCard(game, playerState.position)

	if err != nil {
		return
	}

	updates = append(updates, game.Apply(mutation))

	if mutation[0].gameStateChange != nil {

	}

	return
}
