import { ApiStateChangedAction } from "./ApiStateChangedAction";
import { GameUpdateReceivedAction } from "./GameUpdateReceivedAction";
import { InitialGameStateReceivedAction } from "./InitialGameStateReceivedAction";
import { PreviewCardChangedAction } from "./PreviewCardChangedAction";

export type Action =
	ApiStateChangedAction
	| InitialGameStateReceivedAction
	| GameUpdateReceivedAction
	| PreviewCardChangedAction;