package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	muhjong "riichi.moe/muhjong/proto"
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
}

type MuhjongRpcServer struct {
	muhjong.UnimplementedMuhjongServer
}

func (MuhjongRpcServer) SayHello(context context.Context, request *muhjong.HelloRequest) (*muhjong.HelloReply, error) {
	log.Print(request.Name)
	return &muhjong.HelloReply{Message: "Hello " + request.Name}, nil
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
	return &WebSocketConn{
		Conn:   conn,
		reader: websocket.JoinMessages(conn, ""),
	}, nil
}

func (WebSocketListener) Close() error {
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

func main() {
	flag.Parse()
	log.SetFlags(0)

	listener := NewWebSocketListener()

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		listener.websocketOpened <- c
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Fatal(http.ListenAndServe(*addr, nil))

	server := grpc.NewServer()
	muhjong.RegisterMuhjongServer(server, &MuhjongRpcServer{})
	server.Serve(listener)
}
