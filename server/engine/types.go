package engine

import (
	lovelove "hanafuda.moe/lovelove/proto"
)

type GameState int64

const (
	GameState_Teyaku GameState = iota
	GameState_HandCardPlay
	GameState_DeckCardPlay
	GameState_ShoubuOpportunity
)

type YakuCategory int64

const (
	YakuCategory_Other YakuCategory = iota
	YakuCategory_Tane
	YakuCategory_Tanzaku
	YakuCategory_Hikari
)

type cardState struct {
	location CardLocation
	order    int
	card     *lovelove.Card
}

var monthToHanaMap = map[lovelove.Month]lovelove.Hana{
	lovelove.Month_UnknownMonth: lovelove.Hana_UnknownSeason,
	lovelove.Month_January:      lovelove.Hana_Matsu,
	lovelove.Month_February:     lovelove.Hana_Ume,
	lovelove.Month_March:        lovelove.Hana_Sakura,
	lovelove.Month_April:        lovelove.Hana_Fuji,
	lovelove.Month_May:          lovelove.Hana_Ayame,
	lovelove.Month_June:         lovelove.Hana_Botan,
	lovelove.Month_July:         lovelove.Hana_Hagi,
	lovelove.Month_August:       lovelove.Hana_Susuki,
	lovelove.Month_September:    lovelove.Hana_Kiku,
	lovelove.Month_October:      lovelove.Hana_Momiji,
	lovelove.Month_November:     lovelove.Hana_Yanagi,
	lovelove.Month_December:     lovelove.Hana_Kiri,
}

func getHana(month lovelove.Month) lovelove.Hana {
	return monthToHanaMap[month]
}

func getOpponentPosition(playerPosition lovelove.PlayerPosition) lovelove.PlayerPosition {
	if playerPosition == lovelove.PlayerPosition_Red {
		return lovelove.PlayerPosition_White
	}
	return lovelove.PlayerPosition_Red
}
