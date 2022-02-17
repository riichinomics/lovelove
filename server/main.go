package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	engine "hanafuda.moe/lovelove/engine"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// TODO: no idea what this does
	flag.Parse()
	log.SetFlags(0)

	// TODO: separate server from interceptor
	interceptor := engine.NewLoveLoveRpcInterceptor()
	loveloveRpcServer := engine.NewLoveLoveRpcServer()
	loveloveRpcServer.LoadTestGames()
	websocketRpcServer := rpc.NewWebSocketRpcServer(rpc.UnaryInterceptor(interceptor.Interceptor))

	lovelove.RegisterLoveLoveServer(websocketRpcServer, loveloveRpcServer)

	addr := flag.String("addr", "0.0.0.0:6482", "http service address")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(request *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		log.Print("New Websocket connection from ", r.RemoteAddr)
		websocketRpcServer.HandleConnection(c)
	})

	http.HandleFunc("/proto/lovelove.proto", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "./proto/lovelove.proto")
	})

	log.Print("starting")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
