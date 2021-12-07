package rpc

import (
	"context"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/grpc"
	lovelove "hanafuda.moe/lovelove/proto"
)

type connContextKey struct {
	key string
}

type serviceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*grpc.MethodDesc
	streams     map[string]*grpc.StreamDesc
	mdata       interface{}
}

type connectionMeta struct {
	ConnId   string
	Messages chan proto.Message
	closed   chan rxgo.Item
	Closed   rxgo.Observable
}

// Connection between autogenerated RPC code and the websocket listener.
// Implements the grpc service registrar interface, then decodes and dispatches
// the messages from the web socket.
type webSocketRpcServer struct {
	// lis      map[net.Listener]bool
	// conns    map[transport.ServerTransport]bool
	serve    bool
	services map[string]*serviceInfo
}

func NewWebSocketRpcServer() *webSocketRpcServer {
	return &webSocketRpcServer{
		services: make(map[string]*serviceInfo),
	}
}

func (server *webSocketRpcServer) RegisterService(serviceDesc *grpc.ServiceDesc, ss interface{}) {
	// if server.serve {
	// 	logger.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", serviceDesc.ServiceName)
	// }

	if _, ok := server.services[serviceDesc.ServiceName]; ok {
		log.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", serviceDesc.ServiceName)
	}

	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*grpc.MethodDesc),
		streams:     make(map[string]*grpc.StreamDesc),
		mdata:       serviceDesc.Metadata,
	}
	for i := range serviceDesc.Methods {
		d := &serviceDesc.Methods[i]
		info.methods[d.MethodName] = d
	}
	for i := range serviceDesc.Streams {
		d := &serviceDesc.Streams[i]
		info.streams[d.StreamName] = d
	}
	server.services[serviceDesc.ServiceName] = info
}

func (server *webSocketRpcServer) HandleConnection(connection *websocket.Conn) {
	go func(conn *websocket.Conn) {
		connMeta := &connectionMeta{
			ConnId:   uuid.New().String(),
			Messages: make(chan proto.Message),
			closed:   make(chan rxgo.Item),
		}
		connMeta.Closed = rxgo.FromChannel(connMeta.closed, rxgo.WithPublishStrategy())
		connMeta.Closed.Connect(context.TODO())

		sendChannel := make(chan []byte)

		defer close(connMeta.Messages)
		defer close(connMeta.closed)
		defer close(sendChannel)

		go func(sendChannel chan []byte) {
			for data := range sendChannel {
				conn.WriteMessage(websocket.BinaryMessage, data)
			}
		}(sendChannel)

		go func(cm *connectionMeta, sendChannel chan []byte) {
			sequence := int32(0)
			for message := range cm.Messages {
				valueData, _ := proto.Marshal(message)

				wrapperData, _ := proto.Marshal(&lovelove.Wrapper{
					Type:        lovelove.MessageType_Broadcast,
					Sequence:    sequence,
					ContentType: string(message.ProtoReflect().Descriptor().Name()),
					Data:        valueData,
				})
				sequence = sequence + 1
				sendChannel <- wrapperData
			}
		}(connMeta, sendChannel)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Print(err)
				break
			}
			log.Print(messageType)
			log.Print(data)
			wrapper := new(lovelove.Wrapper)
			proto.Unmarshal(data, wrapper)
			log.Print(wrapper.Type, wrapper.Sequence, wrapper.ContentType, wrapper.Data)

			method := wrapper.ContentType
			if method != "" && method[0] == '.' {
				method = method[1:]
			}

			lastDotPosition := strings.LastIndex(method, ".")

			if lastDotPosition == -1 {
				log.Print("Malformed method ", wrapper.Type)
				continue
			}

			serviceName := method[:lastDotPosition]
			methodName := method[lastDotPosition+1:]

			serviceInfo, serviceIsKnown := server.services[serviceName]
			if serviceIsKnown {
				if methodInfo, ok := serviceInfo.methods[methodName]; ok {
					value, _ := methodInfo.Handler(
						serviceInfo.serviceImpl,
						context.WithValue(context.Background(), connContextKey{
							key: "connId",
						}, connMeta),
						func(message interface{}) error {
							return proto.Unmarshal(wrapper.Data, message.(proto.Message))
						},
						nil,
					)

					valueData, _ := proto.Marshal(value.(proto.Message))

					wrapperData, _ := proto.Marshal(&lovelove.Wrapper{
						Sequence:    wrapper.Sequence,
						Type:        lovelove.MessageType_Transact,
						ContentType: wrapper.ContentType,
						Data:        valueData,
					})

					sendChannel <- wrapperData
				}
			}
		}
	}(connection)
}

func GetConnectionMeta(context context.Context) *connectionMeta {
	return context.Value(connContextKey{key: "connId"}).(*connectionMeta)
}
