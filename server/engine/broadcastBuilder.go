package engine

import lovelove "hanafuda.moe/lovelove/proto"

type broadcastBuilder struct {
	gameContext         *gameContext
	yakuTracker         *yakuTracker
	gameMutationContext *gameMutationContext
}

func NewBroadcastBuilder(gameContext *gameContext) *broadcastBuilder {
	return &broadcastBuilder{
		gameContext:         gameContext,
		gameMutationContext: NewGameMutationContext(gameContext.GameState),
	}
}

func (builder *broadcastBuilder) TrackYaku() {
	builder.yakuTracker = NewYakuTracker(builder.gameContext.GameState)
}

func (builder *broadcastBuilder) Broadcast() {
	movedCards := builder.gameMutationContext.MovedCards()

	playOptionsUpdateMap := buildPlayOptionsUpdate(builder.gameContext.GameState, movedCards)

	updates := builder.gameMutationContext.UpdatesByPosition()

	yakuUpdates := make(map[lovelove.PlayerPosition]*lovelove.YakuUpdate)
	if builder.yakuTracker != nil {
		yakuUpdates = builder.yakuTracker.GetYakuUpdates(builder.gameMutationContext.MovedCards())
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)

		updatePart := &lovelove.GameStateUpdatePart{
			PlayOptionsUpdate:  playOptionsUpdateMap[position],
			YakuUpdate:         yakuUpdates[position],
			OpponentYakuUpdate: yakuUpdates[getOpponentPosition(position)],
		}

		updates[position] = append(updates[position], updatePart)
	}

	builder.gameContext.BroadcastUpdates(updates)
}
