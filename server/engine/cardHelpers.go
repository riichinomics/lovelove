package engine

import lovelove "hanafuda.moe/lovelove/proto"

func cardIdFromCardDetails(hana int32, variation int32) int32 {
	return (hana-1)*4 + (variation - 1) + 1
}

type CardType int64

const (
	CardType_None CardType = iota
	CardType_Kasu
	CardType_Tane
	CardType_Tanzaku
	CardType_Hikari
)

var cardTypeMap = map[lovelove.Hana]map[lovelove.Variation]CardType{
	lovelove.Hana_UnknownSeason: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_None,
		lovelove.Variation_Second:           CardType_None,
		lovelove.Variation_Third:            CardType_None,
		lovelove.Variation_Fourth:           CardType_None,
	},
	lovelove.Hana_Ayame: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Botan: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Fuji: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Hagi: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Kiku: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Kiri: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Kasu,
		lovelove.Variation_Fourth:           CardType_Hikari,
	},
	lovelove.Hana_Matsu: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tanzaku,
		lovelove.Variation_Fourth:           CardType_Hikari,
	},
	lovelove.Hana_Momiji: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Sakura: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tanzaku,
		lovelove.Variation_Fourth:           CardType_Hikari,
	},
	lovelove.Hana_Susuki: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Hikari,
	},
	lovelove.Hana_Ume: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Kasu,
		lovelove.Variation_Third:            CardType_Tane,
		lovelove.Variation_Fourth:           CardType_Tanzaku,
	},
	lovelove.Hana_Yanagi: {
		lovelove.Variation_UnknownVariation: CardType_None,
		lovelove.Variation_First:            CardType_Kasu,
		lovelove.Variation_Second:           CardType_Tane,
		lovelove.Variation_Third:            CardType_Tanzaku,
		lovelove.Variation_Fourth:           CardType_Hikari,
	},
}

func getCardType(card *lovelove.Card) CardType {
	return cardTypeMap[card.Hana][card.Variation]
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

func WillAccept(receiver *lovelove.Card, option *lovelove.Card) bool {
	return receiver.Hana == option.Hana
}
