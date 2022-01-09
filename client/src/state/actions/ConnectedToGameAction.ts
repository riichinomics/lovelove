import { Action } from "redux";
import { lovelove } from "../../rpc/proto/lovelove";
import { ActionType } from "./ActionType";

export interface GameCreatedOnServerAction extends Action<ActionType.ConnectedToGame> {
	position: lovelove.PlayerPosition,
	opponentDisconnected: boolean;
}