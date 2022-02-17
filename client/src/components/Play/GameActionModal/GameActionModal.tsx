import { stylesheet } from "astroturf";
import * as React from "react";
import { lovelove } from "../../../rpc/proto/lovelove";
import { RoundEndInformation } from "../../../state/IState";
import { getTeyakuName, getYakuName } from "../../../utils";
import { CenterModal, CenterModalActionPanel } from "./../CenterModal";
import { ShoubuTotal } from "./ShoubuValueRow";
import { IYakuSummary, ShoubuYakuDisplay } from "./ShoubuYakuDisplay";

const styles = stylesheet`
	.caption {
		text-align: center;
		font-size: 30px;
	}

	.koikoi {
		background-color: #d81e1e !important;
	}
`;

export interface IGameModalActions {
	onKoikoiChosen?: () => void,
	onShoubuChosen?: (teyaku: boolean) => void,
	onContinueChosen?: () => void,
	teyakuResolved: boolean,
}

export const GameActionModal = (props: {
	teyaku: lovelove.TeyakuId,
	shoubuOpportunity: lovelove.IShoubuOpportunity,
	roundEndView: RoundEndInformation,
	position: lovelove.PlayerPosition,
	yakuInformation: lovelove.IYakuData[],
	opponentYakuInformation: lovelove.IYakuData[],
	collection: lovelove.ICard[],
	opponentCollection: lovelove.ICard[],
	hand: lovelove.ICard[],
} & IGameModalActions) => {
	if (props.teyakuResolved || (!props.shoubuOpportunity && !props.roundEndView && !props.teyaku)) {
		return null;
	}

	const yakuTarget = props.roundEndView
		? props.roundEndView.winner
		: props.shoubuOpportunity
			? props.position
			: null;
	const collection = yakuTarget == props.position ? props.collection : props.opponentCollection;
	const yakuInformation = yakuTarget == props.position ? props.yakuInformation : props.opponentYakuInformation;

	const yakuSummaries = React.useMemo<IYakuSummary[]>(() => {
		if (!yakuTarget || props.roundEndView?.teyaku) {
			return null;
		}

		return yakuInformation.map(yaku => ({
			id: `yaku_${yaku.id}`,
			cards: yaku.cards.map(cardId => collection.find(card => card.id == cardId)),
			name: getYakuName(yaku.id),
			value: yaku.value
		}));
	}, [yakuTarget, collection, yakuInformation]);

	const teyakuSummaries = React.useMemo<IYakuSummary[]>(() => {
		if (props.roundEndView) {
			if (props.roundEndView.teyaku) {
				return props.roundEndView.teyaku.map(teyaku => ({
					id: `teyaku_${teyaku.teyakuId}`,
					name: getTeyakuName(teyaku.teyakuId),
					cards: teyaku.cards,
					value: 6,
				}));
			}

			return null;
		}

		if (props.teyaku) {
			return [{
				cards: props.hand,
				id: `teyaku_${props.teyaku}`,
				name: getTeyakuName(props.teyaku),
				value: 6,
			}];
		}

		return null;
	}, [props.roundEndView, props.teyaku]);

	const summaries = yakuSummaries ?? teyakuSummaries;

	const shoubuTotal = props.roundEndView
		? props.roundEndView.winnings
		: props.shoubuOpportunity?.value;

	return <CenterModal>
		{ (props.roundEndView && props.roundEndView.winner == lovelove.PlayerPosition.UnknownPosition) && <div className={styles.caption}>
			{props.roundEndView.teyaku?.length > 0 ? "引き分け" : "流局"}
		</div>}
		<ShoubuYakuDisplay yakuInformation={summaries} />
		<ShoubuTotal total={shoubuTotal} />
		<CenterModalActionPanel>
			{(!props.roundEndView && props.shoubuOpportunity) && <div className={styles.koikoi} onClick={props.onKoikoiChosen}>
				こいこい
			</div>}

			{(!props.roundEndView && (props.shoubuOpportunity || props.teyaku != lovelove.TeyakuId.UnknownTeyaku))
				&& <div onClick={() => props.onShoubuChosen(props.teyaku != null && props.teyaku != lovelove.TeyakuId.UnknownTeyaku)}>
					勝負
				</div>}

			{props.roundEndView && <div onClick={props.onContinueChosen}>
				確認
			</div>}
		</CenterModalActionPanel>
	</CenterModal>;
};
