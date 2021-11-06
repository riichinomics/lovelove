import * as React from "react";
import { CardBackProps } from "../CardBackProps";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.cardBack {
		&:before {
			content: "";
			display: block;
			width: 100%;
			padding-top: 159.933%;
			border: 2px solid black;
			box-sizing: border-box;
		}

		width: 100px;
		background-color: #222;
	}
`;

export const CardBackComponent: React.FC<CardBackProps> = () => <div className={styles.cardBack} />;
