package engine

import lovelove "hanafuda.moe/lovelove/proto"

func ConcedeGameChange(gameState *gameState, player *playerState) (concessionUpdates GameUpdateMap) {
	concessionUpdates = make(GameUpdateMap)

	gameState.state = GameState_End
	player.conceded = true

	for p := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		concessionUpdates[position] = []*lovelove.GameStateUpdatePart{
			{
				RoundEndResult: &lovelove.RoundEndResult{
					NextRound: gameState.ToCompleteGameState(position),
				},
			},
		}
	}
	return
}
