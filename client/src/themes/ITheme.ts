import { CardBackProps } from "./CardBackProps";
import { CardProps } from "./CardProps";

export interface ITheme {
	cardStackSpacing: number;
	CardComponent: React.FC<CardProps>;
	CardBackComponent: React.FC<CardBackProps>;
}
