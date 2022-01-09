import { stylesheet } from "astroturf";
import clsx from "clsx";
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

		.rematchContainer {
			position: absolute;
			left: 100%;
			top: 0;
			padding: 0px 20px;
			.rematch {
				position: relative;

				word-break: keep-all;
				background-color: white;
				color: black;
				font-weight: 800;

				&.action {
					.rematchBorder {
						position: absolute;
						top: 0;
						right: 0;
						left: 0;
						bottom: 0;
						border: solid 2px white;
						border-radius: 3px;
					}

					background-color: #0005;
					color: white;
					font-weight: 700;
					cursor: pointer;

					&:hover {
						background-color: #0007;
					}
				}
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
	onRematchRequested: () => void;
}) => {
	console.log("redraw curtain");
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
				{opponent.rematchRequested &&
					<div className={styles.rematchContainer}>
						<MetadataBubble className={styles.rematch}>
							もう一回
						</MetadataBubble>
					</div>
				}
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
				<div className={styles.rematchContainer}>
					<MetadataBubble
						className={clsx(styles.rematch, !player.rematchRequested && styles.action)}
						onClick={() => {
							if (!player.rematchRequested) {
								props.onRematchRequested();
							}
						}}
					>
						<div className={styles.rematchBorder}> </div>
						もう一回
					</MetadataBubble>
				</div>
			</div>
		</div>
	</FullScreenCurtain>;
};

