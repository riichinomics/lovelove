import * as React from "react";
import { stylesheet } from "astroturf";
import clsx from "clsx";

const styles = stylesheet`
	.nameTag {
		box-sizing: border-box;
		position: relative;

		&.active {
			min-height: 20px;
			background-color: #d81e1e;
		}
	}
`;

export const PlayerNameTag: React.FC<{active?: boolean}> = (props) => {
	return <div className={clsx(styles.nameTag, props.active && styles.active)}>
		{props.children}
	</div>;
};