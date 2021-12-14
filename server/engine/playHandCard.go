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
	defer broadcastBuilder.Broadcast()

	gameMutationContext := NewGameMutationContext(gameContext.GameState)

	yakuTracker := NewYakuTracker(gameContext.GameState)

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))

	mutation, err = DrawCardMutation(gameContext.GameState)

	if err != nil {
		//TODO: report error
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))

	mutation, err = PlayDrawnCardMutation(gameContext.GameState, playerState.position)

	if err != nil {
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))
	defer func() {
		broadcastBuilder.QueueUpdates(gameMutationContext.BuildPlayOptions())
	}()

	if mutation[0].gameStateChange != nil {
		return
	}

	yakuUpdate := yakuTracker.BuildYakuUpdate(gameMutationContext.MovedCards())
	broadcastBuilder.QueueUpdates(yakuUpdate.gameUpdate)
	_, hasYakuUpdate := yakuUpdate.yakuUpdatesMap[playerState.position]
	if !hasYakuUpdate {
		mutation, err = TurnEndMutation(gameContext.GameState)
		if err != nil {
			return
		}

		broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))
		return
	}

	mutation, err = ShoubuOpportunityMutation(gameContext.GameState)

	if err != nil {
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))

	return
}
