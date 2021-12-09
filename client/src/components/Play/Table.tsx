import * as React from "react";
import { Center } from "./Center";
import { Collection } from "./Collection";
import { OpponentHand } from "./OpponentHand";
import { PlayerHand } from "./PlayerHand";
import clsx from "clsx";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { PlayerNameTag } from "./PlayerNameTag";
import { CardDroppedHandler, oppositePosition } from "./utils";
import { CardMoveContext } from "../../rpc/CardMoveContext";

const styles = stylesheet`
	$collection-peek: 100px;
	$border-weight: 2px;

	.table {
		overflow: hidden;
		display: inline-flex;
		flex-direction: column;
		min-height: 800px;
		height: 100vh;
		width: 100%;

		.nameTag {
			z-index: 100;
			&.opponentNameTag {
				border-bottom: $border-weight solid black;
			}

			&.playerNameTag {
				border-top: $border-weight solid black;
			}
		}

		.collection {
			position: relative;
			box-sizing: border-box;

			min-height: 40px;

			display: flex;
			flex-direction: column;
			align-items: stretch;

			&.opponentCollection {
				.popup {
					> .collectionWrapper {
						border-bottom: $border-weight solid black;
						background-color: #222;
					}
				}
			}

			&.playerCollection {
				.popup {
					> .collectionWrapper {
						border-top: $border-weight solid black;
						background-color: #888;
					}

					max-height: 280px;
					transition: max-height 350ms ease-in;

					&:hover {
						max-height: 100%;
						transition: max-height 350ms ease-out;
					}
				}
			}

			transition: flex-basis 350ms ease-in-out;


			> .popup {
				z-index: 99;
				max-height: 100%;
				transition: max-height 350ms ease-out;
				display: flex;
				flex-direction: column;
				justify-content: flex-end;
				flex: 1;

				&:hover {
					max-height: 280px;
					transition: max-height 350ms ease-in;
				}

				> .collectionWrapper {
					flex: 1;
					padding-top: 10px;
					padding-bottom: 10px;
					box-sizing: border-box;
				}
			}
		}

		.opponentHand {
			position: relative;
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
			z-index: 200;

			display: flex;
			justify-content: center;
		}
	}
`;

type IGameState = lovelove.ICompleteGameState & {
	position: lovelove.PlayerPosition
}

export const Table = ({
	collection = [],
	deck = 0,
	deckFlipCard,
	hand = [],
	opponentCollection = [],
	opponentHand = 0,
	table = [],
	active,
	oya,
	position,
	tablePlayOptions,
	onCardDropped
}: IGameState & {
	onCardDropped: CardDroppedHandler
}) => {
	const opponentPosition = oppositePosition(position);
	const [previewCard, setPreviewCard] = React.useState<lovelove.ICard>();
	const { move } = React.useContext(CardMoveContext);
	React.useEffect(() => {
		setPreviewCard(null);
	}, [!!move]);

	return <div className={styles.table}>
		<div className={styles.opponentHand}>
			<OpponentHand cards={opponentHand} />
		</div>
		<div className={clsx(styles.nameTag, styles.opponentNameTag)}>
			<PlayerNameTag position={opponentPosition} active={opponentPosition === active} oya={opponentPosition === oya} />
		</div>
		<div className={clsx(styles.collection, styles.opponentCollection)}>
			<div className={styles.popup}>
				<div className={styles.collectionWrapper}>
					<Collection cards={opponentCollection} stackUpwards />
				</div>
			</div>
		</div>
		<div className={styles.center}>
			<Center
				cards={table}
				deck={deck}
				drawnCard={deckFlipCard}
				playOptions={tablePlayOptions}
				previewCard={previewCard}
				onPreviewCardChanged={setPreviewCard}
				onCardDropped={onCardDropped}
			/>
		</div>
		<div className={clsx(styles.collection, styles.playerCollection)}>
			<div className={styles.popup}>
				<div className={styles.collectionWrapper}>
					<Collection cards={collection} />
				</div>
			</div>
		</div>
		<div className={clsx(styles.nameTag, styles.playerNameTag)}>
			<PlayerNameTag position={position} active={position === active} oya={position === oya} />
		</div>
		<div className={styles.playerArea}>
			<PlayerHand
				cards={hand}
				onPreviewCardChanged={setPreviewCard}
				canPlay={tablePlayOptions?.acceptedOriginZones?.indexOf(lovelove.PlayerCentricZone.Hand) >= 0}
			/>
		</div>
	</div>;
};
