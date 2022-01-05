import { stylesheet } from "astroturf";
import clsx from "clsx";
import * as React from "react";
import { MetadataBubble } from "./MetadataBubble";

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

	.oya {
		background-color: #6060bf;
	}

	.koikoi {
		background-color: #d81e1e;
	}

	.disconnected {
		background-color: #ff8db6;
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
		<MetadataBubble>{props.score ?? 0}</MetadataBubble>
		{props.oya && <MetadataBubble className={styles.oya}>親</MetadataBubble>}
		{props.koikoi && <MetadataBubble className={styles.koikoi}>こいこい</MetadataBubble>}
		{props.disconnected && <MetadataBubble className={styles.disconnected}>接続遮断</MetadataBubble>}
	</div>;
};
