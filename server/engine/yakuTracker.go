package engine

import lovelove "hanafuda.moe/lovelove/proto"

type yakuTracker struct {
	gameState       *gameState
	initialYakuData map[lovelove.PlayerPosition][]*lovelove.YakuData
}

func NewYakuTracker(gameState *gameState) (tracker *yakuTracker) {
	tracker = &yakuTracker{
		gameState:       gameState,
		initialYakuData: make(map[lovelove.PlayerPosition][]*lovelove.YakuData),
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		if position == lovelove.PlayerPosition_UnknownPosition {
			continue
		}

		tracker.initialYakuData[position] = tracker.gameState.GetYakuData(position)
	}
	return
}

// Yaku updates are visible to both players, so the position determines who's yaku is being updated
func (tracker *yakuTracker) buildYakuUpdate(movedCards MovedCardsMeta) (
	yakuUpdateMap map[lovelove.PlayerPosition]*lovelove.YakuUpdate,
	yakuInfo []*lovelove.YakuData,
) {
	yakuUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.YakuUpdate)
	if tracker.initialYakuData == nil {
		return
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		if position == lovelove.PlayerPosition_UnknownPosition {
			continue
		}

		collectionLocation := GetCollectionLocation(position)
		collectedCards := make([]*cardState, 0)

		for _, card := range movedCards {
			if card.cardState.location == collectionLocation {
				collectedCards = append(collectedCards, card.cardState)
			}
		}

		if len(collectedCards) == 0 {
			continue
		}

		deletedYaku := make(map[lovelove.YakuId]bool)
		newOrUpdatedYaku := make(map[lovelove.YakuId]*lovelove.YakuUpdatePart)

		yakuInfo = tracker.gameState.GetYakuData(position)
		for _, card := range collectedCards {
			yakuContribution := YakuContribution(card.card, tracker.gameState)
			for _, possibleYaku := range yakuContribution {
				for _, yaku := range yakuInfo {
					if possibleYaku.yakuId != yaku.Id {
						continue
					}

					yakuCategory := CategoryFor(yaku.Id)
					sameCategoryExisted := false
					sameCategoryYakuId := lovelove.YakuId_UnknownYaku
					yakuExisted := false
					for _, existingYaku := range tracker.initialYakuData[position] {
						if yakuCategory != YakuCategory_Other && yakuCategory == CategoryFor(existingYaku.Id) {
							sameCategoryExisted = true
							sameCategoryYakuId = existingYaku.Id

							if existingYaku.Id == yaku.Id {
								yakuExisted = true
							}

							break
						}

						if existingYaku.Id == yaku.Id {
							yakuExisted = true
							sameCategoryExisted = true
							sameCategoryYakuId = existingYaku.Id
							break
						}
					}

					if !yakuExisted && sameCategoryExisted {
						deletedYaku[sameCategoryYakuId] = true
					}

					existingUpdate, updateExisted := newOrUpdatedYaku[yaku.Id]
					if !updateExisted {
						existingUpdate = &lovelove.YakuUpdatePart{
							YakuId: yaku.Id,
							Value:  yaku.Value,
						}
						newOrUpdatedYaku[yaku.Id] = existingUpdate
					}

					if yakuExisted {
						if existingUpdate.CardIds == nil {
							existingUpdate.CardIds = make([]int32, 0)
						}
						existingUpdate.CardIds = append(existingUpdate.CardIds, card.card.Id)
						continue
					}

					existingUpdate.CardIds = yaku.Cards
				}
			}
		}

		if len(deletedYaku) == 0 && len(newOrUpdatedYaku) == 0 {
			continue
		}

		update := &lovelove.YakuUpdate{
			DeletedYaku:      make([]lovelove.YakuId, 0),
			NewOrUpdatedYaku: make([]*lovelove.YakuUpdatePart, 0),
		}
		yakuUpdateMap[position] = update

		for yaku := range deletedYaku {
			update.DeletedYaku = append(update.DeletedYaku, yaku)
		}

		for _, yaku := range newOrUpdatedYaku {
			update.NewOrUpdatedYaku = append(update.NewOrUpdatedYaku, yaku)
		}
	}

	return
}

type yakuUpdate struct {
	gameUpdate       GameUpdateMap
	completeYakuInfo []*lovelove.YakuData
	yakuUpdatesMap   map[lovelove.PlayerPosition]*lovelove.YakuUpdate
}

func (tracker *yakuTracker) BuildYakuUpdate(movedCards MovedCardsMeta) *yakuUpdate {
	yakuUpdatesMap, completeYakuUpdate := tracker.buildYakuUpdate(movedCards)

	gameUpdate := make(GameUpdateMap)
	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)

		yakuUpdate := yakuUpdatesMap[position]
		opponentYakuUpdate := yakuUpdatesMap[getOpponentPosition(position)]
		if yakuUpdate == nil && opponentYakuUpdate == nil {
			continue
		}

		gameUpdate[position] = []*lovelove.GameStateUpdatePart{
			{
				YakuUpdate:         yakuUpdate,
				OpponentYakuUpdate: opponentYakuUpdate,
			},
		}
	}

	return &yakuUpdate{
		gameUpdate,
		completeYakuUpdate,
		yakuUpdatesMap,
	}
}
