package engine

import (
	"context"

	lovelove "hanafuda.moe/lovelove/proto"
)

func (server loveLoveRpcServer) PlayDrawnCard(context context.Context, request *lovelove.PlayDrawnCardRequest) (response *lovelove.PlayDrawnCardResponse, rpcError error) {
	response = &lovelove.PlayDrawnCardResponse{
		Status: lovelove.GenericResponseCode_Error,
	}
	rpcError = nil

	return
}
