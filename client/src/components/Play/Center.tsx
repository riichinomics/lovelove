import * as React from "react";
import { CardStack } from "./CardStack";
import { ICard } from "../ICard";
import { cardKey } from "./utils";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.center {
		display: flex;

		> * {
			&:not(:last-child) {
				margin-right: 20px;
			}

			> *:not(:last-child) {
				margin-bottom: 20px;
			}
		}

		.deck {
			display: flex;
			flex-direction: column;
		}

		.cards {
			display: flex;
			flex: 1;
			flex-direction: column;

			.cardRow {
				display: flex;
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
			<div>

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
