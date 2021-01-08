package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	lovelove "hanafuda.moe/lovelove/proto"
)

var addr = flag.String("addr", "localhost:6482", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		log.Print("Reading")
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	log.Print("Closed")
}

type LoveLoveRpcServer struct {
	lovelove.UnimplementedLoveLoveServer
}

func (LoveLoveRpcServer) SayHello(context context.Context, request *lovelove.HelloRequest) (*lovelove.HelloReply, error) {
	log.Print(request.Name)
	return &lovelove.HelloReply{Message: "Hello " + request.Name}, nil
}

type WebSocketListener struct {
	websocketOpened chan *websocket.Conn
}

type WebSocketConn struct {
	*websocket.Conn
	reader io.Reader
}

func NewWebSocketListener() *WebSocketListener {
	return &WebSocketListener{
		websocketOpened: make(chan *websocket.Conn),
	}
}

func (conn *WebSocketConn) Read(buffer []byte) (int, error) {
	return conn.reader.Read(buffer)
}

func (conn *WebSocketConn) Write(buffer []byte) (int, error) {
	err := conn.WriteMessage(websocket.BinaryMessage, buffer)
	if err != nil {
		return 0, err
	}
	return len(buffer), nil
}

func (conn *WebSocketConn) SetDeadline(t time.Time) error {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return conn.SetWriteDeadline(t)
}

func (listener WebSocketListener) Accept() (net.Conn, error) {
	conn := <-listener.websocketOpened
	log.Print("Accepted")
	return &WebSocketConn{
		Conn:   conn,
		reader: websocket.JoinMessages(conn, ""),
	}, nil
}

func (WebSocketListener) Close() error {
	log.Print("Close")
	return nil
}

func (WebSocketListener) Addr() net.Addr {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			return addr
		}
	}
	return nil
}

type serviceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*grpc.MethodDesc
	streams     map[string]*grpc.StreamDesc
	mdata       interface{}
}

type WebSocketRpcServer struct {
	// lis      map[net.Listener]bool
	// conns    map[transport.ServerTransport]bool
	serve    bool
	services map[string]*serviceInfo
}

func (server *WebSocketRpcServer) RegisterService(serviceDesc *grpc.ServiceDesc, ss interface{}) {
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

func main() {
	flag.Parse()
	log.SetFlags(0)

	server := &WebSocketRpcServer{
		services: make(map[string]*serviceInfo),
	}
	lovelove.RegisterLoveLoveServer(server, &LoveLoveRpcServer{})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		log.Print("Connected")
		server.Connect(c)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("starting")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (server *WebSocketRpcServer) Connect(connection *websocket.Conn) {
	go func(conn *websocket.Conn) {
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
			log.Print(wrapper.Sequence, wrapper.Type, wrapper.Data)

			method := wrapper.Type
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
					methodInfo.Handler(
						serviceInfo.serviceImpl,
						context.TODO(),
						func(message interface{}) error {
							return proto.Unmarshal(wrapper.Data, message.(proto.Message))
						},
						nil,
					)
					return
				}
			}
		}
	}(connection)
}
