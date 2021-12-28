import { stylesheet } from "astroturf";
import * as React from "react";
import { useHref, useLocation, useResolvedPath } from "react-router";
import { Link } from "react-router-dom";

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
				.title {
					font-size: 48px;
					font-weight: 700;
				}

				.subtitle {
					font-size: 32px;
					margin-bottom: 24px;
				}

				.link {
					cursor: pointer;
					user-select: all;
					font-size: 20px;
					text-decoration: underline;
					font-weight: 700;
					margin-bottom: 24px;
				}

				.help {
					font-size: 20px;
					a {
						color: white;
						&:hover {
							color: #ddd;
						}

						&:active {
							color: #bbb;
						}
					}
				}
			}
		}
		.bottom {
			min-height: 27%;
			background-color: #50b271;
		}
	}
`;

export const WaitingCurtain = (props: {
	roomFull?: boolean,
}) => {
	return <div className={styles.waitingCurtain}>
		<div className={styles.top} />
		<div className={styles.core}>
			<div className={styles.content}>
				<div className={styles.title}>hanafuda.live</div>
				<div className={styles.subtitle}>An online hanafuda game.</div>

				{ props.roomFull
					? <div className={styles.help}>This room is already in use, please <Link to="/">click here</Link> to generate a new room.</div>
					: <>
						<div className={styles.help}>Send this link to a friend to play:</div>
						<div
							className={styles.link}
							onClick={() => navigator.clipboard.writeText(window.location.href)}
						>
							{window.location.href}
						</div>
					</>
				}
				<div className={styles.help}>
					An explanation of the rules can be <a href="https://fudawiki.org/en/hanafuda/games/koi-koi">found here</a>.
				</div>
			</div>
		</div>
		<div className={styles.bottom} />
	</div>;
};
