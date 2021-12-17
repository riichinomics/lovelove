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

	broadcastBuilder := NewBroadcastBuilder(gameContext)
	defer broadcastBuilder.Broadcast()

	broadcastBuilder.QueueUpdates(
		EndRound(gameContext.GameState, &roundEndChange{
			winner,
			teyakuInfo,
		}),
	)
	return
}
