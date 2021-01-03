import { Subject, Observable } from 'rxjs';
import { Codec } from "./Codec";

export class Connection {
	private readonly messagesSubject = new Subject<any>();
	private socket: WebSocket;

	constructor(private readonly server: string) {}

	public get messages(): Observable<{
		index: number;
		data: Buffer;
	}> {
		return this.messagesSubject;
	}

	public init(): Promise<void> {
		return this.reconnect();
	}

	private reconnect(): Promise<void> {
		if (this.socket) {
			this.socket.close();
		}

		console.log("Connecting to " + this.server);

		return new Promise((resolve) => {
			this.socket = new WebSocket(this.server);
			this.socket.onmessage = (event: MessageEvent) => {
				console.log(event.data);
				const message = Codec.stripIndex(event.data);
				this.messagesSubject.next(message);
			};

			this.socket.onerror = (event: any) => {
				console.log(`websocker onerror`, event);
				process.exit(1);
			}
			this.socket.onclose = (event: any) => {
				console.log(`websocker onclose`, event);
				process.exit(1);
			}
			this.socket.onopen = () => resolve();
		});
	}

	public send(data: Uint8Array): void {
		if(this.socket.readyState !== WebSocket.OPEN) {
			throw new Error("Connection is not opened");
		}

		this.socket.send(data);
	}

	public close() {
		this.socket.close()
	}
}
