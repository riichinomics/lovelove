import { ApiStateChangedAction } from "./ApiStateChangedAction";
import { InitialGameStateReceivedAction } from "./InitialGameStateReceivedAction";
import { PreviewCardChangedAction } from "./PreviewCardChangedAction";

export type Action =
	ApiStateChangedAction
	| InitialGameStateReceivedAction
	| PreviewCardChangedAction;