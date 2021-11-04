import { CardNumber } from "../../themes/CardNumber";
import { ICard } from "../ICard";
import { Season } from "../../themes/Season";
import { Month } from "../../themes/Month";

export enum CardType {
	Kasu,
	Tane,
	Tanzaku,
	Hikari,
}

export function cardKey(card: ICard, extra?: any): string {
	return `${card?.season}_${card?.variation}_${extra}`;
}

export function createRandomCard(): ICard {
	return {
		season: Math.random() * 12 | 0,
		variation: Math.random() * 4 | 0,
	};
}

const CARD_TYPE_MAP: Record<Season, Record<CardNumber, CardType>> = {
	[Season.Ayame]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Botan]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Fuji]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Hagi]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Kiku]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Kiri]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Kasu,
		[CardNumber.Fourth]: CardType.Hikari,
	},
	[Season.Matsu]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tanzaku,
		[CardNumber.Fourth]: CardType.Hikari,
	},
	[Season.Momiji]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Sakura]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tanzaku,
		[CardNumber.Fourth]: CardType.Hikari,
	},
	[Season.Susuki]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Hikari,
	},
	[Season.Ume]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Kasu,
		[CardNumber.Third]: CardType.Tane,
		[CardNumber.Fourth]: CardType.Tanzaku,
	},
	[Season.Yanagi]: {
		[CardNumber.First]: CardType.Kasu,
		[CardNumber.Second]: CardType.Tane,
		[CardNumber.Third]: CardType.Tanzaku,
		[CardNumber.Fourth]: CardType.Hikari,
	}
};

export function getCardType(card: ICard): CardType {
	const type = CARD_TYPE_MAP[card.season as Season][card.variation as CardNumber];
	if (type === undefined) {
		console.log(card);
	}
	return type;
}

const SEASON_MONTH_MAP = {
	[Season.Matsu]: Month.January,
	[Season.Ume]: Month.February,
	[Season.Sakura]: Month.March,
	[Season.Fuji]: Month.April,
	[Season.Ayame]: Month.May,
	[Season.Botan]: Month.June,
	[Season.Hagi]: Month.July,
	[Season.Susuki]: Month.August,
	[Season.Kiku]: Month.September,
	[Season.Momiji]: Month.October,
	[Season.Yanagi]: Month.November,
	[Season.Kiri]: Month.December,
};

export function getSeasonMonth(season: Season): Month {
	return SEASON_MONTH_MAP[season];
}
