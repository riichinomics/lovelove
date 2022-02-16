import * as React from "react";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.links {
		z-index: 999;
		pointer-events: auto;
		position: fixed;
		top: 10px;
		right: 10px;
		display: flex;
		> a {
			font-size: 16px;
			line-height: 16px;
			font-weight: 700;
			cursor: pointer;

			color: black;

			&:hover {
				text-decoration: underline;
				color: #222;
			}

			&:active {
				text-decoration: underline;
				color: #333;
			}
		}
	}
`;

export const InfoLinksWrapper: React.FC = ({children}) => {
	return <>
		<div className={styles.links}>
			<a href="https://github.com/riichinomics/lovelove">More Info</a>
		</div>
		{children}
	</>;
};
