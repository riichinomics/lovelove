import * as React from "react";
import { CardStack } from "./CardStack";
import { cardKey, CardZone } from "./utils";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { useDispatch } from "react-redux";
import { PreviewCardChangedAction } from "../../state/actions/PreviewCardChangedAction";
import { ActionType } from "../../state/actions/ActionType";

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
	onPreviewCardChanged?: (card: lovelove.ICard) => void;
}) => {
	const onNoCardSelected = React.useCallback(() => props.onPreviewCardChanged?.(null), [props.onPreviewCardChanged]);

	const cards = [null, ...props.cards];
	return <div className={styles.playerHand} onMouseLeave={onNoCardSelected}>
		{cards.map((card, index) => <div className={styles.handCard} key={cardKey(card, index)}>
			<CardStack cards={[card]} onCardSelected={props.onPreviewCardChanged} canDrag zone={CardZone.Hand} />
		</div>)}
	</div>;
};
