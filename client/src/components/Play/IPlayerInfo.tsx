import { lovelove } from "../../rpc/proto/lovelove";

export interface IPlayerInfo {
	position: lovelove.PlayerPosition;
	oya?: boolean;
	name?: string;
	active?: boolean;
}
