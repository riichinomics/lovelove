package engine

import lovelove "hanafuda.moe/lovelove/proto"

func cardIdFromCardDetails(hana int32, variation int32) int32 {
	return (hana-1)*4 + (variation - 1)
}

func cardIdFromCard(card lovelove.Card) int32 {
	return cardIdFromCardDetails(int32(card.Hana), int32(card.Variation))
}

func moveCards(cardMap map[int32]*cardState, cards []*lovelove.Card, location CardLocation) {
	for i, card := range cards {
		cardMap[cardIdFromCard(*card)] = &cardState{
			order:    i,
			card:     card,
			location: location,
		}
	}
}

func createCardMoveUpdate(movingCard *cardState, targetLocation CardLocation, playerPosition lovelove.PlayerPosition, targetIndex int32) *lovelove.CardMoveUpdate {
	return &lovelove.CardMoveUpdate{
		MovedCard: movingCard.card,
		OriginSlot: &lovelove.CardSlot{
			Zone:  movingCard.location.ToPlayerCentricZone(playerPosition),
			Index: int32(movingCard.order),
		},
		DestinationSlot: &lovelove.CardSlot{
			Zone:  targetLocation.ToPlayerCentricZone(playerPosition),
			Index: targetIndex,
		},
	}
}
