import * as React from "react";
import { ICard } from "../ICard";
import { ThemeContext } from "../../themes/ThemeContext";
import { cardKey } from "./utils";
import clsx from "clsx";
import { stylesheet } from "astroturf";

const styles = stylesheet`
	.collectionGroup {
		display: flex;
		align-items: flex-start;
		position: relative;
		&.upwards {
			flex-direction: column-reverse;
		}
	}

	.collectionGroupItem {
		position: absolute;
	}
`;

export const CardStack = (props: {
	cards: ICard[],
	stackUpwards?: boolean,
	concealed?: boolean,
	stackDepth?: number,
}) => {
	const { stackDepth = 1 } = props;
	const {
		CardComponent,
		CardPlaceholderComponent,
		CardBackComponent,
		cardStackSpacing
	} = React.useContext(ThemeContext).theme;
	const [selectedIndex, setSelectedIndex] = React.useState(0);
	const [selectedLayerIndex, setSelectedLayerIndex] = React.useState(0);

	const layers = React.useMemo(() => {
		const layers: ICard[][] = [];
		for (let i = 0; i < props.cards.length; i++) {
			if (i % stackDepth === 0) {
				layers.push([]);
			}
			layers[layers.length - 1].push(props.cards[i]);
		}
		return layers;
	}, [props.cards, stackDepth]);

	const cardStackHorizontalSpacing = 30;
	const cardStackVerticalSpacing = 30;
	const cardStackLayerOffset = 20;

	const horizontalPadding = cardStackHorizontalSpacing * (layers[0].length - 1) + cardStackLayerOffset * layers.slice(1).filter(layer => layer.length >= stackDepth).length;
	const verticalPadding = cardStackVerticalSpacing * (layers.length - 1);

	return <div
		className={clsx(styles.collectionGroup, props.stackUpwards && styles.upwards)}
		style={{
			paddingRight: horizontalPadding,
			paddingBottom: props.stackUpwards ? null : verticalPadding,
			paddingTop: props.stackUpwards ? verticalPadding : null,
			zIndex: props.cards.length
		}}
	>
		{layers.map((layer, layerIndex) => layer.map((card, index) => {
			const distanceToSelectedLayerIndex = Math.abs(selectedLayerIndex - layerIndex);
			const distanceToSelectedIndex = Math.abs(selectedIndex - index);
			return <div
				className={clsx((index !== 0 || layerIndex !== 0) && styles.collectionGroupItem)}
				key={cardKey(card, index)}
				style={{
					top: props.stackUpwards ? null : cardStackVerticalSpacing * layerIndex,
					left: cardStackHorizontalSpacing * index + cardStackLayerOffset * layerIndex,
					bottom: props.stackUpwards ? cardStackVerticalSpacing * layerIndex : null,
					zIndex: layers.length - distanceToSelectedLayerIndex + (distanceToSelectedLayerIndex === 0
						? layer.length - distanceToSelectedIndex
						: 0
					)
				}}
				onMouseEnter={() => {
					setSelectedIndex(index);
					setSelectedLayerIndex(layerIndex);
				}}
			>
				{props.concealed
					? <CardBackComponent />
					: card
						? <CardComponent card={card.variation} season={card.season} />
						: <CardPlaceholderComponent />
				}
			</div>;
		}))}
	</div>;
};
