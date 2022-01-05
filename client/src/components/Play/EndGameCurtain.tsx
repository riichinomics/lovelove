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

	.endGameCurtainContent {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.playerArea {
		position: relative;

		.concedeContainer {
			position: absolute;
			right: 100%;
			padding: 0px 20px;
			.concede {
				word-break: keep-all;
				background-color: #0005;
			}
		}

		.score {
			display: inline-block;
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
		<div className={styles.endGameCurtainContent}>
			<div className={styles.playerArea}>
				{opponent.conceded &&
					<div className={styles.concedeContainer}>
						<MetadataBubble className={styles.concede}>
							諦めた
						</MetadataBubble>
					</div>
				}
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
				{player.conceded &&
					<div className={styles.concedeContainer}>
						<MetadataBubble className={styles.concede}>
							諦めた
						</MetadataBubble>
					</div>
				}
				<MetadataBubble className={styles.score}>
					{player.score}
				</MetadataBubble>
			</div>
		</div>
	</FullScreenCurtain>;
};

