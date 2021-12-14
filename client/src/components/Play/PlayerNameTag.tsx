import * as React from "react";
import { stylesheet } from "astroturf";
import clsx from "clsx";

const styles = stylesheet`
	.nameTag {
		box-sizing: border-box;
		position: relative;
		min-height: 0;

		transition: min-height 0.3s ease-in;
		background-color: #d81e1e;

		&.active {
			min-height: 20px;
		}
	}
`;

export const PlayerNameTag: React.FC<{active?: boolean}> = (props) => {
	return <div className={clsx(styles.nameTag, props.active && styles.active)}>
		{props.children}
	</div>;
};