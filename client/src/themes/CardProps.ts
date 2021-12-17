import { lovelove } from "../rpc/proto/lovelove";

export type CardProps = lovelove.ICard & {
	hideHints?: boolean,
}
