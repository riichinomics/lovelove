import { stylesheet } from "astroturf";
import * as React from "react";
import { lovelove } from "../../rpc/proto/lovelove";
import { ThemeContext } from "../../themes/ThemeContext";
import { jpNumeral } from "../../utils";

const styles = stylesheet`
	.deckContainer {
		position: relative;
		.floatingCard {
			position: absolute;
			.cardOverlay {
				position: absolute;
				left: 0;
				right: 0;
				top: 0;
				bottom: 0;
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;

				.monthDisplay {
					display: flex;
					align-items: center;
					justify-content: center;

					> *:not(:first-child) {
						margin-left: 12px;
					}

					.monthName {
						font-size: 16px;
						writing-mode: vertical-rl;
					}
					.monthCard {
						width: 50px;
					}
				}
			}
		}
	}
`;


const paddingLength = 10;

export const Deck = (props: {
	cards: number,
	month: lovelove.Month,
	monthHana: lovelove.Hana,
}) => {
	const {theme} = React.useContext(ThemeContext);

	return <div className={styles.deckContainer} style={{paddingBottom: paddingLength * 2, paddingRight: paddingLength * 2}}>
		<theme.CardBackComponent />
		<div className={styles.floatingCard} style={{top: paddingLength, left: paddingLength}}><theme.CardBackComponent /></div>
		<div className={styles.floatingCard} style={{top: paddingLength * 2, left: paddingLength * 2}}>
			<theme.CardBackComponent />
			{props.monthHana &&
				<div className={styles.cardOverlay}>
					<div className={styles.monthDisplay}>
						<div className={styles.monthCard}>
							<theme.CardComponent hana={props.monthHana} variation={lovelove.Variation.First} hideHints />
						</div>
						<div className={styles.monthName}>
							{jpNumeral(props.month)}
							月
						</div>
					</div>
				</div>
			}
		</div>
	</div>;
};
