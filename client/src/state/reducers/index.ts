import {produce, immerable} from "immer";
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

function applyYakuUpdate(yakuInformation: lovelove.IYakuData[], yakuUpdate: lovelove.IYakuUpdate) {
	for (const deletedYaku of yakuUpdate.deletedYaku) {
		const indexOfYaku = yakuInformation.findIndex(yaku => yaku.id == deletedYaku);
		if (indexOfYaku >= 0) {
			yakuInformation.splice(indexOfYaku, 1);
		}
	}

	for (const newOrUpdatedYaku of yakuUpdate.newOrUpdatedYaku) {
		const existingYaku = yakuInformation.find(yaku => yaku.id == newOrUpdatedYaku.yakuId);
		if (!existingYaku) {
			yakuInformation.push({
				id: newOrUpdatedYaku.yakuId,
				cards: newOrUpdatedYaku.cardIds,
				value: newOrUpdatedYaku.value,
			});
			continue;
		}

		existingYaku.cards.push(...newOrUpdatedYaku.cardIds);
		existingYaku.value = newOrUpdatedYaku.value;
	}
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
			(action.gameState as any)[immerable] = true;
			return {
				...state,
				gamePosition: action.position,
				gameState: action.gameState
			};
		} case ActionType.GameUpdateReceived: {
			return {
				...state,
				gameState: produce(state.gameState, gameState => {
					for (const update of action.update.updates) {
						if (update.cardMoveUpdates) {
							for (const cardMove of update.cardMoveUpdates) {
								switch(cardMove.originSlot.zone) {
									case lovelove.PlayerCentricZone.UnknownZone: {
										break;
									}
									case lovelove.PlayerCentricZone.Table: {
										if (gameState.table[cardMove.originSlot.index ?? 0]?.card) {
											gameState.table[cardMove.originSlot.index ?? 0].card = null;
										}
										break;
									}
									case lovelove.PlayerCentricZone.Hand: {
										gameState.hand = removeCard(gameState.hand, cardMove.movedCard.id);
										break;
									}
									case lovelove.PlayerCentricZone.OpponentHand: {
										gameState.opponentHand--;
										break;
									}
									case lovelove.PlayerCentricZone.Deck: {
										gameState.deck--;
										break;
									}
									case lovelove.PlayerCentricZone.Collection: {
										gameState.collection = removeCard(gameState.collection, cardMove.movedCard.id);
										break;
									}
									case lovelove.PlayerCentricZone.OpponentCollection: {
										gameState.opponentCollection = removeCard(gameState.opponentCollection, cardMove.movedCard.id);
										break;
									}
									case lovelove.PlayerCentricZone.Drawn: {
										gameState.deckFlipCard = null;
										break;
									}
								}

								switch(cardMove.destinationSlot.zone) {
									case lovelove.PlayerCentricZone.UnknownZone: {
										break;
									}
									case lovelove.PlayerCentricZone.Table: {
										// TODO: Animation Float
										gameState.table[cardMove.destinationSlot.index ?? 0] = {
											card: cardMove.movedCard
										};
										break;
									}
									case lovelove.PlayerCentricZone.Hand: {
										gameState.hand = [...gameState.hand ?? [], cardMove.movedCard];
										break;
									}
									case lovelove.PlayerCentricZone.OpponentHand: {
										gameState.opponentHand++;
										break;
									}
									case lovelove.PlayerCentricZone.Deck: {
										gameState.deck++;
										break;
									}
									case lovelove.PlayerCentricZone.Collection: {
										gameState.collection = [...gameState.collection ?? [], cardMove.movedCard];
										break;
									}
									case lovelove.PlayerCentricZone.OpponentCollection: {
										gameState.opponentCollection = [...gameState.opponentCollection ?? [], cardMove.movedCard];
										break;
									}
									case lovelove.PlayerCentricZone.Drawn: {
										gameState.deckFlipCard = cardMove.movedCard;
										break;
									}
								}
							}
						}

						if (update.playOptionsUpdate) {
							if (update.playOptionsUpdate.defunctOptions) {
								for (const defuctOption of update.playOptionsUpdate.defunctOptions) {
									if (!defuctOption.originCardId) {
										delete gameState.tablePlayOptions.playOptions[defuctOption.targetCardId.cardId];
										continue;
									}

									if (!defuctOption.targetCardId) {
										const originOptionIndex = gameState.tablePlayOptions.noTargetPlayOptions.options.indexOf(defuctOption.originCardId.cardId);
										if (originOptionIndex >= 0) {
											gameState.tablePlayOptions.noTargetPlayOptions.options.splice(originOptionIndex, 1);
										}
										continue;
									}

									const optionIndex = gameState.tablePlayOptions.playOptions[defuctOption.targetCardId.cardId].options.indexOf(defuctOption.originCardId.cardId);
									gameState.tablePlayOptions.noTargetPlayOptions.options.splice(optionIndex, 1);
								}
							}

							if (update.playOptionsUpdate.newOptions) {
								for (const newOption of update.playOptionsUpdate.newOptions) {
									if (!newOption.targetCardId) {
										gameState.tablePlayOptions.noTargetPlayOptions.options.push(newOption.originCardId.cardId);
										continue;
									}

									if (!gameState.tablePlayOptions.playOptions[newOption.targetCardId.cardId]) {
										gameState.tablePlayOptions.playOptions[newOption.targetCardId.cardId] = {
											options: []
										};
									}

									gameState.tablePlayOptions.playOptions[newOption.targetCardId.cardId].options.push(newOption.originCardId.cardId);
								}
							}

							if (update.playOptionsUpdate.updatedAcceptedOriginZones) {
								gameState.tablePlayOptions.acceptedOriginZones = update.playOptionsUpdate.updatedAcceptedOriginZones.zones;
							}
						}

						if (update.yakuUpdate) {
							applyYakuUpdate(gameState.yakuInformation, update.yakuUpdate);
						}

						if (update.opponentYakuUpdate) {
							applyYakuUpdate(gameState.opponentYakuInformation, update.opponentYakuUpdate);
						}

						if (update.shoubuOpportunityUpdate) {
							gameState.shoubuOpportunity = update.shoubuOpportunityUpdate.available;
						}

						if (update.activePlayerUpdate) {
							gameState.active = update.activePlayerUpdate.position;
						}

						if (update.koikoiUpdate) {
							if (update.koikoiUpdate.self) {
								gameState.koikoi = true;
							}

							if (update.koikoiUpdate.opponent) {
								gameState.opponentKoikoi = true;
							}
						}
					}
				}),
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
