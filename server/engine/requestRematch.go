package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) RequestRematch(context context.Context, request *lovelove.RequestRematchRequest) (response *lovelove.RequestRematchResponse, rpcError error) {
	response = &lovelove.RequestRematchResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	log.Print("RequestRematch")

	connMeta := GetConnectionMeta(context)
	gameContext := GetGameContext(context)

	if len(connMeta.userId) == 0 {
		log.Print("Player not identified")
		return
	}

	if gameContext == nil || gameContext.GameState == nil {
		log.Print("Not connected to room")
		return
	}

	playerState, playerStateFound := gameContext.GameState.playerState[connMeta.userId]

	if !playerStateFound {
		log.Print("Player not in game")
		return
	}

	if gameContext.GameState.state != GameState_End {
		log.Print("Game in wrong state")
		return
	}

	if playerState.requestedRematch {
		log.Print("Player Already Requested Rematch")
		return
	}

	response = &lovelove.RequestRematchResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	defer broadcastBuilder.Broadcast()

	broadcastBuilder.QueueUpdates(RequestRematch(gameContext.GameState, playerState))
	return
}
