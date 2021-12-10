import { Month } from "../../themes/Month";
import { lovelove } from "../../rpc/proto/lovelove";
import { CardProps } from "../../themes/CardProps";

export enum CardType {
	None,
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
	[lovelove.Hana.UnknownSeason]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.None,
		[lovelove.Variation.Second]: CardType.None,
		[lovelove.Variation.Third]: CardType.None,
		[lovelove.Variation.Fourth]: CardType.None,
	},
	[lovelove.Hana.Ayame]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Botan]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Fuji]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Hagi]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Kiku]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Kiri]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Kasu,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Matsu]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tanzaku,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Momiji]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Sakura]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tanzaku,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Susuki]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Hikari,
	},
	[lovelove.Hana.Ume]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
		[lovelove.Variation.First]: CardType.Kasu,
		[lovelove.Variation.Second]: CardType.Kasu,
		[lovelove.Variation.Third]: CardType.Tane,
		[lovelove.Variation.Fourth]: CardType.Tanzaku,
	},
	[lovelove.Hana.Yanagi]: {
		[lovelove.Variation.UnknownVariation]: CardType.None,
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
	[lovelove.Hana.UnknownSeason]: Month.January,
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

export enum CardZone {
	Hand,
	Table,
	Drawn,
}

export interface Vector2 {
	x: number;
	y: number;
}

export interface CardWithOffset extends lovelove.ICard, CardProps {
	offset?: Vector2;
}

export interface CardLocation {
	card: CardWithOffset;
	index?: number;
	zone: CardZone;
}

export type CardDroppedHandler = (
	move: CardMove
) => void;

export interface CardMove {
	from: CardLocation;
	to: CardLocation;
	offset: Vector2;
}

const yamuNameMap: Record<lovelove.YakuId, string> = {
	[lovelove.YakuId.UnknownYaku]: "不明",
	[lovelove.YakuId.Gokou]: "五光",
	[lovelove.YakuId.Shikou]: "四光",
	[lovelove.YakuId.Ameshikou]: "雨四光",
	[lovelove.YakuId.Sankou]: "三光",
	[lovelove.YakuId.Inoshikachou]: "猪鹿蝶",
	[lovelove.YakuId.Tane]: "タネ",
	[lovelove.YakuId.AkatanAotanNoChoufuku]: "赤短・青短の重複",
	[lovelove.YakuId.Akatan]: "赤短",
	[lovelove.YakuId.Aotan]: "青短",
	[lovelove.YakuId.Tanzaku]: "短冊",
	[lovelove.YakuId.Hanamizake]: "花見酒",
	[lovelove.YakuId.Tsukimizake]: "月見酒",
	[lovelove.YakuId.Tsukifuda]: "月札",
	[lovelove.YakuId.Kasu]: "カス",
};

export function getYakuName(yakuId: lovelove.YakuId): string {
	return yamuNameMap[yakuId];
}