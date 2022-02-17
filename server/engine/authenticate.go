package engine

import (
	"context"
	"log"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) Authenticate(context context.Context, request *lovelove.AuthenticateRequest) (*lovelove.AuthenticateResponse, error) {
	connMeta := GetConnectionMeta(context)
	log.Print("User Authenticated")
	connMeta.userId = request.UserId
	return &lovelove.AuthenticateResponse{}, nil
}
