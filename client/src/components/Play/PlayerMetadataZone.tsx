import { stylesheet } from "astroturf";
import clsx from "clsx";
import * as React from "react";

const styles = stylesheet`
	.metadataZone {
		position: absolute;
		display: flex;
		justify-content: flex-start;

		bottom: 100%;
		left: 0;
		padding: 10px;
		&.opponent {
			bottom: auto;
			top: 100%;
		}

		> *:not(:first-child) {
			margin-left: 10px;
		}
	}


	.metadataBubble {
		border-radius: 3px;
		background-color: #333;
		min-width: 80px;
		text-align: center;
		font-weight: 800;
		padding: 6px 12px 8px;
		font-size: 24px;

		&.oya {
			min-width: auto;
			background-color: #6060bf;
		}

		&.koikoi {
			min-width: auto;
			background-color: #d81e1e;
		}
	}
`;


export const PlayerMetadataZone = (props: {
	opponent?: boolean,
	score: number,
	oya?: boolean,
	koikoi?: boolean,
}) => {
	return <div className={clsx(styles.metadataZone, props.opponent && styles.opponent)}>
		<div className={styles.metadataBubble}>{props.score ?? 0}</div>
		{props.oya && <div className={clsx(styles.metadataBubble, styles.oya)}>親</div>}
		{props.koikoi && <div className={clsx(styles.metadataBubble, styles.koikoi)}>こいこい</div>}
	</div>;
};
