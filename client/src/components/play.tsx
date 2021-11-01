import * as React from "react";
import { ThemeContext } from "../themes/ThemeContext"
interface ICard {
	revealed?: boolean;
	season: number;
	variation: number;
}

export namespace Play {
	function Card({ card }: { card: ICard }) {
		const { CardComponent } = React.useContext(ThemeContext).theme;
		return <CardComponent card={card.variation} season={card.season} />
	}

	function Center(props: {
		cards: ICard[]
	}): JSX.Element {
		return <div className={`${""}`}/>
	}

	function PlayerArea(props: {
		cards: ICard[]
	}): JSX.Element {
		return <div className={`${""}`}>
			{ props.cards.map(card => <Card card={card} />) }
		</div>
	}

	export function Table(): JSX.Element {
		return <div>
			<PlayerArea cards={[
				{
					season: 2,
					variation: 2
				},
				{
					season: 5,
					variation: 1
				},
				{
					season: 8,
					variation: 3
				},
				{
					season: 0,
					variation: 0
				},
				{
					season: 11,
					variation: 2
				},
			]}/>
			<Center cards={[]}/>
			<PlayerArea cards={[]}/>
		</div>
	}
}
