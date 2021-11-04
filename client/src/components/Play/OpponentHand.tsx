import * as React from "react";
import { ThemeContext } from "../../themes/ThemeContext";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.opponentHandWrapper {
		position: relative;
		height: 50px;
		overflow: hidden;
	}

	.opponentHand {
		min-height: 50px;
		background-color: #eee;
		padding-bottom: 10px;
		border-bottom: 2px solid black;
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;

		display: flex;
		justify-content: center;

		> * {
			&:not(:last-child) {
				margin-right: 10px;
			}
		}
	}
`;

export const OpponentHand = (props: {
	cards: number;
}) => {
	const { CardBackComponent } = React.useContext(ThemeContext).theme;
	return <div className={styles.opponentHandWrapper}>
		<div className={styles.opponentHand}>
			{[...Array(props.cards).keys()].map(index => <CardBackComponent key={index} />)}
		</div>
	</div>;
};