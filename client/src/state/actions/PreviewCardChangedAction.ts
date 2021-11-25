import { Action } from "redux";
import { ActionType } from "./ActionType";
import { lovelove } from "../../rpc/proto/lovelove";

export interface PreviewCardChangedAction extends Action<ActionType.PreviewCardChanged> {
	card?: lovelove.ICard;
}
