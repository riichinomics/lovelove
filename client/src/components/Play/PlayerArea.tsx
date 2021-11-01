import * as React from "react";
import { Card } from "./Card";
import { ICard } from "../ICard";

export const PlayerArea = (props: {
	cards: ICard[];
}) => {
	return <div className={`${""}`}>
		{props.cards.map(card => <Card card={card} key={`${card.season}_${card.variation}`} />)}
	</div>;
};
