package rpc

import (
	"context"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

func (API) Check(_ context.Context, _ *proto.CheckRequest) (*proto.CheckResponse, error) {
	return &proto.CheckResponse{
		Status: proto.CheckResponse_SERVING_STATUS_SERVING,
	}, nil
}
