import { Action } from "redux";
import { ActionType } from "./ActionType";
import { ApiConnection } from "../../rpc/Api";

export interface ApiStateChangedAction extends Action<ActionType.ApiStateChanged> {
	apiConnection?: ApiConnection;
}
