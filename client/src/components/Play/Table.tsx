import * as React from "react";
import { Center } from "./Center";
import { Collection } from "./Collection";
import { ICard } from "../ICard";
import { OpponentHand } from "./OpponentHand";
import { PlayerHand } from "./PlayerHand";
import clsx from "clsx";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	$collection-peek: 100px;
	$border-weight: 2px;

	.table {
		display: flex;
		flex-direction: column;
		height: 100vh;

		.collection {
			position: absolute;
			z-index: 99;
			left: 0;
			right: 0;

			padding-top: 10px;
			padding-bottom: 10px;

			overflow: hidden;
			max-height: $collection-peek;
			box-sizing: border-box;
			transition: max-height 350ms ease-in-out;

			display: flex;
			justify-content: center;

			&:hover {
				max-height: 300px;
			}
		}

		.opponentHand {
			position: relative;
			margin-bottom: $collection-peek;

			.opponentCollection {
				top: 100%;
				background-color: #222;
				border-bottom: $border-weight solid black;
				align-items: flex-end;
			}
		}

		.center {
			flex: 1;
			display: flex;
			align-items: stretch;
		}

		.playerArea {
			display: flex;
			flex-direction: column;

			position: relative;
			bottom: 0;
			left: 0;
			right: 0;
			z-index: 99;

			margin-top: $collection-peek;

			display: flex;
			justify-content: center;

			.playerCollection {
				bottom: 100%;
				background-color: #888;
				border-top: $border-weight solid black;
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
		<div className={clsx(styles.collection, styles.opponentCollection)}>
			<Collection cards={props.opponentCollection} stackUpwards />
		</div>
	</div>
	<div className={styles.center}>
		<Center cards={props.sharedCards} deck={props.deck} drawnCard={props.drawnCard} />
	</div>
	<div className={styles.playerArea}>
		<div className={clsx(styles.collection, styles.playerCollection)}>
			<Collection cards={props.playerCollection} />
		</div>
		<PlayerHand cards={props.playerHand} />
	</div>
</div>;
