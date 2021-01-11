package main

import (
	"context"
	"net"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/users/lib/friend"
	"github.com/ravilushqa/highload/services/users/lib/user"
)

type Api struct {
	logger        *zap.Logger
	userManager   *user.Manager
	friendManager *friend.Manager
}

func NewApi(logger *zap.Logger, userManager *user.Manager, friendManager *friend.Manager) *Api {
	return &Api{logger: logger, userManager: userManager, friendManager: friendManager}
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
	usersGrpc.RegisterUsersServer(s, a)

	reflection.Register(s)

	a.logger.Info("api started..", zap.String("addr", addr))

	defer s.GracefulStop()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return s.Serve(lis)
}

func (a *Api) GetAll(ctx context.Context, req *usersGrpc.GetUsersRequest) (*usersGrpc.GetUsersResponse, error) {
	users, err := a.userManager.GetAll(ctx, req.Filter)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := make([]*usersGrpc.User, 0, len(users))

	for _, v := range users {
		userProto, err := a.user2proto(&v)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}

		res = append(res, userProto)
	}

	return &usersGrpc.GetUsersResponse{Users: res}, err
}

func (a *Api) FriendRequest(ctx context.Context, req *usersGrpc.FriendRequestRequest) (*empty.Empty, error) {
	if err := a.friendManager.FriendRequest(ctx, int(req.RequesterUserId), int(req.AddedUserId)); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return new(empty.Empty), nil
}

func (a *Api) ApproveFriendRequest(ctx context.Context, req *usersGrpc.ApproveFriendRequestRequest) (*empty.Empty, error) {
	if err := a.friendManager.ApproveFriendRequest(ctx, int(req.ApproverUserId), int(req.RequesterUserId)); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return new(empty.Empty), nil
}

func (a *Api) GetById(ctx context.Context, req *usersGrpc.GetByIdRequest) (*usersGrpc.GetByIdResponse, error) {
	u, err := a.userManager.GetByID(ctx, int(req.UserId))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	userProto, err := a.user2proto(u)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &usersGrpc.GetByIdResponse{User: userProto}, err
}

func (a *Api) GetFriendsIds(ctx context.Context, req *usersGrpc.GetFriendsIdsRequest) (*usersGrpc.GetFriendsIdsResponse, error) {
	friendIds, err := a.friendManager.GetFriends(ctx, int(req.UserId))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res := make([]int64, 0, len(friendIds))
	for _, v := range friendIds {
		res = append(res, int64(v))
	}

	return &usersGrpc.GetFriendsIdsResponse{UserIds: res}, err
}

func (a *Api) GetListByIds(ctx context.Context, req *usersGrpc.GetListByIdsRequest) (*usersGrpc.GetListByIdsResponse, error) {
	ids := make([]int, 0, len(req.UserIds))
	for _, v := range req.UserIds {
		ids = append(ids, int(v))
	}
	friends, err := a.userManager.GetListByIds(ctx, ids)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := make([]*usersGrpc.User, 0, len(friends))
	for _, v := range friends {
		userProto, err := a.user2proto(&v)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		res = append(res, userProto)
	}

	return &usersGrpc.GetListByIdsResponse{Users: res}, err
}

func (a *Api) GetRelation(ctx context.Context, req *usersGrpc.GetRelationRequest) (*usersGrpc.GetRelationResponse, error) {
	relation, err := a.friendManager.GetRelation(ctx, int(req.FromUserId), int(req.ToUserId))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.GetRelationResponse{
		Relation: usersGrpc.UserRelation(usersGrpc.UserRelation_value[strings.Title(string(relation))]),
	}, err
}

func (a *Api) GetByEmail(ctx context.Context, req *usersGrpc.GetByEmailRequest) (*usersGrpc.GetByEmailResponse, error) {
	u, err := a.userManager.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res, err := a.user2proto(u)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.GetByEmailResponse{User: res}, err
}

func (a *Api) Store(ctx context.Context, req *usersGrpc.StoreRequest) (*usersGrpc.StoreResponse, error) {
	bd, err := ptypes.Timestamp(req.Birthday)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	userID, err := a.userManager.Store(ctx, &user.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  bd,
		Interests: req.Interests,
		Sex:       user.Sex(strings.ToLower(req.Sex.String())),
		City:      req.City,
	})
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.StoreResponse{Id: int64(userID)}, nil
}

func (a *Api) user2proto(u *user.User) (*usersGrpc.User, error) {
	ca, err := ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return nil, err
	}

	bd, err := ptypes.TimestampProto(u.Birthday)
	if err != nil {
		return nil, err
	}

	var da *timestamp.Timestamp
	if u.DeletedAt.Valid {
		da, err = ptypes.TimestampProto(u.DeletedAt.Time)
		if err != nil {
			return nil, err
		}
	}
	return &usersGrpc.User{
		Id:        int64(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birthday:  bd,
		Interests: u.Interests,
		Sex:       usersGrpc.Sex(usersGrpc.Sex_value[strings.Title(string(u.Sex))]),
		City:      u.City,
		CreatedAt: ca,
		DeletedAt: da,
	}, nil
}
