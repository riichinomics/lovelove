import * as React from "react";
import { CardStack } from "./CardStack";
import { CardDroppedHandler, cardKey, CardZone } from "./utils";
import { lovelove } from "../../rpc/proto/lovelove";
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
			padding-top: 20px;
			padding-bottom: 20px;
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
			padding-top: 20px;
			padding-bottom: 20px;
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

const CenterCardStack = (props: {
	card: lovelove.ICard;
	playOptions: Record<string, lovelove.IPlayOptions>;
	previewCard: lovelove.ICard;
	onCardDropped?: CardDroppedHandler;
}) =>
	<CardStack
		cards={[props.card]}
		playOptions={props.playOptions?.[props.card.id]?.options ?? []}
		onCardDropped={props.onCardDropped}
		previewCard={props.previewCard}
		zone={CardZone.Table}
	/>;

export const Center = (props: {
	deck: number;
	drawnCard?: lovelove.ICard;
	cards: lovelove.ICard[];
	playOptions: Record<string, lovelove.IPlayOptions>;
	previewCard: lovelove.ICard;
	onCardDropped?: CardDroppedHandler;
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
					.filter((_, index) => index % 2 === 0)
					.map((card, index) => <CenterCardStack key={cardKey(card, `center_top_${index}`)} card={card} {...props} />)}
			</div>
			<div className={styles.cardRow}>
				{props.cards
					.filter((_, index) => index % 2 === 1)
					.map((card, index) => <CenterCardStack key={cardKey(card, `center_bottom_${index}`)} card={card} {...props} />)}
			</div>
		</div>
	</div>;
};
