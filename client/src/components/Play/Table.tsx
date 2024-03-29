import * as React from "react";
import { Center } from "./Center";
import { Collection } from "./Collection";
import { OpponentHand } from "./OpponentHand";
import { PlayerHand } from "./PlayerHand";
import clsx from "clsx";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { PlayerNameTag } from "./PlayerNameTag";
import { CardDroppedHandler, oppositePosition } from "../../utils";
import { CardMoveContext } from "../../rpc/CardMoveContext";
import { PlayerMetadataZone } from "./PlayerMetadataZone";
import { RoundEndInformation } from "../../state/IState";
import { GameActionModal, IGameModalActions } from "./GameActionModal/GameActionModal";
import { MetadataBubble } from "./MetadataBubble";
import { CenterModal, CenterModalActionPanel } from "./CenterModal";

const styles = stylesheet`
	$collection-peek: 100px;
	$border-weight: 2px;

	.table {
		display: inline-flex;
		flex-direction: column;
		height: 100vh;
		width: 100%;

		.nameTag {
			z-index: 100;
			flex-basis: auto;
			&.opponentNameTag {
				border-bottom: $border-weight solid black;

				.concedeContainer {
					position: absolute;
					padding: 10px;
					top: 100%;
					right: 0;
				}

				.concede {
					cursor: pointer;
					transition: background-color 0.3s ease-in-out;
					&:hover {
						background-color: #666;
					}
				}
			}

			&.playerNameTag {
				border-top: $border-weight solid black;
			}
		}

		.modalArea {
			pointer-events: none;
			position: fixed;
			top: 0;
			right: 0;
			bottom: 0;
			left: 0;
			z-index: 999;

			.cancelConcede {
				background-color: #d81e1e;
			}
		}

		.collection {
			position: relative;
			box-sizing: border-box;

			min-height: 68px;

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
					text-align: center;
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

export const Table = (props: IGameState & IGameModalActions & {
	opponentDisconnected: boolean,
	gameState: lovelove.ICompleteGameState,
	onCardDropped: CardDroppedHandler,
	roundEndView?: RoundEndInformation,
	onGameConceded: () => void,
}) => {
	const {
		opponentDisconnected,
		position,
		onCardDropped,
		onKoikoiChosen,
		onShoubuChosen,
		onContinueChosen,
		teyakuResolved,
		roundEndView,
		onGameConceded,
		gameState: {
			redPlayer,
			whitePlayer,
			deck = 0,
			deckFlipCard,
			table = [],
			active,
			oya,
			tablePlayOptions,
			shoubuOpportunity,
			month,
			monthHana,
			teyaku,
		} = {}
	} = props;

	const player = position == lovelove.PlayerPosition.Red ? redPlayer : whitePlayer;
	const opponent = position == lovelove.PlayerPosition.Red ? whitePlayer : redPlayer;

	const opponentPosition = oppositePosition(position);
	const [previewCard, setPreviewCard] = React.useState<lovelove.ICard>();
	const [concessionRequested, setConcessionRequested] = React.useState(false);
	const { move } = React.useContext(CardMoveContext);
	React.useEffect(() => {
		setPreviewCard(null);
	}, [!!move]);

	return <div className={styles.table}>
		<div className={styles.modalArea}>
			{ concessionRequested &&
				<CenterModal>
					<div>
						Are you sure you want to concede?
					</div>
					<CenterModalActionPanel>
						<div className={styles.cancelConcede} onClick={() => setConcessionRequested(false)}>いいえ</div>
						<div onClick={onGameConceded}>はい</div>
					</CenterModalActionPanel>
				</CenterModal>
			}

			<GameActionModal
				onKoikoiChosen={onKoikoiChosen}
				onShoubuChosen={onShoubuChosen}
				onContinueChosen={onContinueChosen}
				roundEndView={roundEndView}
				shoubuOpportunity={shoubuOpportunity}
				teyaku={teyaku}
				collection={player?.collection}
				opponentCollection={opponent?.collection}
				yakuInformation={player?.yakuInformation}
				opponentYakuInformation={opponent?.yakuInformation}
				hand={player?.hand?.cards}
				position={position}
				teyakuResolved={teyakuResolved}
			/>
		</div>
		<div className={styles.opponentHand}>
			<OpponentHand cards={opponent?.hand?.numberOfCards} />
		</div>
		<div className={clsx(styles.nameTag, styles.opponentNameTag)}>
			<PlayerNameTag active={opponentPosition === active}>
				<PlayerMetadataZone
					opponent
					oya={opponentPosition === oya}
					score={opponent?.score}
					koikoi={opponent?.koikoi}
					disconnected={opponentDisconnected}
				/>
				<div className={styles.concedeContainer}>
					<MetadataBubble
						onClick={() => setConcessionRequested(true)}
						className={styles.concede}
					>
						諦める
					</MetadataBubble>
				</div>
			</PlayerNameTag>
		</div>
		<div className={clsx(styles.collection, styles.opponentCollection)}>
			<div className={styles.popup}>
				<div className={styles.collectionWrapper}>
					<Collection
						cards={opponent?.collection}
						yakuInformation={opponent?.yakuInformation}
						stackUpwards
					/>
				</div>
			</div>
		</div>
		<div className={styles.center}>
			<Center
				month={month}
				monthHana={monthHana}
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
					<Collection cards={player?.collection} yakuInformation={player?.yakuInformation} />
				</div>
			</div>
		</div>
		<div className={clsx(styles.nameTag, styles.playerNameTag)}>
			<PlayerNameTag active={position === active}>
				<PlayerMetadataZone
					oya={position === oya}
					score={player?.score}
					koikoi={player?.koikoi}
				/>
			</PlayerNameTag>
		</div>
		<div className={styles.playerArea}>
			<PlayerHand
				cards={player?.hand?.cards}
				onPreviewCardChanged={setPreviewCard}
				canPlay={tablePlayOptions?.acceptedOriginZones?.indexOf(lovelove.CardZone.Hand) >= 0}
			/>
		</div>
	</div>;
};
