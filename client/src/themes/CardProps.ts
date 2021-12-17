import { lovelove } from "../rpc/proto/lovelove";

export type CardProps = lovelove.ICard & {
	className?: string,
	hideHints?: boolean,
}
