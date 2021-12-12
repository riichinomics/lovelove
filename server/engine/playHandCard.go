package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) PlayHandCard(context context.Context, request *lovelove.PlayHandCardRequest) (response *lovelove.PlayHandCardResponse, rpcError error) {
	response = &lovelove.PlayHandCardResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	log.Print("PlayHandCard")

	if request.HandCard == nil {
		log.Print("No hand card")
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

	mutation, err := PlayHandCardMutation(gameContext.GameState, request, playerState.position)
	if err != nil {
		log.Print(err.Error())
		return
	}

	response = &lovelove.PlayHandCardResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	broadcastBuilder.TrackYaku()
	defer broadcastBuilder.Broadcast()

	broadcastBuilder.gameMutationContext.Apply(mutation)

	mutation, err = DrawCardMutation(gameContext.GameState)

	if err != nil {
		//TODO: report error
		return
	}

	broadcastBuilder.gameMutationContext.Apply(mutation)

	mutation, err = PlayDrawnCardMutation(gameContext.GameState, playerState.position)

	if err != nil {
		return
	}

	broadcastBuilder.gameMutationContext.Apply(mutation)

	if mutation[0].gameStateChange != nil {
		// TODO: check yaku
	}

	return
}
