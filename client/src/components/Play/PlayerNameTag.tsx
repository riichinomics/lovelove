import * as React from "react";
import { IPlayerInfo } from "./IPlayerInfo";
import { stylesheet } from "astroturf";
import clsx from "clsx";

const styles = stylesheet`
	.nameTag {
		box-sizing: border-box;
		min-height: 10px;

		background-color: white;

		&.red {
			background-color: red;
		}

		&.active {
			min-height: 20px;
		}
	}
`;

export const PlayerNameTag = (props: IPlayerInfo) => {
	return <div className={clsx(styles.nameTag, props.oya && styles.red, props.active && styles.active)} />;
};