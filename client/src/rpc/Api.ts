import { Root, load, rpc } from "protobufjs";
import { Observable, merge } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { Codec } from "./Codec";
import { Connection } from "./Connection";
import { RpcImplementation } from "./RpcImplementation";
import { lovelove } from "./proto/lovelove";

export interface ApiOptions {
	url: string;
}

export class Api {
	private readonly contestObservables: Observable<any>[] = [];
	private readonly notifications: Observable<any>;

	private readonly connection: Connection;
	public lovelove: lovelove.LoveLove;
	private protobufRoot: Root;
	private codec: Codec;
	private rpc: RpcImplementation;

	constructor(private readonly options: ApiOptions) {
		Root.prototype.fetch = function (filename, callback) {
			fetch(`http://${options.url}${filename}`)
				.then(
					response => response.text()
						.then(data => {
							console.log(data);
							callback(null, data)
						}),
					error => callback(error)
				);
		}

		this.connection = new Connection(`ws://${options.url}/echo`);
		// this.notifications = this.connection.messages.pipe(filter(message => message.index !== 0), map(message => this.codec.decode(message.data)));
	}

	public async init(): Promise<void> {
		this.protobufRoot = await load("/proto/lovelove.proto");
		this.rpc = new RpcImplementation(this.connection, this.protobufRoot);
		this.lovelove = this.rpc.createService<lovelove.LoveLove>("lovelove.LoveLove");
		this.codec = new Codec(this.protobufRoot);
		await this.connection.init();
	}

	public dispose() {
		this.connection.close();
	}
}
