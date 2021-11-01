import { Root, load } from "protobufjs";
import { Connection } from "./Connection";
import { RpcImplementation } from "./RpcImplementation";
import { lovelove } from "./proto/lovelove";

export interface ApiOptions {
	url: string;
}

export class Api {
	// private readonly notifications: Observable<any>;

	private connection: Connection;
	public lovelove: lovelove.LoveLove;
	private protobufRoot: Root;
	private rpc: RpcImplementation;

	constructor(private readonly options: ApiOptions) {
		Root.prototype.fetch = function (filename, callback) {
			fetch(`http://${options.url}${filename}`)
				.then(
					response => response.text()
						.then(data => {
							console.log(data);
							callback(null, data);
						}),
					error => callback(error)
				);
		};
		// this.notifications = this.connection.messages.pipe(filter(message => message.index !== 0), map(message => this.codec.decode(message.data)));
	}

	public async init(): Promise<void> {
		this.protobufRoot = await load("/proto/lovelove.proto");
		const Wrapper = this.protobufRoot.lookupType("Wrapper") as unknown as typeof lovelove.Wrapper;
		this.connection = new Connection(`ws://${this.options.url}/echo`, Wrapper);
		this.rpc = new RpcImplementation(this.connection, this.protobufRoot);
		this.lovelove = this.rpc.createService<lovelove.LoveLove>("lovelove.LoveLove");
		await this.connection.init();
	}

	public dispose() {
		this.connection.close();
	}
}
