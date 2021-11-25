import { Action } from "../actions/Action";
import { ActionType } from "../actions/ActionType";
import { IState } from "../IState";

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
		} case ActionType.PreviewCardChanged: {
			return {
				...state,
				previewCard: action.card
			};
		}
	}

	return state;
}

export default mainReducer;