package engine

import (
	"google.golang.org/protobuf/proto"
	lovelove "hanafuda.moe/lovelove/proto"
)

type GameState int64

const (
	GameState_HandCardPlay GameState = iota
	GameState_DeckCardPlay
	GameState_DeclareWin
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

type playerState struct {
	id        string
	position  lovelove.PlayerPosition
	listeners []chan proto.Message
}

var hanaToMonthMap = map[lovelove.Hana]lovelove.Month{
	lovelove.Hana_UnknownSeason: lovelove.Month_UnknownMonth,
	lovelove.Hana_Matsu:         lovelove.Month_January,
	lovelove.Hana_Ume:           lovelove.Month_February,
	lovelove.Hana_Sakura:        lovelove.Month_March,
	lovelove.Hana_Fuji:          lovelove.Month_April,
	lovelove.Hana_Ayame:         lovelove.Month_May,
	lovelove.Hana_Botan:         lovelove.Month_June,
	lovelove.Hana_Hagi:          lovelove.Month_July,
	lovelove.Hana_Susuki:        lovelove.Month_August,
	lovelove.Hana_Kiku:          lovelove.Month_September,
	lovelove.Hana_Momiji:        lovelove.Month_October,
	lovelove.Hana_Yanagi:        lovelove.Month_November,
	lovelove.Hana_Kiri:          lovelove.Month_December,
}

func getMonth(hana lovelove.Hana) lovelove.Month {
	return hanaToMonthMap[hana]
}

func getOpponentPosition(playerPosition lovelove.PlayerPosition) lovelove.PlayerPosition {
	if playerPosition == lovelove.PlayerPosition_Red {
		return lovelove.PlayerPosition_White
	}
	return lovelove.PlayerPosition_Red
}
