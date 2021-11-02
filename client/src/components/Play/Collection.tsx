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

	.opponentHand {
		display: flex;
		justify-content: center;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}
	}
`;

const GROUP_OFFSET = 20;

const Group = (props: {cards: ICard[]}) => {
	const { CardComponent } = React.useContext(ThemeContext).theme;
	const [selectedIndex, setSelectedIndex] = React.useState(props.cards.length - 1);
	const padding = GROUP_OFFSET * (props.cards.length - 1);
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
				left: GROUP_OFFSET * index,
				top: GROUP_OFFSET * index,
				zIndex: selectedIndex - Math.abs(selectedIndex - index)
			}}
			onMouseEnter={() => setSelectedIndex(index)}
		>
			<CardComponent card={card.variation} season={card.season} />
		</div>)}
	</div>;
};

export const Collection = (props: {
	cards: ICard[];
}) => {
	const groups = React.useMemo(
		() => Object.values(props.cards.reduce(
			(total, next) => (total[next.season] ??= [], total[next.season].push(next), total),
			{} as Record<number, ICard[]>
		)).map(group => group.sort((a, b) => a.variation - b.variation)),
		[props.cards]
	);

	return <div className={styles.opponentHand}>
		{groups.map(group => <Group cards={group} key={group[0].season} />)}
	</div>;
};
