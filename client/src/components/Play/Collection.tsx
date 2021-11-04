import * as React from "react";
import { CardStack } from "./CardStack";
import { ICard } from "../ICard";
import { getCardType } from "./utils";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.collection {
		display: flex;
		justify-content: center;
		flex-wrap: wrap;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}
	}
`;

export const Collection = (props: {
	cards: ICard[];
	stackUpwards?: boolean;
}) => {
	const groups = React.useMemo(
		() => Object.values(props.cards.reduce(
			(total, next) => {
				const type = getCardType(next);
				total[type] ??= [];
				total[type].push(next);
				return total;
			},
			{} as Record<number, ICard[]>
		)).map(group => group.sort((a, b) => a.variation - b.variation)),
		[props.cards]
	);

	return <div className={styles.collection}>
		{groups.map(group => <CardStack cards={group} key={group[0].season} stackDepth={5} stackUpwards={props.stackUpwards} />)}
	</div>;
};