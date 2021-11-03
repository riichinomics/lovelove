import * as React from "react";
import { Center } from "./Center";
import { Collection } from "./Collection";
import { ICard } from "../ICard";
import { OpponentHand } from "./OpponentHand";
import { PlayerHand } from "./PlayerHand";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.table {
		flex-direction: column;
		display: flex;

		.center {
			flex: 1;
			margin-top: 20px;
			margin-bottom: 20px;
		}

		.opponentHand {
			margin-bottom: 20px;
		}

		.playerArea {
			display: flex;
			flex-direction: column;

			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			z-index: 99;

			margin-top: 20px;

			display: flex;
			justify-content: center;

			> * {
				&:not(:last-child) {
					margin-right: 10px;
				}
			}

			.handCard {
				transition: margin 200ms ease-in-out;
				&:hover {
					margin-top: -50px;
				}
			}
		}
	}
`;

interface IGameState {
	deck: number;
	sharedCards: ICard[];
	playerHand: ICard[];
	playerCollection: ICard[];
	opponentCards: number;
	opponentCollection: ICard[];
	drawnCard?: ICard;
}

export const Table = (props: IGameState) => <div className={styles.table}>
	<div className={styles.opponentHand}>
		<OpponentHand cards={props.opponentCards} />
	</div>
	<Collection cards={props.opponentCollection} />
	<div className={styles.center}>
		<Center cards={props.sharedCards} deck={props.deck} drawnCard={props.drawnCard} />
	</div>
	<div className={styles.playerArea}>
		<Collection cards={props.playerCollection} />
		<PlayerHand cards={props.playerHand} />
	</div>
</div>;
