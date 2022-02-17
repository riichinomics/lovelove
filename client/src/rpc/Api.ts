import { Root, load, Message } from "protobufjs";
import { Connection } from "./Connection";
import { RpcImplementation } from "./RpcImplementation";
import { lovelove } from "./proto/lovelove";
import { first, firstValueFrom, map, Observable, share, takeUntil } from "rxjs";

export interface ApiOptions {
	production: boolean;
	url: string;
}

export interface ApiConnection {
	lovelove: lovelove.LoveLove;
	broadcastMessages: Observable<Message>;
	closed: Promise<{
		reconnect(): Promise<ApiConnection>
	}>;
}

export class Api {
	private protobufRoot: Root;

	constructor(private readonly options: ApiOptions) {
		Root.prototype.fetch = function (filename, callback) {
			fetch(`${options.production ? "https" : "http"}://${options.url}${options.production ? "" : ":6482"}${filename}`)
				.then(
					response => response.text()
						.then(data => {
							console.log(data);
							callback(null, data);
						}),
					error => callback(error)
				);
		};
	}

	public async init(): Promise<void> {
		this.protobufRoot = await load("/proto/lovelove.proto");
	}

	public async connect(): Promise<ApiConnection> {
		const Wrapper = this.protobufRoot.lookupType("Wrapper") as unknown as typeof lovelove.Wrapper;

		const connection = new Connection(`${this.options.production ? "wss" : "ws"}://${this.options.url}${this.options.production ? "" : ":6482"}/socket`, Wrapper);
		const rpc = new RpcImplementation(connection, this.protobufRoot);
		const loveloveService = rpc.createService<lovelove.LoveLove>("lovelove.LoveLove");

		return this.reconnect(connection, rpc, loveloveService);
	}

	private async reconnect(
		connection: Connection,
		rpc: RpcImplementation,
		loveloveService: lovelove.LoveLove,
	): Promise<ApiConnection> {
		while (await Promise.race([
			connection.reconnect(),
			new Promise((resolve) => setTimeout(() => resolve(true), 5000))
		]) === true) {
			// eslint-disable-next-line no-empty
		}

		const disconnectObservable = connection.disconnect.pipe(first(), share());
		const disconnectPromise = firstValueFrom(disconnectObservable).then(() => ({
			reconnect: () => this.reconnect(connection, rpc, loveloveService)
		}));

		return {
			broadcastMessages: rpc.broadcastMessages.pipe(
				map(wrapper => this.protobufRoot.lookupType(wrapper.contentType).decode(wrapper.data)),
				takeUntil(disconnectObservable)
			),
			lovelove: loveloveService,
			closed: disconnectPromise
		};
	}
}
