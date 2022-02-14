import { stylesheet } from "astroturf";
import clsx from "clsx";
import * as React from "react";

const styles = stylesheet`
	.waitingCurtain {
		background-color: red;
		overflow: hidden;
		min-height: 800px;
		height: 100vh;
		width: 100%;
		display: flex;
		flex-direction: column;
		border: solid 5px #592116;
		box-sizing: border-box;

		> * {
			transition: min-height 0.3s ease-in-out, border-width 0.3s ease-in-out;
			min-height: 27%;
		}

		&.collapsed {
			> * {
				min-height: 0;
			}

			.core {
				border-top-width: 0px;
				border-bottom-width: 0px;
			}
		}

		.top {
			background-color: #6060bf;
		}

		.core {
			border-top: solid 5px #592116;
			border-bottom: solid 5px #592116;
			flex: 1;
			background-color: #d81e1e;
			display: flex;
			align-items: center;
			justify-content: center;

			.content {
				user-select: none;
				text-align: center;
			}
		}

		.bottom {
			background-color: #50b271;
		}
	}
`;

export const FullScreenCurtain: React.FC<{centerOnly?: boolean}> = ({children, centerOnly}) => {
	return <div className={clsx(styles.waitingCurtain, centerOnly && styles.collapsed)}>
		<div className={styles.top} />
		<div className={styles.core}>
			<div className={styles.content}>
				{children}
			</div>
		</div>
		<div className={styles.bottom} />
	</div>;
};
