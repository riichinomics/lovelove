import { Month } from "../../themes/Month";
import { lovelove } from "../../rpc/proto/lovelove";

export enum CardType {
	Kasu,
	Tane,
	Tanzaku,
	Hikari,
}

export function cardKey(card: lovelove.ICard, extra?: any): string {
	return `${card?.hana}_${card?.variation}_${extra}`;
}

export function createRandomCard(): lovelove.ICard {
	return {
		hana: Math.random() * 12 | 0,
		variation: Math.random() * 4 | 0,
	};
}

const CARD_TYPE_MAP: Record<lovelove.Hana, Record<lovelove.Variation, CardType>> = {
	[lovelove.Hana.Ayame]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Botan]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Fuji]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Hagi]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Kiku]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Kiri]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Kasu,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Matsu]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tanzaku,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Momiji]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Sakura]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tanzaku,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Susuki]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Ume]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Yanagi]: {
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Tane,
		[lovelove.Variation.Third]: CardType.Tanzaku,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	}
};

export function getCardType(card: lovelove.ICard): CardType {
	const type = CARD_TYPE_MAP[card.hana][card.variation];
	if (type === undefined) {
		console.log(card);
	}
	return type;
}

const SEASON_MONTH_MAP = {
	[lovelove.Hana.Matsu]: Month.January,
	[lovelove.Hana.Ume]: Month.February,
	[lovelove.Hana.Sakura]: Month.March,
	[lovelove.Hana.Fuji]: Month.April,
	[lovelove.Hana.Ayame]: Month.May,
	[lovelove.Hana.Botan]: Month.June,
	[lovelove.Hana.Hagi]: Month.July,
	[lovelove.Hana.Susuki]: Month.August,
	[lovelove.Hana.Kiku]: Month.September,
	[lovelove.Hana.Momiji]: Month.October,
	[lovelove.Hana.Yanagi]: Month.November,
	[lovelove.Hana.Kiri]: Month.December,
};

export function getSeasonMonth(season: lovelove.Hana): Month {
	return SEASON_MONTH_MAP[season];
}

export function oppositePosition(position: lovelove.PlayerPosition) {
	switch (position) {
		case lovelove.PlayerPosition.Red:
			return lovelove.PlayerPosition.White;
		case lovelove.PlayerPosition.White:
			return lovelove.PlayerPosition.Red;
	}
}