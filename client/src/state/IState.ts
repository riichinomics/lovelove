import { ApiState } from "../rpc/ApiState";
import { lovelove } from "../rpc/proto/lovelove";

export interface IState {
	apiState: ApiState;
	gameState?: lovelove.ICompleteGameState;
}
