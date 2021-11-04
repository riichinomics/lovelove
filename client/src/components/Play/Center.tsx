import * as React from "react";
import { CardStack } from "./CardStack";
import { ICard } from "../ICard";
import { cardKey } from "./utils";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.center {
		display: flex;
		flex: 1;

		/* > * {
			&:not(:last-child) {
				margin-right: 20px;
			}

			> *:not(:last-child) {
				margin-bottom: 20px;
			}
		} */

		.deck {
			padding-left: 120px;
			padding-right: 20px;
			background-color: white;
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			border-right: 2px solid black;
			.deckStack {
				margin-bottom: 20px;
			}
		}

		.cards {
			padding-left: 20px;
			display: flex;
			flex: 1;
			flex-direction: column;
			justify-content: center;

			.cardRow {
				display: flex;
				&:not(:last-child) {
					margin-bottom: 20px;
				}

				> * {
					&:not(:last-child) {
						margin-right: 20px;
					}
				}
			}
		}
	}
`;

export const Center = (props: {
	deck: number;
	drawnCard?: ICard;
	cards: ICard[];
}) => {
	return <div className={styles.center}>
		<div className={styles.deck}>
			<div className={styles.deckStack}>
				<CardStack cards={[...new Array(Math.min(props.deck, 3))]} concealed />
			</div>
			<CardStack cards={[props.drawnCard]} />
		</div>
		<div className={styles.cards}>
			<div className={styles.cardRow}>
				{props.cards
					.filter((card, index) => index % 2 === 0)
					.map((card, index) => <CardStack cards={[card]} key={cardKey(card, `center_top_${index}`)} />)}
			</div>
			<div className={styles.cardRow}>
				{props.cards
					.filter((card, index) => index % 2 === 1)
					.map((card, index) => <CardStack cards={[card]} key={cardKey(card, `center_bottom_${index}`)} />)}
			</div>
		</div>
	</div>;
};
