package engine

import lovelove "hanafuda.moe/lovelove/proto"

type connectionMeta struct {
	userId string
}

type userMeta struct {
	roomId string
}

type loveLoveRpcServer struct {
	lovelove.UnimplementedLoveLoveServer
	games          map[string]*gameState
	connectionMeta map[string]*connectionMeta
	userMeta       map[string]*userMeta
}

func NewLoveLoveRpcServer() *loveLoveRpcServer {
	return &loveLoveRpcServer{
		UnimplementedLoveLoveServer: lovelove.UnimplementedLoveLoveServer{},
		games:                       make(map[string]*gameState),
		connectionMeta:              make(map[string]*connectionMeta),
		userMeta:                    make(map[string]*userMeta),
	}
}
