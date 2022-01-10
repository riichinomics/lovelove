package engine

import lovelove "hanafuda.moe/lovelove/proto"

type CardLocation int64

const (
	CardLocation_Unknown CardLocation = iota
	CardLocation_Deck
	CardLocation_Table
	CardLocation_RedHand
	CardLocation_WhiteHand
	CardLocation_RedCollection
	CardLocation_WhiteCollection
	CardLocation_Drawn
)

func (location CardLocation) ToCardZone() lovelove.CardZone {
	switch location {
	case CardLocation_Deck:
		return lovelove.CardZone_Deck
	case CardLocation_Table:
		return lovelove.CardZone_Table
	case CardLocation_RedHand:
		return lovelove.CardZone_Hand
	case CardLocation_WhiteHand:
		return lovelove.CardZone_Hand
	case CardLocation_RedCollection:
		return lovelove.CardZone_Collection
	case CardLocation_WhiteCollection:
		return lovelove.CardZone_Collection
	case CardLocation_Drawn:
		return lovelove.CardZone_Drawn
	}

	return lovelove.CardZone_UnknownZone
}

func (location CardLocation) ToPlayerPosition() lovelove.PlayerPosition {
	switch location {
	case CardLocation_RedHand:
		return lovelove.PlayerPosition_Red
	case CardLocation_WhiteHand:
		return lovelove.PlayerPosition_White
	case CardLocation_RedCollection:
		return lovelove.PlayerPosition_Red
	case CardLocation_WhiteCollection:
		return lovelove.PlayerPosition_White
	}

	return lovelove.PlayerPosition_UnknownPosition
}

func GetHandLocation(playerPosition lovelove.PlayerPosition) CardLocation {
	if playerPosition == lovelove.PlayerPosition_White {
		return CardLocation_WhiteHand
	}
	return CardLocation_RedHand
}

func GetCollectionLocation(playerPosition lovelove.PlayerPosition) CardLocation {
	if playerPosition == lovelove.PlayerPosition_White {
		return CardLocation_WhiteCollection
	}
	return CardLocation_RedCollection
}

func LocationIsVisible(cardLocation CardLocation, playerPosition lovelove.PlayerPosition) bool {
	switch cardLocation {
	case CardLocation_Deck:
		return false
	case CardLocation_Drawn:
		return true
	case CardLocation_Table:
		return true
	case CardLocation_WhiteCollection:
		return true
	case CardLocation_RedCollection:
		return true
	case CardLocation_RedHand:
		return playerPosition == lovelove.PlayerPosition_Red
	case CardLocation_WhiteHand:
		return playerPosition == lovelove.PlayerPosition_White
	}
	return false
}
