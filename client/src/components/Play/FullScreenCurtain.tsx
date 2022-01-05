import { stylesheet } from "astroturf";
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

		.top {
			min-height: 27%;
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
			min-height: 27%;
			background-color: #50b271;
		}
	}
`;

export const FullScreenCurtain: React.FC = ({children}) => {
	return <div className={styles.waitingCurtain}>
		<div className={styles.top} />
		<div className={styles.core}>
			<div className={styles.content}>
				{children}
			</div>
		</div>
		<div className={styles.bottom} />
	</div>;
};
