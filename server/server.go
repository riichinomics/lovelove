package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

type connDetails struct {
	userId string
	connId string
}

type CardLocation int64

const (
	CardLocation_Deck CardLocation = iota
	CardLocation_Table
	CardLocation_RedHand
	CardLocation_WhiteHand
	CardLocation_RedCollection
	CardLocation_WhiteCollection
	CardLocation_Drawn
)

type GameState int64

const (
	GameState_HandCardPlay GameState = iota
	GameState_DeckCardPlay
	GameState_DeclareWin
)

type cardState struct {
	location CardLocation
	order    int
	card     *lovelove.Card
}

type cardId struct {
	hana      lovelove.Hana
	variation lovelove.Variation
}

type playerState struct {
	id       string
	position lovelove.PlayerPosition
}

type gameState struct {
	state        GameState
	id           string
	activePlayer lovelove.PlayerPosition
	oya          lovelove.PlayerPosition
	cards        map[cardId]cardState
	playerState  map[string]*playerState
}

type LoveLoveServer struct {
	connDetails map[string]*connDetails
	games       map[string]*gameState
}

type LoveLoveRpcServer struct {
	lovelove.UnimplementedLoveLoveServer
	server *LoveLoveServer
}

func (LoveLoveRpcServer) SayHello(context context.Context, request *lovelove.HelloRequest) (*lovelove.HelloReply, error) {
	log.Print(request.Name)
	return &lovelove.HelloReply{Message: "Hello " + request.Name}, nil
}

func (server LoveLoveRpcServer) Authenticate(context context.Context, request *lovelove.AuthenticateRequest) (*lovelove.AuthenticateResponse, error) {
	connId := context.Value(connContextKey{key: "connId"}).(string)
	log.Print(connId, request.UserId)

	if connDetails, ok := server.server.connDetails[connId]; ok {
		connDetails.userId = request.UserId
	}

	return &lovelove.AuthenticateResponse{}, nil
}

func cardIdFromCard(card lovelove.Card) cardId {
	return cardId{
		hana:      card.Hana,
		variation: card.Variation,
	}
}

func moveCards(cardMap map[cardId]cardState, cards []*lovelove.Card, location CardLocation) {
	for i, card := range cards {
		cardMap[cardIdFromCard(*card)] = cardState{
			order:    i,
			card:     card,
			location: location,
		}
	}
}

func createGameStateView(gameState gameState, playerPosition lovelove.PlayerPosition) *lovelove.CompleteGameState {
	zones := make(map[CardLocation][]cardState)

	for _, card := range gameState.cards {
		zone, zoneFound := zones[card.location]
		if !zoneFound {
			zone = make([]cardState, 0, 12*4)
		}
		zones[card.location] = append(zone, card)
	}

	completeGameState := &lovelove.CompleteGameState{
		Deck:               0,
		Table:              make([]*lovelove.Card, 0, 12*4),
		Hand:               make([]*lovelove.Card, 0, 12*4),
		Collection:         make([]*lovelove.Card, 0, 12*4),
		OpponentHand:       0,
		OpponentCollection: make([]*lovelove.Card, 0, 12*4),
		Active:             gameState.activePlayer,
		Oya:                gameState.oya,
	}

	for zoneType, zone := range zones {
		sort.Slice(zone, func(i, j int) bool {
			return zone[i].order < zone[j].order
		})

		cards := make([]*lovelove.Card, 0, 12*4)
		for _, card := range zone {
			cards = append(cards, card.card)
		}

		switch zoneType {
		case CardLocation_Deck:
			completeGameState.Deck = int32(len(zone))
		case CardLocation_Table:
			completeGameState.Table = cards
		case CardLocation_RedCollection:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.Collection = cards
			} else {
				completeGameState.OpponentCollection = cards
			}
		case CardLocation_WhiteCollection:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.OpponentCollection = cards
			} else {
				completeGameState.Collection = cards
			}
		case CardLocation_RedHand:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.Hand = cards
			} else {
				completeGameState.OpponentHand = int32(len(zone))
			}
		case CardLocation_WhiteHand:
			if playerPosition == lovelove.PlayerPosition_Red {
				completeGameState.OpponentHand = int32(len(zone))
			} else {
				completeGameState.Hand = cards
			}
		}
	}

	if gameState.activePlayer == playerPosition && gameState.state == GameState_HandCardPlay {
		completeGameState.Action = &lovelove.PlayerAction{
			Type:        lovelove.PlayerActionType_HandCardPlayOpportunity,
			PlayOptions: make(map[int32]*lovelove.PlayOptions),
		}
		for _, tableCard := range completeGameState.Table {
			playOptions := make([]int32, 0)
			for _, handCard := range completeGameState.Hand {
				if handCard.Hana == tableCard.Hana {
					playOptions = append(playOptions, handCard.Id)
				}
			}

			if len(playOptions) > 0 {
				completeGameState.Action.PlayOptions[tableCard.Id] = &lovelove.PlayOptions{
					Options: playOptions,
				}
			}
		}
	}

	return completeGameState
}

func (server LoveLoveRpcServer) ConnectToGame(context context.Context, request *lovelove.ConnectToGameRequest) (*lovelove.ConnectToGameResponse, error) {
	log.Print(request.RoomId)

	game, gameFound := server.server.games[request.RoomId]

	// TODO: deal with missing connection problem?
	connDetails := server.server.connDetails[context.Value(connContextKey{
		key: "connId",
	}).(string)]

	if len(connDetails.userId) == 0 {
		log.Print("Player not identified")
		return &lovelove.ConnectToGameResponse{}, nil
	}

	if !gameFound {
		deck := make([]*lovelove.Card, 12*4)

		for hana := range lovelove.Hana_name {
			for variation := range lovelove.Variation_name {
				deck[hana*4+variation] = &lovelove.Card{
					Id:        hana*4 + variation,
					Hana:      lovelove.Hana(hana),
					Variation: lovelove.Variation(variation),
				}
			}
		}

		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})

		oya := lovelove.PlayerPosition(rand.Intn(2))

		game = &gameState{
			state:        GameState_HandCardPlay,
			id:           request.RoomId,
			activePlayer: oya,
			cards:        make(map[cardId]cardState),
			playerState:  make(map[string]*playerState),
			oya:          oya,
		}

		game.playerState[connDetails.userId] = &playerState{
			id:       connDetails.userId,
			position: lovelove.PlayerPosition(rand.Intn(2)),
		}

		moveCards(game.cards, deck[0:8], CardLocation_Table)
		moveCards(game.cards, deck[8:16], CardLocation_RedHand)
		moveCards(game.cards, deck[16:24], CardLocation_WhiteHand)
		moveCards(game.cards, deck[24:], CardLocation_Deck)

		server.server.games[game.id] = game
	} else {
		_, playerExists := game.playerState[connDetails.userId]
		if !playerExists && len(game.playerState) < 2 {
			newPlayerPosition := lovelove.PlayerPosition_Red
			for _, playerState := range game.playerState {
				if playerState.position == lovelove.PlayerPosition_Red {
					newPlayerPosition = lovelove.PlayerPosition_White
				}
			}

			game.playerState[connDetails.userId] = &playerState{
				id:       connDetails.userId,
				position: newPlayerPosition,
			}
		}
	}

	playerPosition := game.playerState[connDetails.userId].position

	return &lovelove.ConnectToGameResponse{
		Position:  playerPosition,
		GameState: createGameStateView(*game, playerPosition),
	}, nil
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
	log.Print("Write")

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
	connMap  map[string]*connDetails
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
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	log.SetFlags(0)

	server := &WebSocketRpcServer{
		services: make(map[string]*serviceInfo),
		connMap:  make(map[string]*connDetails),
	}

	lovelove.RegisterLoveLoveServer(server, &LoveLoveRpcServer{
		UnimplementedLoveLoveServer: lovelove.UnimplementedLoveLoveServer{},
		server: &LoveLoveServer{
			connDetails: server.connMap,
			games:       make(map[string]*gameState),
		},
	})

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

type connContextKey struct {
	key string
}

func (server *WebSocketRpcServer) Connect(connection *websocket.Conn) {
	go func(conn *websocket.Conn) {
		connId := uuid.New().String()

		server.connMap[connId] = &connDetails{
			connId: connId,
		}

		defer delete(server.connMap, connId)

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
					value, _ := methodInfo.Handler(
						serviceInfo.serviceImpl,
						context.WithValue(context.Background(), connContextKey{
							key: "connId",
						}, connId),
						func(message interface{}) error {
							return proto.Unmarshal(wrapper.Data, message.(proto.Message))
						},
						nil,
					)

					valueData, _ := proto.Marshal(value.(proto.Message))

					wrapperData, _ := proto.Marshal(&lovelove.Wrapper{
						Sequence: wrapper.Sequence,
						Type:     wrapper.Type,
						Data:     valueData,
					})

					conn.WriteMessage(websocket.BinaryMessage, wrapperData)
				}
			}
		}
	}(connection)
}
