import { ApiState } from "../rpc/ApiState";
import { lovelove } from "../rpc/proto/lovelove";

export interface IState {
	userId: string;
	apiState: ApiState;
	gameState?: lovelove.ICompleteGameState;
}
