import { ApiStateChangedAction } from "./ApiStateChangedAction";
import { InitialGameStateReceivedAction } from "./InitialGameStateReceivedAction";

export type Action =
	ApiStateChangedAction
	| InitialGameStateReceivedAction;