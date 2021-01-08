import { Subject, Observable } from 'rxjs';
import { lovelove } from './proto/lovelove';

export class Connection {
	private readonly messagesSubject = new Subject<lovelove.Wrapper>();
	private socket: WebSocket;

	constructor(
		private readonly server: string,
		private readonly Wrapper: typeof lovelove.Wrapper,
	) {}

	public get messages(): Observable<lovelove.Wrapper> {
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
				this.messagesSubject.next(this.Wrapper.create(event.data));
			};

			this.socket.onerror = (event: any) => {
				console.log(`websocker onerror`, event);
			}
			this.socket.onclose = (event: any) => {
				console.log(`websocker onclose`, event);
			}
			this.socket.onopen = () => resolve();
		});
	}

	public send(sequence: number, type: string, data: Uint8Array): void {
		if(this.socket.readyState !== WebSocket.OPEN) {
			throw new Error("Connection is not opened");
		}

		this.socket.send(this.Wrapper.encode({
			type,
			sequence,
			data
		}).finish());
	}

	public close() {
		this.socket.close()
	}
}
