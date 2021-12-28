import { ApiStateChangedAction } from "./ApiStateChangedAction";
import { GameUpdateReceivedAction } from "./GameUpdateReceivedAction";
import { GameCreatedOnServerAction } from "./GameCreatedOnServerAction";
import { PreviewCardChangedAction } from "./PreviewCardChangedAction";
import { RoundEndClearedAction } from "./RoundEndClearedAction";

export type Action =
	ApiStateChangedAction
	| GameCreatedOnServerAction
	| GameUpdateReceivedAction
	| PreviewCardChangedAction
	| RoundEndClearedAction;