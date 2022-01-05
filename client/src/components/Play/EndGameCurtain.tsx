import { stylesheet } from "astroturf";
import * as React from "react";
import { lovelove } from "../../rpc/proto/lovelove";
import { FullScreenCurtain as FullScreenCurtain } from "./FullScreenCurtain";
import { MetadataBubble } from "./MetadataBubble";

const styles = stylesheet`
	.status {
		font-size: 48px;
		font-weight: 700;
		margin: 48px 0px;
	}

	.playerArea {
		display: flex;
		justify-content: center;
		.score {
			background-color: #0005;
		}
	}
`;

export const EndGameCurtain = (props: {
	position: lovelove.PlayerPosition;
	gameState: lovelove.ICompleteGameState;
}) => {
	const player = props.position === lovelove.PlayerPosition.Red
		? props.gameState.redPlayer
		: props.gameState.whitePlayer;
	const opponent = props.position === lovelove.PlayerPosition.Red
		? props.gameState.whitePlayer
		: props.gameState.redPlayer;
	return <FullScreenCurtain>
		<div className={styles.playerArea}>
			<MetadataBubble className={styles.score}>
				{opponent.score}
			</MetadataBubble>
		</div>
		<div className={styles.status}>
			{props.gameState.gameEnd.gameWinner === lovelove.PlayerPosition.UnknownPosition
				? "引き分け"
				: props.position === props.gameState.gameEnd.gameWinner
					? "勝利"
					: "負け"}
		</div>
		<div className={styles.playerArea}>
			<MetadataBubble className={styles.score}>
				{player.score}
			</MetadataBubble>
		</div>
	</FullScreenCurtain>;
};
