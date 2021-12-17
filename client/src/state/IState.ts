import { ApiState } from "../rpc/ApiState";
import { lovelove } from "../rpc/proto/lovelove";

export interface RoundEndInformation {
	winner: lovelove.PlayerPosition;
	winnings: number;
	gameState: lovelove.ICompleteGameState;
	teyaku: lovelove.IRoundEndResultTeyakuInformation[];
}

export interface IState {
	userId: string;
	apiState: ApiState;
	gamePosition: lovelove.PlayerPosition;
	gameState?: lovelove.ICompleteGameState;
	roundEndView?: RoundEndInformation;
}
