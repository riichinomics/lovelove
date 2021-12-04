import { lovelove } from "../../rpc/proto/lovelove";
import { Action } from "../actions/Action";
import { ActionType } from "../actions/ActionType";
import { IState } from "../IState";

function removeCard(zone: lovelove.ICard[], cardId: number, leaveSpace?: boolean): lovelove.ICard[] {
	if (!zone) {
		return zone;
	}

	const index = zone.findIndex(card => card?.id === cardId);
	if (index < 0) {
		return zone;
	}

	if (leaveSpace) {
		return [...zone.slice(0, index), null, ...zone.slice(index + 1)];
	}

	return [...zone.slice(0, index), ...zone.slice(index + 1)];
}

function mainReducer(state: IState, action: Action): IState {
	// eslint-disable-next-line no-empty
	switch (action.type) {
		case ActionType.ApiStateChanged: {
			return {
				...state,
				apiState: action.apiState
			};
		} case ActionType.InitialGameStateReceived: {
			return {
				...state,
				gamePosition: action.position,
				gameState: action.gameState
			};
		} case ActionType.GameUpdateReceived: {
			const gameStateUpdate = action.update;
			const updatedGameState = {...state.gameState};
			for (const update of gameStateUpdate.updates) {
				if (update.cardMoveUpdates) {
					for (const cardMove of update.cardMoveUpdates) {
						switch(cardMove.originSlot.zone) {
							case lovelove.PlayerCentricZone.UnknownZone: {
								break;
							}
							case lovelove.PlayerCentricZone.Table: {
								updatedGameState.table = removeCard(updatedGameState.table, cardMove.movedCard.id, true);
								break;
							}
							case lovelove.PlayerCentricZone.Hand: {
								updatedGameState.hand = removeCard(updatedGameState.hand, cardMove.movedCard.id);
								break;
							}
							case lovelove.PlayerCentricZone.OpponentHand: {
								updatedGameState.opponentHand--;
								break;
							}
							case lovelove.PlayerCentricZone.Deck: {
								updatedGameState.deck--;
								break;
							}
							case lovelove.PlayerCentricZone.Collection: {
								updatedGameState.collection = removeCard(updatedGameState.collection, cardMove.movedCard.id);
								break;
							}
							case lovelove.PlayerCentricZone.OpponentCollection: {
								updatedGameState.opponentCollection = removeCard(updatedGameState.opponentCollection, cardMove.movedCard.id);
								break;
							}
							case lovelove.PlayerCentricZone.Drawn: {
								updatedGameState.deckFlipCard = null;
								break;
							}
						}

						switch(cardMove.destinationSlot.zone) {
							case lovelove.PlayerCentricZone.UnknownZone: {
								break;
							}
							case lovelove.PlayerCentricZone.Table: {
								//TODO: animation float
								break;
							}
							case lovelove.PlayerCentricZone.Hand: {
								updatedGameState.hand = [...updatedGameState.hand ?? [], cardMove.movedCard];
								break;
							}
							case lovelove.PlayerCentricZone.OpponentHand: {
								updatedGameState.opponentHand++;
								break;
							}
							case lovelove.PlayerCentricZone.Deck: {
								updatedGameState.deck++;
								break;
							}
							case lovelove.PlayerCentricZone.Collection: {
								updatedGameState.collection = [...updatedGameState.collection ?? [], cardMove.movedCard];
								break;
							}
							case lovelove.PlayerCentricZone.OpponentCollection: {
								updatedGameState.opponentCollection = [...updatedGameState.opponentCollection ?? [], cardMove.movedCard];
								break;
							}
							case lovelove.PlayerCentricZone.Drawn: {
								updatedGameState.deckFlipCard = cardMove.movedCard;
								break;
							}
						}
					}
				}
			}
			return {
				...state,
				gameState: updatedGameState,
			};
		} case ActionType.PreviewCardChanged: {
			return {
				...state,
			};
		}
	}

	return state;
}

export default mainReducer;