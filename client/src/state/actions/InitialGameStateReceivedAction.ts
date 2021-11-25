import { Action } from "redux";
import { ActionType } from "./ActionType";
import { lovelove } from "../../rpc/proto/lovelove";

export interface InitialGameStateReceivedAction extends Action<ActionType.InitialGameStateReceived> {
	position?: lovelove.PlayerPosition;
	gameState?: lovelove.ICompleteGameState;
}