package engine

import lovelove "hanafuda.moe/lovelove/proto"

type broadcastBuilder struct {
	gameContext *gameContext
	gameUpdates map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart
}

func NewBroadcastBuilder(gameContext *gameContext) *broadcastBuilder {
	return &broadcastBuilder{
		gameContext: gameContext,
		gameUpdates: make(map[lovelove.PlayerPosition][]*lovelove.GameStateUpdatePart),
	}
}

func (builder *broadcastBuilder) QueueUpdates(updates GameUpdateMap) {
	for position, update := range updates {
		if update == nil {
			continue
		}

		existingUpdates, updatesExist := builder.gameUpdates[position]

		if !updatesExist {
			builder.gameUpdates[position] = update
			continue
		}

		builder.gameUpdates[position] = append(existingUpdates, update...)
	}
}

func (builder *broadcastBuilder) Broadcast() {
	builder.gameContext.BroadcastUpdates(builder.gameUpdates)
}
