import { Method, RPCImplCallback, Root, rpc } from "protobufjs";
import { Connection } from "./Connection";
import { filter, Observable, share, Subscription } from "rxjs";
import { lovelove } from "./proto/lovelove";

export class RpcImplementation {
	public readonly broadcastMessages: Observable<lovelove.Wrapper>;

	private readonly transactionMap: {
		[key: number]: protobuf.RPCImplCallback;
	} = {};

	private readonly dataSubscription: Subscription;
	private index = 0;

	constructor(private readonly connection: Connection, private readonly protobufRoot: Root) {
		this.broadcastMessages = connection.messages.pipe(
			filter(message => message.type === lovelove.MessageType.Broadcast),
			share()
		);


		this.dataSubscription = connection.messages.pipe(
			filter(message => message.type === lovelove.MessageType.Transact),
		).subscribe((message) => {
			if (!message.sequence) {
				return;
			}

			const callback = this.transactionMap[message.sequence];
			delete this.transactionMap[message.sequence];
			if (!callback) {
				return;
			}
			callback(null, message.data);
		});
	}

	public createService<T extends rpc.Service>(service: string): T {
		return this.protobufRoot.lookupService(service).create(this.rpcCall.bind(this), false, false) as T;
	}

	private rpcCall(method: Method, requestData: Uint8Array, callback: RPCImplCallback) {
		const index = this.index++ % 60006 + 1;
		this.transactionMap[index] = callback;
		this.connection.send(index, method.fullName, requestData);
	}
}
