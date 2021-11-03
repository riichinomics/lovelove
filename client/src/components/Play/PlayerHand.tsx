import * as React from "react";
import { CardStack } from "./CardStack";
import { ICard } from "../ICard";
import { cardKey } from "./utils";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.playerHand {
		background-color: #222;

		padding-top: 10px;
		padding-bottom: 10px;

		display: flex;
		justify-content: center;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}

		.handCard {
			transition: margin 150ms ease-in-out;
			&:hover {
				margin-top: -20px;
			}
		}
	}
`;

export const PlayerHand = (props: {
	cards: ICard[];
}) => {
	return <div className={styles.playerHand}>
		{props.cards.map((card, index) => <div className={styles.handCard} key={cardKey(card, index)}>
			<CardStack cards={[card]} />
		</div>)}
	</div>;
};
