import * as React from "react";
import { CardStack } from "./CardStack";
import { CardDroppedHandler, cardKey, CardMove, CardWithOffset, CardZone } from "./utils";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { CardMoveContext } from "../../rpc/CardMoveContext";

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
	index: number;
	move: CardMove;
	playOptions: Record<string, lovelove.IPlayOptions>;
	previewCard: lovelove.ICard;
	onCardDropped?: CardDroppedHandler;
}) => {
	const cards = [props.card] as CardWithOffset[];

	if (props.move) {
		cards.push({
			...props.move.from.card,
			offset: props.move.offset
		});
	}

	return <CardStack
		cards={cards}
		playOptions={props.playOptions?.[props.card?.id]?.options ?? []}
		onCardDropped={props.onCardDropped}
		previewCard={props.previewCard}
		zone={CardZone.Table}
		index={props.index}
		stunted
		laminated
	/>;
};

export const Center = (props: {
	deck: number;
	drawnCard?: lovelove.ICard;
	cards: lovelove.ICard[];
	playOptions: Record<string, lovelove.IPlayOptions>;
	previewCard: lovelove.ICard;
	onCardDropped?: CardDroppedHandler;
}) => {
	const { move } = React.useContext(CardMoveContext);
	const moveDestination = move?.to.zone === CardZone.Table ? move.to : null;
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
					.map((card, index) => {
						const tableIndex = (index * 2);
						return <CenterCardStack
							key={cardKey(card, `center_top_${index}`)}
							card={card}
							index={tableIndex}
							move={tableIndex === moveDestination?.index ? move : null}
							{...props}
						/>;
					})}
			</div>
			<div className={styles.cardRow}>
				{props.cards
					.filter((_, index) => index % 2 === 1)
					.map((card, index) => {
						const tableIndex = (index * 2 + 1);
						return <CenterCardStack
							key={cardKey(card, `center_bottom_${index}`)}
							card={card}
							index={tableIndex}
							move={tableIndex === moveDestination?.index ? move : null}
							{...props}
						/>;
					})}
			</div>
		</div>
	</div>;
};
