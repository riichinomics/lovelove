import * as React from "react";
import { stylesheet } from "astroturf";
import { useContext } from "react";
import { lovelove } from "../../rpc/proto/lovelove";
import { ThemeContext } from "../../themes/ThemeContext";
import { getYakuName } from "./utils";

const styles = stylesheet`
	.shoubuOpportunityContainer {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: flex;
		justify-content: center;
		align-items: center;

		.shoubuOpportunity {
			font-size: 24px;
			line-height: 24px;
			font-weight: bold;


			display: flex;
			flex-direction: column;
			align-items: stretch;
			min-width: 400px;
			padding: 40px;
			background-color: #000e;

			.yaku {
				display: flex;
				align-items: center;

				&:not(:first-child) {
					margin-top: 30px;
				}

				.yakuName {
					font-weight: bold;
				}

				.yakuValue {
					margin-left: 10px;
				}

				.yakuCards {
					flex: 1;
					justify-content: flex-end;
					display: flex;
					.yakuCard {
						margin-left: 10px;
						min-width: 50px;
					}
				}
			}

			.actionButtons {
				margin-top: 40px;
				display: flex;
				text-align: center;

				justify-content: center;

				> * {
					padding: 16px;
					min-width: 160px;
					cursor: pointer;

					&:not(:first-child) {
						margin-left: 50px;
					}

					&:hover {
						opacity: 0.9;
					}

					&:active {
						opacity: 0.7;
					}
				}

				.koikoi {
					background-color: #d81e1e;
				}

				.shoubu {
					background-color: gray;
				}
			}
		}

	}
`;

export interface IShoubuOpportunityHandlers {
	koikoiChosen: () => void,
	shoubuChosen: () => void,
}

export const ShoubuOpportunityDisplay = (props: {
	collection: lovelove.ICard[],
	yakuInformation: lovelove.IYakuData[],
} & IShoubuOpportunityHandlers) => {
	const context = useContext(ThemeContext);
	const CardComponent = context.theme.CardComponent;
	return <div className={styles.shoubuOpportunityContainer}>
		<div className={styles.shoubuOpportunity}>
			{props.yakuInformation.map(yaku => <div key={yaku.id} className={styles.yaku}>
				<div className={styles.yakuName}>{getYakuName(yaku.id)}</div>
				<div className={styles.yakuValue}>{yaku.value}</div>
				<div className={styles.yakuCards}>
					{yaku.cards.map(cardId =>
						<CardComponent
							key={cardId}
							{...props.collection.find(collectedCard => collectedCard.id === cardId)}
							className={styles.yakuCard}
						/>
					)}
				</div>
			</div>)}
			<div className={styles.actionButtons}>
				<div className={styles.koikoi} onClick={props.koikoiChosen}>
					こいこい
				</div>
				<div className={styles.shoubu} onClick={props.shoubuChosen}>
					勝負
				</div>
			</div>
		</div>
	</div>;
};