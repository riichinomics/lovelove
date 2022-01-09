package engine

import (
	"math/rand"

	lovelove "hanafuda.moe/lovelove/proto"
)

func RequestRematch(
	gameState *gameState,
	playerState *playerState,
) (
	rematchRequestedUpdates GameUpdateMap,
) {
	rematchRequestedUpdates = make(GameUpdateMap)
	playerState.requestedRematch = true

	for _, player := range gameState.playerState {
		if player.position == lovelove.PlayerPosition_UnknownPosition {
			continue
		}

		if !player.requestedRematch {
			for p := range lovelove.PlayerPosition_name {
				position := lovelove.PlayerPosition(p)
				rematchRequestedUpdates[position] = []*lovelove.GameStateUpdatePart{
					{
						RematchUpdate: &lovelove.RematchUpdate{
							Player: playerState.position,
						},
					},
				}
			}
			return
		}
	}

	gameState.month = lovelove.Month_January
	for _, player := range gameState.playerState {
		player.score = 0
		player.requestedRematch = false
		player.conceded = false

		player.confirmedTeyaku = false
		player.koikoi = false
	}

	gameState.oya = lovelove.PlayerPosition(rand.Intn(2) + 1)
	gameState.activePlayer = gameState.oya
	gameState.state = GameState_HandCardPlay

	gameState.Deal()

	if len(gameState.GetTeyaku()) > 0 {
		gameState.state = GameState_Teyaku
	}

	for p := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		rematchRequestedUpdates[position] = []*lovelove.GameStateUpdatePart{
			{
				NewGameUpdate: &lovelove.NewGameUpdate{
					GameState: gameState.ToCompleteGameState(position),
				},
			},
		}
	}

	return
}
