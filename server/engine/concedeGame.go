package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) ConcedeGame(context context.Context, request *lovelove.ConcedeGameRequest) (response *lovelove.ConcedeGameResponse, rpcError error) {
	response = &lovelove.ConcedeGameResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	log.Print("ConcedeGame")

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

	if gameContext.GameState.state == GameState_End || gameContext.GameState.state == GameState_Waiting {
		log.Print("Game in wrong state")
		return
	}

	response = &lovelove.ConcedeGameResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	defer broadcastBuilder.Broadcast()

	broadcastBuilder.QueueUpdates(ConcedeGameChange(gameContext.GameState, playerState))
	return
}
