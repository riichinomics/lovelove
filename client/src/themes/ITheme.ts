import { CardBackProps } from "./CardBackProps";
import { CardProps } from "./CardProps";

export interface ITheme {
	CardComponent: React.FC<CardProps>;
	CardBackComponent: React.FC<CardBackProps>;
}
