import { ICard } from "../ICard";

export function cardKey(card: ICard, extra?: any): string {
	return `${card.season}_${card.variation}_${extra}`;
}

export function createRandomCard(): ICard {
	return {
		season: Math.random() * 12 | 0,
		variation: Math.random() * 4 | 0,
	};
}
