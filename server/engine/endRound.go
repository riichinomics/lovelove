package engine

import lovelove "hanafuda.moe/lovelove/proto"

func EndRound(
	gameState *gameState,
	roundEndChange *roundEndChange,
) (
	roundEndUpdates GameUpdateMap,
) {
	roundEndUpdates = make(GameUpdateMap)
	if roundEndChange == nil {
		return
	}

	shoubuValue := int32(0)

	if roundEndChange.winner != lovelove.PlayerPosition_UnknownPosition {
		gameState.oya = roundEndChange.winner

		if roundEndChange.teyakuInformation != nil {
			shoubuValue = 6
		} else {
			yakuInfo := gameState.GetYakuData(roundEndChange.winner)
			shoubuValue = gameState.GetShoubuValue(yakuInfo, roundEndChange.winner)
		}
	}

	for _, player := range gameState.playerState {
		player.koikoi = false
		player.confirmedTeyaku = false
		if player.position != roundEndChange.winner {
			continue
		}

		player.score += shoubuValue
	}

	gameState.month = lovelove.Month(int(gameState.month+1) % len(lovelove.Month_name))

	gameState.activePlayer = gameState.oya
	gameState.state = GameState_HandCardPlay

	gameState.Deal()

	if len(gameState.GetTeyaku()) > 0 {
		gameState.state = GameState_Teyaku
	}

	teyakuInformation := make([]*lovelove.RoundEndResultTeyakuInformation, 0)
	if roundEndChange.teyakuInformation != nil {
		for player, playerTeyakuInfo := range roundEndChange.teyakuInformation {
			cards := make([]*lovelove.Card, 0)
			for _, card := range gameState.Hand(player) {
				cards = append(cards, card.card)
			}

			teyakuInformation = append(teyakuInformation, &lovelove.RoundEndResultTeyakuInformation{
				TeyakuId: playerTeyakuInfo.id,
				Cards:    cards,
			})
		}
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		roundEndUpdates[position] = []*lovelove.GameStateUpdatePart{
			{
				RoundEndResult: &lovelove.RoundEndResult{
					Winner:            roundEndChange.winner,
					Winnings:          shoubuValue,
					NextRound:         gameState.ToCompleteGameState(position),
					TeyakuInformation: teyakuInformation,
				},
			},
		}
	}

	return
}
