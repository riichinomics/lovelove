import { ApiStateChangedAction } from "./ApiStateChangedAction";
import { GameUpdateReceivedAction } from "./GameUpdateReceivedAction";
import { GameCreatedOnServerAction } from "./ConnectedToGameAction";
import { PreviewCardChangedAction } from "./PreviewCardChangedAction";
import { RoundEndClearedAction } from "./RoundEndClearedAction";
import { EnteredNewRoomAction } from "./EnteredNewRoomAction";

export type Action =
	ApiStateChangedAction
	| GameCreatedOnServerAction
	| EnteredNewRoomAction
	| GameUpdateReceivedAction
	| PreviewCardChangedAction
	| RoundEndClearedAction;