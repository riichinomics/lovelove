package engine

import lovelove "hanafuda.moe/lovelove/proto"

// Play options are private per player, the position key indicates which player should receive the update
func buildPlayOptionsUpdate(gameState *gameState, movedCards MovedCardsMeta) (playOptionsUpdateMap map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate) {
	playOptionsUpdateMap = make(map[lovelove.PlayerPosition]*lovelove.PlayOptionsUpdate)

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		playOptionsUpdateMap[position] = &lovelove.PlayOptionsUpdate{
			DefunctOptions: make([]*lovelove.PlayOption, 0),
			NewOptions:     make([]*lovelove.PlayOption, 0),
		}
	}

	tableCards := gameState.Table()

	for _, movingCard := range movedCards {
		if movingCard.cardState.location != CardLocation_Drawn {
			continue
		}

		foundMatch := false

		for _, tableCard := range tableCards {
			if tableCard == nil || !WillAccept(tableCard.card, movingCard.cardState.card) {
				continue
			}

			foundMatch = true

			for _, playOptions := range playOptionsUpdateMap {
				playOptions.NewOptions = append(
					playOptions.NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
						TargetCardId: &lovelove.CardId{CardId: tableCard.card.Id},
					},
				)
			}
		}

		if !foundMatch {
			for _, playOptions := range playOptionsUpdateMap {
				playOptions.NewOptions = append(
					playOptions.NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
					},
				)
			}
		}

		break
	}

	for p, _ := range lovelove.PlayerPosition_name {
		position := lovelove.PlayerPosition(p)
		handCards := gameState.Hand(position)

		for _, movingCard := range movedCards {
			if movingCard.cardState.location != CardLocation_Table {
				continue
			}

			for _, tableCard := range tableCards {
				if tableCard == nil {
					continue
				}

				if movingCard.cardState.card.Id == tableCard.card.Id {
					continue
				}

				if !WillAccept(tableCard.card, movingCard.cardState.card) {
					continue
				}

				playOptionsUpdateMap[position].DefunctOptions = append(
					playOptionsUpdateMap[position].DefunctOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
						TargetCardId: &lovelove.CardId{CardId: tableCard.card.Id},
					},
				)
			}

			for _, handCard := range handCards {
				if !WillAccept(movingCard.cardState.card, handCard.card) {
					continue
				}

				wasUnmatched := true

				for _, tableCard := range tableCards {
					if tableCard == nil {
						continue
					}

					if movingCard.cardState.card.Id == tableCard.card.Id {
						continue
					}

					if WillAccept(tableCard.card, handCard.card) {
						wasUnmatched = false
						break
					}
				}

				if wasUnmatched {
					playOptionsUpdateMap[position].DefunctOptions = append(
						playOptionsUpdateMap[position].DefunctOptions,
						&lovelove.PlayOption{
							OriginCardId: &lovelove.CardId{CardId: handCard.card.Id},
						},
					)
				}

				playOptionsUpdateMap[position].NewOptions = append(
					playOptionsUpdateMap[position].NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{CardId: handCard.card.Id},
						TargetCardId: &lovelove.CardId{CardId: movingCard.cardState.card.Id},
					},
				)
			}
		}

		for _, movingCard := range movedCards {
			if movingCard.originalLocation != CardLocation_Table {
				continue
			}

			playOptionsUpdateMap[position].DefunctOptions = append(
				playOptionsUpdateMap[position].DefunctOptions,
				&lovelove.PlayOption{
					TargetCardId: &lovelove.CardId{
						CardId: movingCard.cardState.card.Id,
					},
				},
			)

			cardsForValidation := make([]*cardState, 0)
			for _, handCard := range handCards {
				if WillAccept(movingCard.cardState.card, handCard.card) {
					cardsForValidation = append(cardsForValidation, handCard)
				}
			}

		CARD_REVALIDATION:
			for _, card := range cardsForValidation {
				for _, tableCard := range tableCards {
					if tableCard != nil && WillAccept(tableCard.card, card.card) {
						continue CARD_REVALIDATION
					}
				}

				playOptionsUpdateMap[position].NewOptions = append(
					playOptionsUpdateMap[position].NewOptions,
					&lovelove.PlayOption{
						OriginCardId: &lovelove.CardId{
							CardId: card.card.Id,
						},
					},
				)
			}
		}
	}

	return
}
