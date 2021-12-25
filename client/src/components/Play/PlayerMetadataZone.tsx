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
		min-width: 49px;
		box-sizing: border-box;
		text-align: center;
		font-weight: 800;
		padding: 6px 12px 8px;
		font-size: 24px;

		&.oya {
			background-color: #6060bf;
		}

		&.koikoi {
			background-color: #d81e1e;
		}

		&.disconnected {
			background-color: #ff8db6;
		}
	}
`;


export const PlayerMetadataZone = (props: {
	opponent?: boolean,
	score: number,
	oya?: boolean,
	koikoi?: boolean,
	disconnected?: boolean,
}) => {
	return <div className={clsx(styles.metadataZone, props.opponent && styles.opponent)}>
		<div className={styles.metadataBubble}>{props.score ?? 0}</div>
		{props.oya && <div className={clsx(styles.metadataBubble, styles.oya)}>親</div>}
		{props.koikoi && <div className={clsx(styles.metadataBubble, styles.koikoi)}>こいこい</div>}
		{props.disconnected && <div className={clsx(styles.metadataBubble, styles.disconnected)}>接続遮断</div>}
	</div>;
};
