package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) PlayDrawnCard(context context.Context, request *lovelove.PlayDrawnCardRequest) (response *lovelove.PlayDrawnCardResponse, rpcError error) {
	response = &lovelove.PlayDrawnCardResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	log.Print("PlayDrawnCard")

	if request.TableCard == nil {
		log.Print("No target card")
		return
	}

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

	mutation, err := SelectDrawnCardPlayOptionMutation(gameContext.GameState, request, playerState.position)
	if err != nil {
		log.Print(err.Error())
		return
	}

	response = &lovelove.PlayDrawnCardResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	mutationContext := NewGameMutationContext(gameContext.GameState)
	mutationContext.TrackYaku()
	defer mutationContext.BroadcastUpdates()

	mutationContext.Apply(mutation)

	return
}
