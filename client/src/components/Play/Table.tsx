import * as React from "react";
import { Center } from "./Center";
import { PlayerArea } from "./PlayerArea";

export const Table = () => <div>
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
	]}
	/>
	<Center cards={[]} />
	<PlayerArea cards={[]} />
</div>;
