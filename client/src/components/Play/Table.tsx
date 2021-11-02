import * as React from "react";
import { Center } from "./Center";
import { Collection } from "./Collection";
import { ICard } from "../ICard";
import { OpponentHand } from "./OpponentHand";
import { PlayerArea } from "./PlayerArea";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.table {
		display: flex;
		flex-direction: column;

		> * {
			&:not(:last-child) {
				margin-bottom: 20px;
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
}

export const Table = (props: IGameState) => <div className={styles.table}>
	<OpponentHand cards={props.opponentCards} />
	<Collection cards={props.opponentCollection} />
	<PlayerArea cards={[
		{
			season: 2,
			variation: 2
		},
		{
			season: 5,
			variation: 1
		},
		{
			season: 8,
			variation: 3
		},
		{
			season: 0,
			variation: 0
		},
		{
			season: 11,
			variation: 2
		},
	]}
	/>
	<Center cards={[]} />
	<PlayerArea cards={[]} />
</div>;
