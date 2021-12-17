
import * as React from "react";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.shoubuValueRow {
		margin-top: 40px;
		display: flex;
		.shoubuValue {
			text-align: right;
			flex: 1;
		}
	}
`;

export const ShoubuTotal = (props: {
	total?: number
}) => {
	if (!props.total) {
		return null;
	}
	return <div className={styles.shoubuValueRow}>
		<div>合計</div>
		<div className={styles.shoubuValue}>{props.total}</div>
	</div>;
};
