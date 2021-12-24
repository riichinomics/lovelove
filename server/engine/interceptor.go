package engine

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

type connectionMeta struct {
	userId string
	roomId string
}

type loveLoveRpcInterceptor struct {
	gamesMutex     sync.Mutex
	games          map[string]*gameContext
	connectionMeta map[string]*connectionMeta
}

func NewLoveLoveRpcInterceptor() *loveLoveRpcInterceptor {
	return &loveLoveRpcInterceptor{
		gamesMutex:     sync.Mutex{},
		games:          make(map[string]*gameContext),
		connectionMeta: make(map[string]*connectionMeta),
	}
}

type loveloveRpcServerContextKey struct {
	key string
}

type gameContext struct {
	id           string
	GameState    *gameState
	listeners    map[string][]chan proto.Message
	requestQueue chan func()
}

type handlerResponse struct {
	value interface{}
	err   error
}

func (interceptor *loveLoveRpcInterceptor) Interceptor(context context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	rpcConnectionMeta := rpc.GetConnectionMeta(context)

	connMeta, hasConnMeta := interceptor.connectionMeta[rpcConnectionMeta.ConnId]
	if !hasConnMeta {
		connMeta = &connectionMeta{}
		interceptor.connectionMeta[rpcConnectionMeta.ConnId] = connMeta
		rpcConnectionMeta.Closed.DoOnCompleted(func() {
			delete(interceptor.connectionMeta, rpcConnectionMeta.ConnId)
		})
	}

	context = withConnectionMeta(context, connMeta)

	var roomId string

	if request, ok := req.(*lovelove.ConnectToGameRequest); ok {
		log.Print("room request", request.RoomId)
		roomId = request.RoomId
	}

	//TODO: leave old room

	if len(roomId) == 0 {
		roomId = connMeta.roomId
	}

	if len(roomId) != 0 {
		interceptor.gamesMutex.Lock()
		game, gameExists := interceptor.games[roomId]
		if !gameExists {
			game = &gameContext{
				id:           roomId,
				requestQueue: make(chan func()),
				listeners:    make(map[string][]chan proto.Message),
			}

			interceptor.games[roomId] = game

			go func(game *gameContext) {
				for request := range game.requestQueue {
					request()
				}
			}(game)
		}
		interceptor.gamesMutex.Unlock()

		context = withGameContext(context, game)

		responseChan := make(chan *handlerResponse)
		game.requestQueue <- func() {
			value, err := handler(context, req)
			responseChan <- &handlerResponse{
				value,
				err,
			}
		}
		response := <-responseChan
		return response.value, response.err
	}

	return handler(context, req)
}

func withGameContext(ctx context.Context, gameContext *gameContext) context.Context {
	return context.WithValue(ctx, loveloveRpcServerContextKey{
		key: "gameContext",
	}, gameContext)
}

func GetGameContext(context context.Context) *gameContext {
	value, ok := context.Value(loveloveRpcServerContextKey{
		key: "gameContext",
	}).(*gameContext)

	if !ok {
		return nil
	}
	return value
}

func withConnectionMeta(ctx context.Context, connectionMeta *connectionMeta) context.Context {
	return context.WithValue(ctx, loveloveRpcServerContextKey{
		key: "connectionMeta",
	}, connectionMeta)
}

func GetConnectionMeta(context context.Context) *connectionMeta {
	value, ok := context.Value(loveloveRpcServerContextKey{
		key: "connectionMeta",
	}).(*connectionMeta)

	if !ok {
		return nil
	}
	return value
}

func (context *gameContext) BroadcastUpdates(gameUpdates map[string][]*lovelove.GameStateUpdatePart) {
	for playerId, updates := range gameUpdates {
		listeners, ok := context.listeners[playerId]
		if !ok {
			continue
		}

		for _, listener := range listeners {
			listener <- &lovelove.GameStateUpdate{
				Updates: updates,
			}
		}
	}
}
