package main

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/axengine/go-saga"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatsGrpc "github.com/ravilushqa/highload/services/chats/api/grpc"
	"github.com/ravilushqa/highload/services/chats/lib/chat"
	chatuser "github.com/ravilushqa/highload/services/chats/lib/chat_user"
	"github.com/ravilushqa/highload/services/chats/lib/message"
	countersGrpc "github.com/ravilushqa/highload/services/counters/api/grpc"
)

type Api struct {
	chatsGrpc.UnimplementedChatsServer
	chatUserManager *chatuser.Manager
	chatManager     *chat.Manager
	messageManager  *message.Manager
	logger          *zap.Logger
	saga            *saga.ExecutionCoordinator
	countersClient  countersGrpc.CountersClient
}

func NewApi(chatUserManager *chatuser.Manager, chatManager *chat.Manager, messageManager *message.Manager, logger *zap.Logger, saga *saga.ExecutionCoordinator, countersClient countersGrpc.CountersClient) *Api {
	a := &Api{chatUserManager: chatUserManager, chatManager: chatManager, messageManager: messageManager, logger: logger, saga: saga, countersClient: countersClient}
	a.initSagas()

	return a
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
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcprometheus.UnaryServerInterceptor,
			grpczap.UnaryServerInterceptor(a.logger.Named("grpc_unary")),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	chatsGrpc.RegisterChatsServer(s, a)

	reflection.Register(s)

	a.logger.Info("api started..", zap.String("addr", addr))

	defer s.GracefulStop()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return s.Serve(lis)
}

func (a *Api) GetUserChats(ctx context.Context, req *chatsGrpc.GetUserChatsRequest) (*chatsGrpc.GetUserChatsResponse, error) {
	chatIDs, err := a.chatUserManager.GetUserChats(ctx, req.UserId)
	if err != nil {
		a.logger.Error("failed get chats", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	chatIDsRes := make([]string, 0, len(chatIDs))
	for _, v := range chatIDs {
		chatIDsRes = append(chatIDsRes, v)
	}

	return &chatsGrpc.GetUserChatsResponse{
		ChatIds: chatIDsRes,
	}, nil
}

func (a *Api) GetChatMessages(ctx context.Context, req *chatsGrpc.GetChatMessagesRequest) (*chatsGrpc.GetChatMessagesResponse, error) {
	messages, err := a.messageManager.GetChatMessages(ctx, []string{req.ChatId})
	if err != nil {
		a.logger.Error("failed get messages", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	_, err = a.countersClient.FlushChatCounter(ctx, &countersGrpc.FlushChatCounterRequest{
		UserId: req.UserId,
		ChatId: req.ChatId,
	})

	if err != nil {
		a.logger.Error("failed flush chat counter", zap.Error(err))
	}

	resMessages, err := a.messagesToProto(messages)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &chatsGrpc.GetChatMessagesResponse{
		Messages: resMessages,
	}, nil
}

func (a *Api) StoreMessage(ctx context.Context, req *chatsGrpc.StoreMessageRequest) (*empty.Empty, error) {
	userIDs, err := a.chatUserManager.GetChatMembers(ctx, req.ChatId)
	if err != nil {
		a.logger.Error("failed get chat members", zap.Error(err))
		return &empty.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}

	receivers := make([]string, 0, len(userIDs)-1)

	for i, userID := range userIDs {
		if userID == req.UserId {
			receivers = append(userIDs[:i], userIDs[i+1:]...)
			break
		}
	}

	err = a.saga.StartSaga(ctx, strconv.Itoa(time.Now().Nanosecond())).
		ExecSub("store_message", req.UserId, req.ChatId, req.Text).
		ExecSub("update_counter", req.ChatId, receivers).
		EndSaga()

	if err != nil {
		return &empty.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}

	return &empty.Empty{}, nil
}

func (a *Api) FindOrCreateChat(ctx context.Context, req *chatsGrpc.FindOrCreateChatRequest) (*chatsGrpc.FindOrCreateChatResponse, error) {
	chatID, err := a.chatUserManager.GetUsersDialogChat(ctx, req.UserId_1, req.UserId_2)
	if err != nil {
		a.logger.Error("failed get users dialog chat", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if chatID == "" {
		chatID, err = a.chatManager.Insert(ctx, &chat.Chat{
			Type: "dialog",
		})
		if err != nil {
			a.logger.Error("failed create chat", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		err = a.chatUserManager.Insert(ctx, &chatuser.ChatUser{
			UserID: req.UserId_1,
			ChatID: chatID,
		})

		if err != nil {
			a.logger.Error("failed create chat user", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		err = a.chatUserManager.Insert(ctx, &chatuser.ChatUser{
			UserID: req.UserId_2,
			ChatID: chatID,
		})

		if err != nil {
			a.logger.Error("failed create chat user", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &chatsGrpc.FindOrCreateChatResponse{
		ChatId: chatID,
	}, err
}

func (a *Api) messagesToProto(messages []message.Message) ([]*chatsGrpc.Message, error) {
	resMessages := make([]*chatsGrpc.Message, 0, len(messages))
	for _, v := range messages {
		ca, err := ptypes.TimestampProto(v.CreatedAt)
		if err != nil {
			return nil, err
		}
		var ua *timestamp.Timestamp
		if v.UpdatedAt != nil {
			ua = timestamppb.New(*v.UpdatedAt)
			if err != nil {
				return nil, err
			}
		}
		var da *timestamp.Timestamp
		if v.DeletedAt != nil {
			da = timestamppb.New(*v.DeletedAt)
			if err != nil {
				return nil, err
			}
		}
		resMessages = append(resMessages, &chatsGrpc.Message{
			Uuid:      v.UUID,
			UserId:    v.UserID,
			ChatId:    v.ChatID,
			Text:      v.Text,
			CreatedAt: ca,
			UpdatedAt: ua,
			DeletedAt: da,
		})
	}
	return resMessages, nil
}

func (a *Api) initSagas() {
	a.saga.AddSubTxDef(
		"store_message",
		func(ctx context.Context, userID, chatID string, text string) error {
			_, err := a.messageManager.Insert(ctx, &message.Message{
				UserID: userID,
				ChatID: chatID,
				Text:   text,
			})

			return err
		},
		func(ctx context.Context, userID, chatID string, text string) error {
			return a.messageManager.HardDeleteLastMessage(ctx, chatID, userID, text)
		})

	a.saga.AddSubTxDef(
		"update_counter",
		func(ctx context.Context, chatID string, receivers []string) error {
			_, err := a.countersClient.IncrementUnreadMessageCounter(ctx, &countersGrpc.IncrementUnreadMessageCounterRequest{
				UserIds: receivers,
				ChatId:  chatID,
			})

			return err
		},
		func(ctx context.Context, chatID string, receivers []string) error {
			_, err := a.countersClient.DecrementUnreadMessageCounter(ctx, &countersGrpc.DecrementUnreadMessageCounterRequest{
				UserIds: receivers,
				ChatId:  chatID,
			})

			return err
		})
}
