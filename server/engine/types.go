package engine

import (
	"strings"

	lovelove "hanafuda.moe/lovelove/proto"
)

type GameState int64

const (
	GameState_Waiting GameState = iota
	GameState_Teyaku
	GameState_HandCardPlay
	GameState_DeckCardPlay
	GameState_ShoubuOpportunity
	GameState_End
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

func monthToHana(month lovelove.Month) lovelove.Hana {
	return monthToHanaMap[month]
}

var monthNameMap = map[string]lovelove.Month{
	"january":   lovelove.Month_January,
	"february":  lovelove.Month_February,
	"march":     lovelove.Month_March,
	"april":     lovelove.Month_April,
	"may":       lovelove.Month_May,
	"june":      lovelove.Month_June,
	"july":      lovelove.Month_July,
	"august":    lovelove.Month_August,
	"september": lovelove.Month_September,
	"october":   lovelove.Month_October,
	"november":  lovelove.Month_November,
	"december":  lovelove.Month_December,
}

func stringToMonth(monthName string) lovelove.Month {
	month, monthFound := monthNameMap[strings.ToLower(monthName)]
	if monthFound {
		return month
	}
	return lovelove.Month_UnknownMonth
}

var hanaNameMap = map[string]lovelove.Hana{
	"matsu":  lovelove.Hana_Matsu,
	"ume,":   lovelove.Hana_Ume,
	"sakura": lovelove.Hana_Sakura,
	"fuji":   lovelove.Hana_Fuji,
	"ayame":  lovelove.Hana_Ayame,
	"botan":  lovelove.Hana_Botan,
	"hagi":   lovelove.Hana_Hagi,
	"susuki": lovelove.Hana_Susuki,
	"kiku":   lovelove.Hana_Kiku,
	"momiji": lovelove.Hana_Momiji,
	"yanagi": lovelove.Hana_Yanagi,
	"kiri":   lovelove.Hana_Kiri,
}

func stringToHana(hanaName string) lovelove.Hana {
	hana, hanaFound := hanaNameMap[strings.ToLower(hanaName)]
	if hanaFound {
		return hana
	}
	return lovelove.Hana_UnknownSeason
}

func getOpponentPosition(playerPosition lovelove.PlayerPosition) lovelove.PlayerPosition {
	if playerPosition == lovelove.PlayerPosition_Red {
		return lovelove.PlayerPosition_White
	}
	return lovelove.PlayerPosition_Red
}
