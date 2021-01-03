import { Root, Method, Type } from "protobufjs";
export class Codec {
	public static stripIndex(data: Buffer): {
		index: number;
		data: Buffer;
	} {
		return {
			index: data[0] | data[1] << 8,
			data: data.slice(2)
		};
	}

	public static addIndex(index: number, data: Uint8Array): Uint8Array {
		return new Uint8Array([
			index & 0xff,
			index >> 8,
			...data,
		]);
	}

	// public static decode(root: Root, wrapper: Type, data: Buffer): any {
	// 	const message = wrapper.decode(data);
	// 	const type = root.lookupType(message["name"]);
	// 	return type.decode(message["data"]);
	// }

	constructor(private readonly protobufRoot: Root) {
	}

	// public decode(data: Buffer): any {
	// 	return Codec.decode(this.protobufRoot, this.wrapper, data);
	// }

	// public decodeMessage(message: Buffer, methodName?: string): {
	// 	type: MessageType;
	// 	index?: number;
	// 	data: any;
	// } {
	// 	const { type, data: wrappedData } = Codec.stripMessageType(message);
	// 	if (type === MessageType.Notification) {
	// 		return {
	// 			type,
	// 			data: this.decode(wrappedData)
	// 		};
	// 	}
	// 	if (type !== MessageType.Response && type !== MessageType.Request) {
	// 		console.log(`Unknown Message Type ${type}`);
	// 		return;
	// 	}
	// 	const { index, data } = Codec.stripIndex(wrappedData);
	// 	const unwrappedMessage = this.wrapper.decode(data);
	// 	const method = this.lookupMethod(methodName || unwrappedMessage["name"]);
	// 	return {
	// 		type,
	// 		index,
	// 		data: this.protobufRoot.lookupType(type === MessageType.Response
	// 			? method.requestType
	// 			: method.requestType).decode(unwrappedMessage["data"])
	// 	};
	// }

	private lookupMethod(path: string): Method {
		const sections = path.split(".");
		const service = this.protobufRoot.lookupService(sections.slice(0, -1));
		const name = sections[sections.length - 1];
		return service.methods[name];
	}
}
