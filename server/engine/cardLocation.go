package engine

import lovelove "hanafuda.moe/lovelove/proto"

type CardLocation int64

const (
	CardLocation_Deck CardLocation = iota
	CardLocation_Table
	CardLocation_RedHand
	CardLocation_WhiteHand
	CardLocation_RedCollection
	CardLocation_WhiteCollection
	CardLocation_Drawn
)

func (location CardLocation) ToPlayerCentricZone(playerPosition lovelove.PlayerPosition) lovelove.PlayerCentricZone {
	switch location {
	case CardLocation_Deck:
		return lovelove.PlayerCentricZone_Deck
	case CardLocation_Table:
		return lovelove.PlayerCentricZone_Table
	case CardLocation_RedHand:
		if playerPosition == lovelove.PlayerPosition_Red {
			return lovelove.PlayerCentricZone_Hand
		}
		return lovelove.PlayerCentricZone_OpponentHand
	case CardLocation_WhiteHand:
		if playerPosition == lovelove.PlayerPosition_White {
			return lovelove.PlayerCentricZone_Hand
		}
		return lovelove.PlayerCentricZone_OpponentHand
	case CardLocation_RedCollection:
		if playerPosition == lovelove.PlayerPosition_Red {
			return lovelove.PlayerCentricZone_Collection
		}
		return lovelove.PlayerCentricZone_OpponentCollection
	case CardLocation_WhiteCollection:
		if playerPosition == lovelove.PlayerPosition_White {
			return lovelove.PlayerCentricZone_Collection
		}
		return lovelove.PlayerCentricZone_OpponentCollection
	case CardLocation_Drawn:
		return lovelove.PlayerCentricZone_Drawn
	}

	return lovelove.PlayerCentricZone_UnknownZone
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
