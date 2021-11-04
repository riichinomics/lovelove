import { CardBackProps } from "./CardBackProps";
import { CardProps } from "./CardProps";

export interface ITheme {
	cardStackSpacing: number;
	CardComponent: React.FC<CardProps>;
	CardPlaceholderComponent: React.FC;
	CardBackComponent: React.FC<CardBackProps>;
}
