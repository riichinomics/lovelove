import * as React from "react";
import { stylesheet } from "astroturf";
import { useContext } from "react";
import { lovelove } from "../../../rpc/proto/lovelove";
import { ThemeContext } from "../../../themes/ThemeContext";

const styles = stylesheet`
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
			margin-left: 60px;
			flex: 1;
			justify-content: flex-end;
			display: flex;
			.yakuCard {
				margin-left: 10px;
				width: 50px;
			}
		}
	}
`;

export interface IYakuSummary {
	id: string;
	name: string;
	value: number;
	cards: lovelove.ICard[];
}

export const ShoubuYakuDisplay = (props: {
	yakuInformation: IYakuSummary[]
}) => {
	const context = useContext(ThemeContext);
	const CardComponent = context.theme.CardComponent;

	if (!props.yakuInformation) {
		return null;
	}

	return <>
		{props.yakuInformation.map(yaku => <div key={yaku.id} className={styles.yaku}>
			<div className={styles.yakuName}>{yaku.name}</div>
			<div className={styles.yakuValue}>{yaku.value}</div>
			<div className={styles.yakuCards}>
				{yaku.cards.map(card =>
					<div key={card.id} className={styles.yakuCard}>
						<CardComponent hideHints {...card} />
					</div>
				)}
			</div>
		</div>)}
	</>;
};
