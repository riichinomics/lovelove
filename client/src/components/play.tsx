import React from "react"
import styles from "./play.sass"

export namespace Play {
	function Card(props: {
	}): JSX.Element {
		return <div className={`${styles.card}`}>
			Card
		</div>
	}

	export function Table(): JSX.Element {
		return <div>
			<Card/>
		</div>
	}
}
