import { stylesheet } from "astroturf";
import * as React from "react";
import { Link } from "react-router-dom";
import { FullScreenCurtain } from "./FullScreenCurtain";

const styles = stylesheet`
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
`;

export const WaitingCurtain = (props: {
	roomFull?: boolean,
	connected?: boolean,
}) => {
	return <FullScreenCurtain centerOnly={!props.connected}>
		<div className={styles.title}>hanafuda.live</div>
		<div className={styles.subtitle}>An online hanafuda game.</div>

		{ props.connected
			? props.roomFull
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
			: <div className={styles.subtitle}>Connecting...</div>
		}
		<div className={styles.help}>
			An explanation of the rules can be <a href="https://fudawiki.org/en/hanafuda/games/koi-koi">found here</a>.
		</div>
	</FullScreenCurtain>;
};
