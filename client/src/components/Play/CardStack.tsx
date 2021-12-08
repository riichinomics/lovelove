import * as React from "react";
import { ThemeContext } from "../../themes/ThemeContext";
import { CardDroppedHandler, cardKey, CardLocation, CardWithOffset, CardZone, Vector2 } from "./utils";
import clsx from "clsx";
import { lovelove } from "../../rpc/proto/lovelove";
import { stylesheet } from "astroturf";
import { DragDropTypes } from "./DragDropTypes";
import { useDrag, useDrop } from "react-dnd";

export const Card = (props: {
	card?: CardWithOffset,
} & ISharedCardProps) => {
	const {
		cardStackSpacing,
		CardComponent,
		CardPlaceholderComponent,
		CardBackComponent
	} = React.useContext(ThemeContext).theme;

	const wrapperRef = React.useRef<HTMLDivElement>(null);

	const [{ isDragging }, drag] = useDrag(() => ({
		type: DragDropTypes.Card,
		collect: (monitor) => ({
			isDragging: monitor.isDragging(),
		}),
		canDrag: () => props.canDrag,
		item: () => {
			const rect = wrapperRef.current.getBoundingClientRect();
			return {
				card: {
					zone: props.zone,
					index: props.index,
					card: {
						...(props.card ?? {})
					},
				},
				offset: {
					x: rect.left,
					y: rect.top,
				}
			};
		}
	}), [props.card, props.zone, props.canDrag, wrapperRef]);

	const [{canDrop, isOver}, drop] = useDrop(() => ({
		accept: DragDropTypes.Card,
		canDrop: (item: { card: CardLocation, offset: Vector2 }) => props.playOptions?.indexOf(item.card.card.id) >= 0,
		drop: (item) => {
			const rect = wrapperRef.current.getBoundingClientRect();
			props.onCardDropped?.(
				{
					from: item.card,
					to: {
						zone: props.zone,
						index: props.index,
						card: {
							...(props.card ?? {})
						},
					},
					offset: {
						x: item.offset.x - rect.left - cardStackSpacing.horizontal,
						y: item.offset.y - rect.top - cardStackSpacing.vertical
					}
				}
			);
		},
		collect: (monitor) => ({
			canDrop: monitor.canDrop(),
			isOver: monitor.isOver(),
		})
	}), [props.playOptions, props.onCardDropped, props.card, props.zone, cardStackSpacing]);

	const dragDropRef = drag(drop(wrapperRef));

	return <div
		ref={dragDropRef as any}
		className={clsx(
			styles.cardDropWrapper,
			(!isDragging && canDrop) || (props.card && props.previewCard && props.playOptions?.indexOf(props.previewCard.id) < 0) && styles.previewCardNotAccepted,
			isOver && canDrop && styles.cardHoverOver
		)}
		style={{
			transform: props.card?.offset && `translate(${props.card.offset.x}px, ${props.card.offset.y}px)`
		}}
	>
		{isDragging
			? <CardPlaceholderComponent />
			: props.concealed
				? <CardBackComponent />
				: props.card
					? <CardComponent hana={props.card.hana} variation={props.card.variation} />
					: <CardPlaceholderComponent />
		}
	</div>;
};

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

	.cardDropWrapper {
		position: relative;
		&:after {
			content: "";
			position: absolute;
			/* transition: top 0.3s ease-in-out, bottom 0.3s ease-in-out, right 0.3s ease-in-out, left 0.3s ease-in-out; */
			top: 0;
			bottom: 0;
			left: 0;
			right: 0;
			z-index: -1;
		}

		&.cardHoverOver {
			&:after {
				top: -4px;
				bottom: -4px;
				left: -4px;
				right: -4px;
				background-color: gold;
			}
		}

		&.previewCardNotAccepted {
			background-color: black;
			> * {
				opacity: 0.5;
			}
		}

		> * {
			transition: opacity 0.3s ease-in-out;
		}
	}

`;

interface ISharedCardProps {
	index?: number;
	playOptions?: number[];
	canDrag?: boolean;
	concealed?: boolean;
	zone?: CardZone;
	previewCard?: lovelove.ICard;
	onCardDropped?: CardDroppedHandler;
}

export const CardStack = (props: {
	cards: CardWithOffset[],
	stackUpwards?: boolean,
	stackDepth?: number,
	move?: CardLocation,
	onCardSelected?: (card: lovelove.ICard) => void,
	onMouseLeave?: () => void,
	stunted?: boolean;
	laminated?: boolean
} & ISharedCardProps ) => {
	const { stackDepth = 1 } = props;
	const [selectedIndexReal, setSelectedIndex] = React.useState(-1);
	const [selectedLayerIndexReal, setSelectedLayerIndex] = React.useState(-1);

	const layers = React.useMemo(() => {
		const layers: CardWithOffset[][] = [];
		for (let i = 0; i < props.cards.length; i++) {
			if (i % stackDepth === 0) {
				layers.push([]);
			}
			layers[layers.length - 1].push(props.cards[i]);
		}
		return layers;
	}, [props.cards, stackDepth]);

	const selectedLayerIndex = (layers.length + selectedLayerIndexReal) % layers.length;
	const selectedIndex = (layers[selectedLayerIndex]?.length ?? 0 + selectedIndexReal) % (layers[selectedLayerIndex]?.length ?? 1);

	const {
		cardStackSpacing,
	} = React.useContext(ThemeContext).theme;

	const cardStackHorizontalSpacing = cardStackSpacing.vertical ?? 30;
	const cardStackVerticalSpacing = cardStackSpacing.vertical ?? 30;
	const cardStackLayerOffset = cardStackSpacing.horizontal ?? 20;

	const horizontalPadding = cardStackHorizontalSpacing * ((layers[0]?.length ?? 1) - 1) + cardStackLayerOffset * layers.slice(1).filter(layer => layer.length >= stackDepth).length;
	const verticalPadding = cardStackVerticalSpacing * (layers.length - 1);

	return (
		<div
			className={clsx(styles.collectionGroup, props.stackUpwards && styles.upwards)}
			style={{
				paddingRight: !props.stunted && horizontalPadding,
				paddingBottom: !props.stunted && (props.stackUpwards ? null : verticalPadding),
				paddingTop: !props.stunted && (props.stackUpwards ? verticalPadding : null),
				zIndex: props.cards.length
			}}
			onMouseLeave={props.onMouseLeave}
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
						if (!props.concealed && !props.laminated) {
							setSelectedIndex(index);
							setSelectedLayerIndex(layerIndex);
							props.onCardSelected?.(layers[layerIndex][index]);
						}
					}}
				>
					<Card
						card={card}
						concealed={props.concealed}
						playOptions={props.playOptions}
						canDrag={props.canDrag}
						previewCard={props.previewCard}
						onCardDropped={props.onCardDropped}
						zone={props.zone}
						index={props.index}
					/>
				</div>;
			}))}
		</div>
	);
};
