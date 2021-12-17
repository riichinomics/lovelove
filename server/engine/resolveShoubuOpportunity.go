package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) ResolveShoubuOpportunity(
	context context.Context,
	request *lovelove.ResolveShoubuOpportunityRequest,
) (response *lovelove.ResolveShoubuOpportunityResponse, rpcError error) {
	response = &lovelove.ResolveShoubuOpportunityResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	log.Print("ResolveShoubuOpportunity")

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

	if gameContext.GameState.activePlayer != playerState.position {
		log.Print("Player is not active")
		return
	}

	if gameContext.GameState.state != GameState_ShoubuOpportunity {
		log.Print("Game in wrong state")
		return
	}

	response = &lovelove.ResolveShoubuOpportunityResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	defer broadcastBuilder.Broadcast()

	gameMutationContext := NewGameMutationContext(gameContext.GameState)

	if request.Shoubu {
		broadcastBuilder.QueueUpdates(EndRound(gameContext.GameState, &roundEndChange{
			winner: playerState.position,
		}))
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply([]*gameStateMutation{
		{
			koikoiChange: map[lovelove.PlayerPosition]*koikoiChange{
				playerState.position: {
					koikoiStatus: true,
				},
			},
		},
	}))

	mutation, err := TurnEndMutation(gameContext.GameState)
	if err != nil {
		log.Print(err.Error())
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))
	return
}
