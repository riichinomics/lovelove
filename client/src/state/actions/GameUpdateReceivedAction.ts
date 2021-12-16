import { Action } from "redux";
import { ActionType } from "./ActionType";
import { lovelove } from "../../rpc/proto/lovelove";

export interface GameUpdateReceivedAction extends Action<ActionType.GameUpdateReceived> {
	update: lovelove.GameStateUpdate;
}
