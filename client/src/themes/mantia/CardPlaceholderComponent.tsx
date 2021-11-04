import * as React from "react";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.cardPlaceholder {
		&:before {
			content: "";
			display: block;
			width: 100%;
			padding-top: 163.933%;
		}

		width: 100px;
	}
`;

export const CardPlaceholderComponent: React.FC = () => <div className={styles.cardPlaceholder} />;
