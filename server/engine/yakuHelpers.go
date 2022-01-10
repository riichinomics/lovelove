package engine

import lovelove "hanafuda.moe/lovelove/proto"

type yakuPart struct {
	id         lovelove.YakuId
	cards      []*cardState
	bonusCards []*cardState
}

func CategoryFor(yakuId lovelove.YakuId) YakuCategory {
	yakuDefinition, ok := yakuDefinitions[yakuId]
	if !ok {
		return YakuCategory_Other
	}

	return yakuDefinition.category
}

type yakuDefinition struct {
	id        lovelove.YakuId
	category  YakuCategory
	hasBonus  bool
	minCards  int
	baseValue int
}

var yakuDefinitions = map[lovelove.YakuId]*yakuDefinition{
	lovelove.YakuId_Tane: {
		id:        lovelove.YakuId_Tane,
		category:  YakuCategory_Tane,
		baseValue: 1,
		minCards:  5,
		hasBonus:  true,
	},

	lovelove.YakuId_Inoshikachou: {
		id:        lovelove.YakuId_Inoshikachou,
		category:  YakuCategory_Tane,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},

	lovelove.YakuId_AkatanAotanNoChoufuku: {
		id:        lovelove.YakuId_AkatanAotanNoChoufuku,
		category:  YakuCategory_Tanzaku,
		baseValue: 10,
		minCards:  6,
		hasBonus:  true,
	},
	lovelove.YakuId_Akatan: {
		id:        lovelove.YakuId_Akatan,
		category:  YakuCategory_Tanzaku,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},
	lovelove.YakuId_Aotan: {
		id:        lovelove.YakuId_Akatan,
		category:  YakuCategory_Tanzaku,
		baseValue: 5,
		minCards:  3,
		hasBonus:  true,
	},
	lovelove.YakuId_Tanzaku: {
		id:        lovelove.YakuId_Tanzaku,
		category:  YakuCategory_Tanzaku,
		baseValue: 1,
		minCards:  5,
		hasBonus:  true,
	},

	lovelove.YakuId_Gokou: {
		id:        lovelove.YakuId_Gokou,
		category:  YakuCategory_Hikari,
		baseValue: 15,
		minCards:  5,
	},
	lovelove.YakuId_Shikou: {
		id:        lovelove.YakuId_Shikou,
		category:  YakuCategory_Hikari,
		baseValue: 8,
		minCards:  4,
	},
	lovelove.YakuId_Ameshikou: {
		id:        lovelove.YakuId_Ameshikou,
		category:  YakuCategory_Hikari,
		baseValue: 7,
		minCards:  4,
	},
	lovelove.YakuId_Sankou: {
		id:        lovelove.YakuId_Sankou,
		category:  YakuCategory_Hikari,
		baseValue: 6,
		minCards:  3,
	},

	lovelove.YakuId_Tsukifuda: {
		id:        lovelove.YakuId_Tsukifuda,
		baseValue: 4,
		minCards:  4,
	},
	lovelove.YakuId_Tsukimizake: {
		id:        lovelove.YakuId_Tsukimizake,
		baseValue: 5,
		minCards:  2,
	},
	lovelove.YakuId_Hanamizake: {
		id:        lovelove.YakuId_Hanamizake,
		baseValue: 5,
		minCards:  2,
	},
	lovelove.YakuId_Kasu: {
		id:        lovelove.YakuId_Kasu,
		baseValue: 1,
		minCards:  10,
		hasBonus:  true,
	},
}

func (yakuPart *yakuPart) IsComplete() bool {
	yakuDefinition, ok := yakuDefinitions[yakuPart.id]
	if !ok {
		return false
	}

	return len(yakuPart.cards) >= yakuDefinition.minCards
}

func (yakuPart *yakuPart) Value() int32 {
	yakuDefinition, ok := yakuDefinitions[yakuPart.id]
	if !ok {
		return 0
	}

	if !yakuDefinition.hasBonus {
		return int32(yakuDefinition.baseValue)
	}

	return int32(yakuDefinition.baseValue + len(yakuPart.cards) + len(yakuPart.bonusCards) - yakuDefinition.minCards)
}

func (yakuPart *yakuPart) AllCardIds() []int32 {
	cards := make([]int32, 0)
	for _, card := range yakuPart.cards {
		cards = append(cards, card.card.Id)
	}

	for _, card := range yakuPart.bonusCards {
		cards = append(cards, card.card.Id)
	}
	return cards
}

func (yakuPart *yakuPart) ToYakuData() *lovelove.YakuData {
	return &lovelove.YakuData{
		Id:    yakuPart.id,
		Cards: yakuPart.AllCardIds(),
		Value: yakuPart.Value(),
	}
}

func cardIsRainman(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Yanagi && card.Variation == lovelove.Variation_Fourth
}

func cardIsSakeCup(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Kiku && card.Variation == lovelove.Variation_Third
}

func cardIsMoon(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Susuki && card.Variation == lovelove.Variation_Fourth
}

func cardIsSakuraCurtain(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Sakura && card.Variation == lovelove.Variation_Fourth
}

func cardIsBoar(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Hagi && card.Variation == lovelove.Variation_Third
}

func cardIsDeer(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Momiji && card.Variation == lovelove.Variation_Third
}

func cardIsButterfly(card *lovelove.Card) bool {
	return card.Hana == lovelove.Hana_Botan && card.Variation == lovelove.Variation_Third
}

func cardIsAotan(card *lovelove.Card) bool {
	return getCardType(card) == CardType_Tanzaku &&
		(card.Hana == lovelove.Hana_Botan ||
			card.Hana == lovelove.Hana_Kiku ||
			card.Hana == lovelove.Hana_Momiji)
}

func cardIsAkatan(card *lovelove.Card) bool {
	return getCardType(card) == CardType_Tanzaku &&
		(card.Hana == lovelove.Hana_Ume ||
			card.Hana == lovelove.Hana_Sakura ||
			card.Hana == lovelove.Hana_Matsu)
}

type yakuContribution struct {
	yakuId      lovelove.YakuId
	isBonusCard bool
}

func YakuContribution(card *lovelove.Card, gameState *gameState) []*yakuContribution {
	yaku := make([]*yakuContribution, 0)
	if card.Hana == monthToHana(gameState.month) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Tsukifuda})
	}

	cardType := getCardType(card)

	switch cardType {
	case CardType_Hikari:
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Gokou},
			&yakuContribution{yakuId: lovelove.YakuId_Ameshikou},
		)

		if !cardIsRainman(card) {
			yaku = append(
				yaku,
				&yakuContribution{yakuId: lovelove.YakuId_Sankou},
				&yakuContribution{yakuId: lovelove.YakuId_Shikou},
			)
		}
	case CardType_Tane:
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Tane},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Inoshikachou,
				isBonusCard: !cardIsDeer(card) && !cardIsBoar(card) && !cardIsButterfly(card),
			},
		)
	case CardType_Tanzaku:
		isAkatan := cardIsAkatan(card)
		isAotan := cardIsAotan(card)
		yaku = append(
			yaku,
			&yakuContribution{yakuId: lovelove.YakuId_Tanzaku},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Akatan,
				isBonusCard: !isAkatan,
			},
			&yakuContribution{
				yakuId:      lovelove.YakuId_Aotan,
				isBonusCard: !isAotan,
			},
			&yakuContribution{
				yakuId:      lovelove.YakuId_AkatanAotanNoChoufuku,
				isBonusCard: !isAotan && !isAkatan,
			},
		)
	case CardType_Kasu:
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Kasu})
	}

	if cardIsSakeCup(card) || cardIsMoon(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Tsukimizake})
	}

	if cardIsSakeCup(card) || cardIsSakuraCurtain(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Hanamizake})
	}

	if cardIsSakeCup(card) {
		yaku = append(yaku, &yakuContribution{yakuId: lovelove.YakuId_Kasu})
	}

	return yaku
}

func GetTeyaku(cards []*cardState) lovelove.TeyakuId {
	if len(cards) < 8 {
		return lovelove.TeyakuId_UnknownTeyaku
	}

	hanaMap := make(map[lovelove.Hana][]*cardState)

	for _, card := range cards {
		hanaCards, hanaExists := hanaMap[card.card.Hana]

		if !hanaExists {
			hanaCards = make([]*cardState, 0)
		}

		hanaMap[card.card.Hana] = append(hanaCards, card)
	}

	completeSets := 0
	pairs := 0

	for _, hanaCards := range hanaMap {
		if len(hanaCards) >= 4 {
			completeSets++
		}

		if len(hanaCards) >= 2 {
			pairs++
		}
	}

	if completeSets >= 2 {
		return lovelove.TeyakuId_Teshi
	}

	if pairs >= 4 {
		return lovelove.TeyakuId_Kuttsuki
	}

	return lovelove.TeyakuId_UnknownTeyaku
}
