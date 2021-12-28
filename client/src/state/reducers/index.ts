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

function immerate<T>(object: T): T{
	(object as any)[immerable] = true;
	return object;
}

function mainReducer(state: IState, action: Action): IState {
	// eslint-disable-next-line no-empty
	switch (action.type) {
		case ActionType.ApiStateChanged: {
			return {
				...state,
				apiState: action.apiState
			};
		} case ActionType.GameUpdateReceived: {
			return produce(state, state => {
				const gameState = state.gameState;
				const player = state.gamePosition == lovelove.PlayerPosition.Red ? gameState?.redPlayer : gameState?.whitePlayer;
				const opponent = state.gamePosition == lovelove.PlayerPosition.Red ? gameState?.whitePlayer : gameState?.redPlayer;
				for (const update of action.update.updates) {
					if (update.cardMoveUpdates) {
						for (const cardMove of update.cardMoveUpdates) {
							if (cardMove.originSlot.player == lovelove.PlayerPosition.UnknownPosition) {
								switch(cardMove.originSlot.zone) {
									case lovelove.CardZone.Table: {
										if (gameState.table[cardMove.originSlot.index ?? 0]?.card) {
											gameState.table[cardMove.originSlot.index ?? 0].card = null;
										}
										break;
									}
									case lovelove.CardZone.Deck: {
										gameState.deck--;
										break;
									}
									case lovelove.CardZone.Drawn: {
										gameState.deckFlipCard = null;
										break;
									}
								}
							} else if (cardMove.originSlot.player == state.gamePosition) {
								switch(cardMove.originSlot.zone) {
									case lovelove.CardZone.Hand: {
										player.hand.numberOfCards--;
										player.hand.cards = removeCard(player.hand.cards, cardMove.movedCard.id);
										break;
									}
									case lovelove.CardZone.Collection: {
										player.collection = removeCard(player.collection, cardMove.movedCard.id);
										break;
									}
								}
							} else {
								switch(cardMove.originSlot.zone) {
									case lovelove.CardZone.Hand: {
										opponent.hand.numberOfCards--;
										break;
									}
									case lovelove.CardZone.Collection: {
										opponent.collection = removeCard(opponent.collection, cardMove.movedCard.id);
										break;
									}
								}
							}

							if (cardMove.destinationSlot.player == lovelove.PlayerPosition.UnknownPosition) {
								switch(cardMove.destinationSlot.zone) {
									case lovelove.CardZone.Table: {
										// TODO: Animation Float
										gameState.table[cardMove.destinationSlot.index ?? 0] = {
											card: cardMove.movedCard
										};
										break;
									}
									case lovelove.CardZone.Deck: {
										gameState.deck++;
										break;
									}
									case lovelove.CardZone.Drawn: {
										gameState.deckFlipCard = cardMove.movedCard;
										break;
									}
								}
							} else if (cardMove.destinationSlot.player == state.gamePosition) {
								switch(cardMove.destinationSlot.zone) {
									case lovelove.CardZone.Hand: {
										player.hand.numberOfCards++;
										player.hand.cards = [...player.hand.cards, cardMove.movedCard];
										break;
									}
									case lovelove.CardZone.Collection: {
										player.collection = [...player.collection, cardMove.movedCard];
										break;
									}
								}
							} else {
								switch(cardMove.destinationSlot.zone) {
									case lovelove.CardZone.Hand: {
										opponent.hand.numberOfCards++;
										break;
									}
									case lovelove.CardZone.Collection: {
										opponent.collection = [...opponent.collection, cardMove.movedCard];
										break;
									}
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
						applyYakuUpdate(player.yakuInformation, update.yakuUpdate);
					}

					if (update.opponentYakuUpdate) {
						applyYakuUpdate(opponent.yakuInformation, update.opponentYakuUpdate);
					}

					if (update.shoubuOpportunityUpdate) {
						if (!update.shoubuOpportunityUpdate.available) {
							gameState.shoubuOpportunity = null;
						} else {
							gameState.shoubuOpportunity = {
								value: update.shoubuOpportunityUpdate.value ?? 0
							};
						}
					}

					if (update.activePlayerUpdate) {
						gameState.active = update.activePlayerUpdate.position;
					}

					if (update.koikoiUpdate) {
						if (update.koikoiUpdate.self) {
							player.koikoi = true;
						}

						if (update.koikoiUpdate.opponent) {
							opponent.koikoi = true;
						}
					}

					if (update.roundEndResult) {
						state.roundEndView = {
							gameState: gameState,
							winner: update.roundEndResult.winner,
							winnings: update.roundEndResult.winnings,
							teyaku: update.roundEndResult.teyakuInformation,
						};

						state.gameState = immerate(update.roundEndResult.nextRound);

					}

					if (update.connectionStatusUpdate) {
						if (update.connectionStatusUpdate.player != state.gamePosition) {
							state.opponentDisconnected = !update.connectionStatusUpdate.connected;
						}
					}

					if (update.gameConnectionData) {
						state.gameState = immerate(update.gameConnectionData.gameState);
						state.gamePosition = update.gameConnectionData.position;
						state.opponentDisconnected = update.gameConnectionData.opponentDisconnected;
					}
				}
			});
		} case ActionType.PreviewCardChanged: {
			return {
				...state,
			};
		} case ActionType.RoundEndCleared: {
			return {
				...state,
				roundEndView: null,
			};
		} case ActionType.GameCreatedOnServer: {
			return {
				...state,
				gamePosition: null,
				gameState: null,
				opponentDisconnected: null,
			};
		}
	}

	return state;
}

export default mainReducer;
