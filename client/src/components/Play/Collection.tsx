import * as React from "react";
import { CardType, getCardType, getSeasonMonth } from "./utils";
import { CardStack } from "./CardStack";
import { lovelove } from "../../rpc/proto/lovelove";
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
	cards: lovelove.ICard[];
	stackUpwards?: boolean;
}) => {
	const groups = React.useMemo(
		() => {
			const groups = Object.values(props.cards.reduce(
				(total, next) => {
					const type = getCardType(next);
					total[type] ??= {
						type,
						cards: [],
					};
					total[type].cards.push(next);
					return total;
				},
				{} as Record<number, {
					type: CardType,
					cards: lovelove.ICard[]
				}>
			));

			for (const group of groups) {
				group.cards.sort((a, b) => getSeasonMonth(a.hana) - getSeasonMonth(b.hana));
			}
			return groups;
		},
		[props.cards]
	);

	return <div className={styles.collection}>
		{groups.map(group => <CardStack cards={group.cards} key={group.type} stackDepth={5} stackUpwards={props.stackUpwards} />)}
	</div>;
};
