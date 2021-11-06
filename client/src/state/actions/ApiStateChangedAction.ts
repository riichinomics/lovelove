import { Action } from "redux";
import { ActionType } from "./ActionType";
import { ApiState } from "../../rpc/ApiState";

export interface ApiStateChangedAction extends Action<ActionType.ApiStateChanged> {
	apiState?: ApiState;
}
