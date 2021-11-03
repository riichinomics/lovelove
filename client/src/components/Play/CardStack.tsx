import * as React from "react";
import { ICard } from "../ICard";
import { ThemeContext } from "../../themes/ThemeContext";
import clsx from "clsx";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.collectionGroup {
		position: relative;
	}

	.collectionGroupItem {
		position: absolute;
	}
`;

export const CardStack = (props: {cards: ICard[]}) => {
	const { CardComponent, cardStackSpacing } = React.useContext(ThemeContext).theme;
	const [selectedIndex, setSelectedIndex] = React.useState(props.cards.length - 1);
	const padding = cardStackSpacing * (props.cards.length - 1);
	return <div
		className={styles.collectionGroup}
		style={{
			paddingRight: padding,
			paddingBottom: padding,
			zIndex: props.cards.length
		}}
	>
		{props.cards.map((card, index) => <div
			className={clsx(index && styles.collectionGroupItem)}
			key={`${card.season}_${card.variation}_${index}`}
			style={{
				left: cardStackSpacing * index,
				top: cardStackSpacing * index,
				zIndex: selectedIndex - Math.abs(selectedIndex - index)
			}}
			onMouseEnter={() => setSelectedIndex(index)}
		>
			<CardComponent card={card.variation} season={card.season} />
		</div>)}
	</div>;
};
