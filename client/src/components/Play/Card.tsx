import * as React from "react";
import { ICard } from "../ICard";
import { ThemeContext } from "../../themes/ThemeContext";

export const Card = ({ card }: { card: ICard; }) => {
	const { CardComponent } = React.useContext(ThemeContext).theme;
	return <CardComponent card={card.variation} season={card.season} />;
};
