package engine

import (
	lovelove "hanafuda.moe/lovelove/proto"
)

type GameState int64

const (
	GameState_HandCardPlay GameState = iota
	GameState_DeckCardPlay
	GameState_DeclareWin
)

type cardState struct {
	location CardLocation
	order    int
	card     *lovelove.Card
}

type playerState struct {
	id       string
	position lovelove.PlayerPosition
}
