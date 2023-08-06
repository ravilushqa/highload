package main

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/users/lib/user"
)

type Api struct {
	usersGrpc.UnimplementedUsersServer
	userManager *user.Manager
}

func NewApi(userManager *user.Manager) *Api {
	return &Api{userManager: userManager}
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
	userID, err := primitive.ObjectIDFromHex(req.RequesterUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	subscriberID, err := primitive.ObjectIDFromHex(req.AddedUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := a.userManager.Subscribe(ctx, userID, subscriberID); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return new(empty.Empty), nil
}

func (a *Api) ApproveFriendRequest(ctx context.Context, req *usersGrpc.ApproveFriendRequestRequest) (*empty.Empty, error) {
	userID, err := primitive.ObjectIDFromHex(req.ApproverUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	subscriberID, err := primitive.ObjectIDFromHex(req.RequesterUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := a.userManager.Subscribe(ctx, userID, subscriberID); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return new(empty.Empty), nil
}

func (a *Api) GetById(ctx context.Context, req *usersGrpc.GetByIdRequest) (*usersGrpc.GetByIdResponse, error) {
	u, err := a.userManager.GetByID(ctx, req.UserId)
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
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	friendIds, err := a.userManager.GetFriends(ctx, userID)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	// @TODO
	friendIdsStr := make([]string, 0, len(friendIds))
	for _, v := range friendIds {
		friendIdsStr = append(friendIdsStr, v.ID.Hex())
	}
	return &usersGrpc.GetFriendsIdsResponse{UserIds: friendIdsStr}, err
}

func (a *Api) GetListByIds(ctx context.Context, req *usersGrpc.GetListByIdsRequest) (*usersGrpc.GetListByIdsResponse, error) {
	friends, err := a.userManager.GetListByIds(ctx, req.UserIds)
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
	fromUserID, err := primitive.ObjectIDFromHex(req.FromUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	toUserID, err := primitive.ObjectIDFromHex(req.ToUserId)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	relation, err := a.userManager.GetRelations(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.GetRelationResponse{
		Relation: usersGrpc.UserRelation(usersGrpc.UserRelation_value[cases.Title(language.Und).String(string(relation))]),
	}, err
}

func (a *Api) GetByEmail(ctx context.Context, req *usersGrpc.GetByEmailRequest) (*usersGrpc.GetByEmailResponse, error) {
	u, err := a.userManager.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if u == nil {
		return nil, status.New(codes.NotFound, "user not found").Err()
	}

	res, err := a.user2proto(u)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.GetByEmailResponse{User: res}, err
}

func (a *Api) Store(ctx context.Context, req *usersGrpc.StoreRequest) (*usersGrpc.StoreResponse, error) {
	userID, err := a.userManager.Store(ctx, &user.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  req.Birthday.AsTime(),
		Interests: req.Interests,
		Sex:       user.Sex(strings.ToLower(req.Sex.String())),
		City:      req.City,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &usersGrpc.StoreResponse{Id: userID.Hex()}, nil
}

func (a *Api) user2proto(u *user.User) (*usersGrpc.User, error) {
	ca := timestamppb.New(u.CreatedAt)
	bd := timestamppb.New(u.Birthday)

	var da *timestamp.Timestamp
	if u.DeletedAt != nil {
		da = timestamppb.New(*u.DeletedAt)
	}
	return &usersGrpc.User{
		Id:        u.ID.Hex(),
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birthday:  bd,
		Interests: u.Interests,
		Sex:       usersGrpc.Sex(usersGrpc.Sex_value[cases.Title(language.Und).String(string(u.Sex))]),
		City:      u.City,
		CreatedAt: ca,
		DeletedAt: da,
	}, nil
}
