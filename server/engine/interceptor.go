package engine

import (
	"context"
	"errors"
	"log"
	"sync"

	"google.golang.org/grpc"
	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

type connectionMeta struct {
	userId            string
	roomId            string
	roomChangedNotify func()
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

type handlerResponse struct {
	value interface{}
	err   error
}

func (interceptor *loveLoveRpcInterceptor) Interceptor(rpcContext context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	rpcConnectionMeta := rpc.GetConnectionMeta(rpcContext)

	connMeta, hasConnMeta := interceptor.connectionMeta[rpcConnectionMeta.ConnId]
	if !hasConnMeta {
		connMeta = &connectionMeta{}
		interceptor.connectionMeta[rpcConnectionMeta.ConnId] = connMeta
		rpcConnectionMeta.Closed.DoOnCompleted(func() {
			delete(interceptor.connectionMeta, rpcConnectionMeta.ConnId)
		})
	}

	rpcContext = withConnectionMeta(rpcContext, connMeta)

	var roomId string

	connectToGameRequest, isConnectToGameRequest := req.(*lovelove.ConnectToGameRequest)
	if isConnectToGameRequest {
		log.Print("room request", connectToGameRequest.RoomId)
		roomId = connectToGameRequest.RoomId
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
				id:               roomId,
				players:          make(map[string]*playerMeta),
				requestQueue:     make(chan func()),
				cleanupRequested: make(chan context.Context),
			}

			interceptor.games[roomId] = game

			go func(game *gameContext) {
				for {
					select {
					case request := <-game.requestQueue:
						request()
					case cleanupContext := <-game.cleanupRequested:
						log.Print("Game cleaned up: ", game.id)
						interceptor.gamesMutex.Lock()
						game.activityState.mutex.Lock()

						if !errors.Is(cleanupContext.Err(), context.Canceled) {
							delete(interceptor.games, game.id)
							game.activityState.cancelCleanup()
							game.activityState.mutex.Unlock()
							interceptor.gamesMutex.Unlock()
							close(game.requestQueue)
							return
						}

						game.activityState.cleanupCancelation = nil
						game.activityState.mutex.Unlock()
						interceptor.gamesMutex.Unlock()
					}
				}
			}(game)
		}

		if isConnectToGameRequest {
			game.StartRequest()
		}

		interceptor.gamesMutex.Unlock()

		rpcContext = withGameContext(rpcContext, game)

		responseChan := make(chan *handlerResponse)
		game.requestQueue <- func() {
			value, err := handler(rpcContext, req)
			if isConnectToGameRequest {
				game.EndRequest()
			}
			responseChan <- &handlerResponse{
				value,
				err,
			}
		}
		response := <-responseChan
		return response.value, response.err
	}

	return handler(rpcContext, req)
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
