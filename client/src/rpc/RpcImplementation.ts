import { Root, Method, RPCImplCallback, Type, rpc } from "protobufjs";
import { Subscription } from 'rxjs';
import { Codec } from "./Codec";
import { Connection } from "./Connection";

export class RpcImplementation {
	private readonly transactionMap: {
		[key: number]: protobuf.RPCImplCallback;
	} = {};

	private readonly dataSubscription: Subscription;
	private index = 0;

	constructor(private readonly connection: Connection, private readonly protobufRoot: Root) {
		this.dataSubscription = connection.messages.subscribe((message) => {
			if (!message.index) {
				return;
			}
			const { index, data } = message;

			const callback = this.transactionMap[index];
			delete this.transactionMap[index];
			if (!callback) {
				return;
			}
			callback(null, data);
		});
	}

	public createService<T extends rpc.Service>(service: string): T {
		return this.protobufRoot.lookupService(service).create(this.rpcCall.bind(this), false, false) as T;
	}

	private rpcCall(method: Method, requestData: Uint8Array, callback: RPCImplCallback) {
		const index = this.index++ % 60006 + 1;
		this.transactionMap[index] = callback;
		this.connection.send(Codec.addIndex(index, requestData));
	}
}
