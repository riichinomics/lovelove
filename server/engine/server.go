package engine

import lovelove "hanafuda.moe/lovelove/proto"

type loveLoveRpcServer struct {
	lovelove.UnimplementedLoveLoveServer
}

func NewLoveLoveRpcServer() *loveLoveRpcServer {
	return &loveLoveRpcServer{
		UnimplementedLoveLoveServer: lovelove.UnimplementedLoveLoveServer{},
	}
}
