package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) ResolveTeyaku(context context.Context, request *lovelove.ResolveTeyakuRequest) (response *lovelove.ResolveTeyakuResponse, rpcError error) {
	response = &lovelove.ResolveTeyakuResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

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

	if gameContext.GameState.state != GameState_Teyaku {
		log.Print("Game in wrong state")
		return
	}

	playerState, playerStateFound := gameContext.GameState.playerState[connMeta.userId]

	if !playerStateFound {
		log.Print("Player not in game")
		return
	}

	teyakuInfo := gameContext.GameState.GetTeyaku()
	playerTeyakuInfo, playerHasTeyaku := teyakuInfo[playerState.position]
	if !playerHasTeyaku {
		log.Print("Player doesn't have teyaku")
		return
	}

	response = &lovelove.ResolveTeyakuResponse{
		Status: lovelove.GenericResponseCode_Ok,
	}

	if playerTeyakuInfo.confirmed {
		log.Print("Player already confirmed")
		return
	}

	playerState.confirmedTeyaku = true
	playerTeyakuInfo.confirmed = true

	opponentTeyakuInfo, opponentHasTeyaku := teyakuInfo[getOpponentPosition(playerState.position)]
	if opponentHasTeyaku && !opponentTeyakuInfo.confirmed {
		return
	}

	winner := playerState.position
	if opponentTeyakuInfo.confirmed {
		winner = lovelove.PlayerPosition_UnknownPosition
	}

	EndRound(gameContext.GameState, &roundEndChange{
		winner,
		teyakuInfo,
	})
	return
}

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

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	defer broadcastBuilder.Broadcast()

	gameMutationContext := NewGameMutationContext(gameContext.GameState)

	yakuTracker := NewYakuTracker(gameContext.GameState)

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))
	broadcastBuilder.QueueUpdates(gameMutationContext.BuildPlayOptions())

	yakuUpdate := yakuTracker.BuildYakuUpdate(gameMutationContext.MovedCards())
	broadcastBuilder.QueueUpdates(yakuUpdate.gameUpdate)

	mutation, err = CheckForYakuMutation(yakuUpdate, playerState, gameContext.GameState)

	if err != nil {
		return
	}

	broadcastBuilder.QueueUpdates(gameMutationContext.Apply(mutation))

	return
}
