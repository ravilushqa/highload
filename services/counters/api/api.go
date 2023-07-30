package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	countersGrpc "github.com/ravilushqa/highload/services/counters/api/grpc"
)

const unreadmessageCacheKey = "unread_message_%s" // userID

type Api struct {
	countersGrpc.UnimplementedCountersServer
	logger *zap.Logger
	redis  *redis.Client
}

func NewApi(logger *zap.Logger, redis *redis.Client) *Api {
	return &Api{logger: logger, redis: redis}
}

func (a *Api) IncrementUnreadMessageCounter(ctx context.Context, req *countersGrpc.IncrementUnreadMessageCounterRequest) (*empty.Empty, error) {
	// err := a.redis.ZAdd("chat_" + strconv.Itoa(int(userID)), redis.Z{Score: 1, Member: chatID}).Err() @todo think about it
	for _, userID := range req.UserIds {
		err := a.redis.HIncrBy(fmt.Sprintf(unreadmessageCacheKey, userID), req.ChatId, 1).Err()
		if err != nil {
			return new(empty.Empty), status.New(codes.Internal, err.Error()).Err()
		}
	}

	return new(empty.Empty), nil
}

func (a *Api) DecrementUnreadMessageCounter(ctx context.Context, req *countersGrpc.DecrementUnreadMessageCounterRequest) (*empty.Empty, error) {
	for _, userID := range req.UserIds {
		res, err := a.redis.HGet(fmt.Sprintf(unreadmessageCacheKey, userID), req.ChatId).Result()
		if err != nil {
			return new(empty.Empty), status.New(codes.Internal, err.Error()).Err()
		}
		if res == "nil" {
			a.logger.Info("decrement on empty hash")
			return new(empty.Empty), nil
		}
		v, err := strconv.Atoi(res)
		if err != nil {
			return new(empty.Empty), status.New(codes.Internal, err.Error()).Err()
		}
		if v <= 0 {
			a.logger.Info("decrement on >= 0 value")
			err := a.redis.HDel(fmt.Sprintf(unreadmessageCacheKey, userID), req.ChatId).Err()
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}

			return new(empty.Empty), nil
		}

		err = a.redis.HIncrBy(fmt.Sprintf(unreadmessageCacheKey, userID), req.ChatId, -1).Err()

		if err != nil {
			return new(empty.Empty), status.New(codes.Internal, err.Error()).Err()
		}
	}
	return new(empty.Empty), nil
}

func (a *Api) UnreadChatsCount(ctx context.Context, req *countersGrpc.UnreadChatsCountRequest) (*countersGrpc.UnreadChatsCountResponse, error) {
	hashMap, err := a.redis.HGetAll(fmt.Sprintf(unreadmessageCacheKey, req.UserId)).Result()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := make(map[string]int64, len(hashMap))
	for chatID, v := range hashMap {
		count, err := strconv.Atoi(v)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		res[chatID] = int64(count)
	}

	return &countersGrpc.UnreadChatsCountResponse{
		ChatsUnreadMessages: res,
	}, err
}

func (a *Api) FlushChatCounter(ctx context.Context, req *countersGrpc.FlushChatCounterRequest) (*empty.Empty, error) {
	err := a.redis.HDel(fmt.Sprintf(unreadmessageCacheKey, req.UserId), req.ChatId).Err()
	if err != nil {
		return new(empty.Empty), status.New(codes.Internal, err.Error()).Err()
	}

	return new(empty.Empty), nil
}
