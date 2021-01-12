package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	kafkaproducerprovider "github.com/ravilushqa/highload/providers/kafka-producer"
	postsGrpc "github.com/ravilushqa/highload/services/posts/api/grpc"
	"github.com/ravilushqa/highload/services/posts/lib/post"
)

const cacheKey = "feed:user_id:%d"

type Api struct {
	logger        *zap.Logger
	redis         *redis.Client
	postManager   *post.Manager
	kafkaProducer *kafkaproducerprovider.KafkaProducer
}

func NewApi(logger *zap.Logger, redis *redis.Client, postManager *post.Manager, kafkaProducer *kafkaproducerprovider.KafkaProducer) *Api {
	return &Api{logger: logger, redis: redis, postManager: postManager, kafkaProducer: kafkaProducer}
}

func (a *Api) Run(ctx context.Context) error {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr) //@todo
	if err != nil {
		return err
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcprometheus.StreamServerInterceptor,
			grpczap.StreamServerInterceptor(a.logger.Named("grpc_stream")),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcprometheus.UnaryServerInterceptor,
			grpczap.UnaryServerInterceptor(a.logger.Named("grpc_unary")),
		)),
	)
	postsGrpc.RegisterPostsServer(s, a)

	reflection.Register(s)

	a.logger.Info("api started..", zap.String("addr", addr))

	defer s.GracefulStop()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return s.Serve(lis)
}

func (a *Api) GetFeed(ctx context.Context, req *postsGrpc.GetFeedRequest) (*postsGrpc.GetFeedResponse, error) {
	list, err := a.redis.LRange(fmt.Sprintf(cacheKey, req.UserId), 0, 1000).Result()
	if err != nil {
		a.logger.Error("failed to get feed from cache", zap.Error(err), zap.Int64("user_id", req.UserId))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	posts := make([]*postsGrpc.Post, 0, len(list))
	for _, jsonPost := range list {
		var p post.Post
		err = json.Unmarshal([]byte(jsonPost), &p)
		if err != nil {
			a.logger.Error("failed unmarshal post", zap.Error(err), zap.Int64("user_id", req.UserId))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		ca, err := ptypes.TimestampProto(p.CreatedAt)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}

		var da *timestamp.Timestamp
		if p.DeletedAt.Valid {
			da, err = ptypes.TimestampProto(p.DeletedAt.Time)
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}
		}

		posts = append(posts, &postsGrpc.Post{
			Id:        int64(p.ID),
			UserId:    int64(p.UserID),
			Text:      p.Text,
			CreatedAt: ca,
			DeletedAt: da,
		})
	}

	return &postsGrpc.GetFeedResponse{
		Posts: posts,
	}, err
}

func (a *Api) GetByUserID(ctx context.Context, req *postsGrpc.GetByUserIDRequest) (*postsGrpc.GetByUserIDResponse, error) {
	posts, err := a.postManager.GetOwnPosts(ctx, int(req.UserId))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := make([]*postsGrpc.Post, 0, len(posts))

	for _, p := range posts {
		ca, err := ptypes.TimestampProto(p.CreatedAt)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}

		var da *timestamp.Timestamp
		if p.DeletedAt.Valid {
			da, err = ptypes.TimestampProto(p.DeletedAt.Time)
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}
		}

		res = append(res, &postsGrpc.Post{
			Id:        int64(p.ID),
			UserId:    int64(p.UserID),
			Text:      p.Text,
			CreatedAt: ca,
			DeletedAt: da,
		})
	}

	return &postsGrpc.GetByUserIDResponse{Posts: res}, err
}

func (a *Api) Store(ctx context.Context, req *postsGrpc.StoreRequest) (*postsGrpc.StoreResponse, error) {
	p := &post.Post{
		UserID: int(req.UserId),
		Text:   req.Text,
	}
	p, err := a.postManager.Insert(ctx, p)
	if err != nil {
		a.logger.Error("failed to insert post", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	ca, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	var da *timestamp.Timestamp
	if p.DeletedAt.Valid {
		da, err = ptypes.TimestampProto(p.DeletedAt.Time)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	message, err := json.Marshal(p)
	if err != nil {
		a.logger.Error("failed to marshal post", zap.Error(err))
	} else if err = a.kafkaProducer.SendMessage(message, nil); err != nil {
		a.logger.Error("failed to send message to kafka", zap.Error(err))
	}

	return &postsGrpc.StoreResponse{Post: &postsGrpc.Post{
		Id:        int64(p.ID),
		UserId:    int64(p.UserID),
		Text:      p.Text,
		CreatedAt: ca,
		DeletedAt: da,
	}}, nil
}
