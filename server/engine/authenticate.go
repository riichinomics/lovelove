package engine

import (
	"context"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) Authenticate(context context.Context, request *lovelove.AuthenticateRequest) (*lovelove.AuthenticateResponse, error) {
	connMeta := GetConnectionMeta(context)
	connMeta.userId = request.UserId
	return &lovelove.AuthenticateResponse{}, nil
}
