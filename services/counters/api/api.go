package main

import (
	"context"

	"github.com/go-redis/redis"
	"go.uber.org/zap"

	countersGrpc "github.com/ravilushqa/highload/services/counters/api/grpc"
)

type Api struct {
	logger *zap.Logger
	redis  *redis.Client
}

func NewApi(logger *zap.Logger, redis *redis.Client) *Api {
	return &Api{logger: logger, redis: redis}
}

func (a *Api) UnreadMessages(ctx context.Context, req *countersGrpc.UnreadMessagesRequest) (*countersGrpc.UnreadMessagesResponse, error) {
	return &countersGrpc.UnreadMessagesResponse{Count: 9999}, nil
}
