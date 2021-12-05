package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
	"hanafuda.moe/lovelove/rpc"
)

func (server loveLoveRpcServer) Authenticate(context context.Context, request *lovelove.AuthenticateRequest) (*lovelove.AuthenticateResponse, error) {
	rpcConnMeta := rpc.GetConnectionMeta(context)
	log.Print(rpcConnMeta.ConnId, request.UserId)

	connMeta, ok := server.connectionMeta[rpcConnMeta.ConnId]

	if !ok {
		connMeta = &connectionMeta{}
		server.connectionMeta[rpcConnMeta.ConnId] = connMeta
		rpcConnMeta.Closed.DoOnCompleted(func() {
			delete(server.connectionMeta, rpcConnMeta.ConnId)
		})
	}
	connMeta.userId = request.UserId

	return &lovelove.AuthenticateResponse{}, nil
}
