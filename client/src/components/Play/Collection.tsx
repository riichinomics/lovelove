import * as React from "react";
import { CardType, getCardType, getSeasonMonth, getYakuName } from "./utils";
import { CardStack } from "./CardStack";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import clsx from "clsx";

const styles = stylesheet`
	.collection {
		position: relative;
		display: inline-flex;
		justify-content: center;
		flex-wrap: wrap;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}


		.yakuSelectorContainer {
			position: absolute;
			left: 100%;
			min-width: 200px;
			text-align: left;

			padding-left: 16px;
			padding-top: 4px;

			.yakuSelector {
				cursor: pointer;
				user-select: none;

				display: inline-flex;
				flex-direction: column;
				justify-content: stretch;

				> .yaku {
					margin-bottom: 8px;

					padding: 4px 18px;
					border-radius: 2px;
					background-color: #0005;
					transition: background-color 0.1s ease-in;

					&:hover {
						background-color: #0009;
					}

					&.yakuSelected{
						background-color: #0009;
						&:hover {
							background-color: #0007;
						}
					}

					display: flex;

					line-height: 24px;
					font-size: 18px;

					> .yakuName {
						/* font-weight: bold; */
						flex: 1;
					}

					> .yakuValue {
						margin-left: 32px;
					}
				}
			}
		}
	}
`;

export const Collection = (props: {
	cards: lovelove.ICard[];
	stackUpwards?: boolean;
	yakuInformation?: lovelove.IYakuData[]
}) => {
	const [previewYakuId, setPreviewYakuId] = React.useState<number>(null);
	const [selectedYakuId, setSelectedYakuId] = React.useState<number>(null);
	const yakuId = previewYakuId ?? selectedYakuId;
	const yakuCards = props.yakuInformation?.find(yaku => yaku.id === yakuId)?.cards;
	const groups = React.useMemo(
		() => {
			const groups = Object.values(props.cards.reduce(
				(total, next) => {
					const type = getCardType(next);
					total[type] ??= {
						type,
						cards: [],
					};
					total[type].cards.push(next);
					return total;
				},
				{} as Record<number, {
					type: CardType,
					cards: lovelove.ICard[]
				}>
			));

			for (const group of groups) {
				group.cards.sort((a, b) => getSeasonMonth(a.hana) - getSeasonMonth(b.hana));
			}
			return groups;
		},
		[props.cards]
	);

	return <div className={styles.collection}>
		<div className={styles.yakuSelectorContainer}>
			<div className={styles.yakuSelector}>
				{(props.yakuInformation ?? []).map(yaku =>
					<div
						key={yaku.id}
						className={clsx(styles.yaku, selectedYakuId === yaku.id && styles.yakuSelected)}
						onMouseOver={() => setPreviewYakuId(yaku.id)}
						onMouseOut={() => setPreviewYakuId(null)}
						onClick={() => {
							if (selectedYakuId == yaku.id) {
								setSelectedYakuId(null);
								return;
							}
							setSelectedYakuId(yaku.id);
						}}>
						<div className={styles.yakuName}>{getYakuName(yaku.id)}</div>
						<div className={styles.yakuValue}>{yaku.value}</div>
					</div>
				)}
			</div>
		</div>
		{groups.map(group => <CardStack
			key={group.type}
			cards={group.cards}
			stackDepth={5}
			stackUpwards={props.stackUpwards}
			highlightedCardsIds={yakuCards}
		/>)}
	</div>;
};
