import { CardBackProps } from "./CardBackProps";
import { CardProps } from "./CardProps";

export interface CardStackSpacing {
	vertical: number;
	horizontal: number;
}

export interface ITheme {
	cardStackSpacing: CardStackSpacing;
	CardComponent: React.FC<CardProps>;
	CardPlaceholderComponent: React.FC;
	CardBackComponent: React.FC<CardBackProps>;
}
