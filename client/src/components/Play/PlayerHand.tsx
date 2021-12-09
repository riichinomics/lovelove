import * as React from "react";
import { CardStack } from "./CardStack";
import { cardKey, CardZone } from "./utils";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { CardMoveContext } from "../../rpc/CardMoveContext";
import clsx from "clsx";

const styles = stylesheet`
	.playerHand {
		background-color: #222;
		border-top: 2px solid black;

		padding-top: 10px;
		padding-bottom: 10px;

		margin-left: -110px;

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
	cards: lovelove.ICard[];
	canPlay: boolean;
	onPreviewCardChanged?: (card: lovelove.ICard) => void;
}) => {
	const onNoCardSelected = React.useCallback(() => props.onPreviewCardChanged?.(null), [props.onPreviewCardChanged]);
	const { move } = React.useContext(CardMoveContext);
	const moveOrigin = move?.from?.zone === CardZone.Hand ? move.from : null;
	const cards = [null, ...props.cards];
	return <div className={styles.playerHand} onMouseLeave={onNoCardSelected}>
		{cards.map((card, index) => <div className={clsx(props.canPlay && styles.handCard)} key={cardKey(card, index)}>
			<CardStack
				cards={[(moveOrigin?.index === index - 1) ? null : card]}
				onCardSelected={props.onPreviewCardChanged}
				canDrag={props.canPlay}
				zone={CardZone.Hand}
				index={index-1}
			/>
		</div>)}
	</div>;
};
