import { stylesheet } from "astroturf";
import * as React from "react";

const styles = stylesheet`
	.centerModalContainer {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: flex;
		justify-content: center;
		align-items: center;

		pointer-events: none;

		.centerModal {
			pointer-events: auto;
			font-size: 24px;
			font-weight: 800;


			display: flex;
			flex-direction: column;
			align-items: stretch;
			min-width: 400px;
			padding: 40px;
			background-color: #000e;
		}

		.actionButtons {
			font-size: 32px;
			margin-top: 40px;
			display: flex;
			text-align: center;

			justify-content: center;

			* {
				padding: 2px 8px 8px;
				min-width: 160px;
				cursor: pointer;

				background-color: gray;

				&:not(:first-child) {
					margin-left: 50px;
				}

				&:hover {
					opacity: 0.9;
				}

				&:active {
					opacity: 0.7;
				}
			}
		}
	}
`;

export const CenterModal: React.FC = ({children}) => {
	return <div className={styles.centerModalContainer}>
		<div className={styles.centerModal}>{children}</div>
	</div>;
};


export const CenterModalActionPanel: React.FC = ({children}) => {
	return <div className={styles.actionButtons}>
		{children}
	</div>;
};
