import * as React from "react";
import { CardStack } from "./CardStack";
import { ICard } from "../ICard";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.collection {
		display: flex;
		justify-content: center;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}
	}
`;

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

	return <div className={styles.collection}>
		{groups.map(group => <CardStack cards={group} key={group[0].season} />)}
	</div>;
};
