package grpc

import (
	"context"
	"net"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	//grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	//grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/ravilushqa/highload/services/chats/repository/chat"
	chatuser "github.com/ravilushqa/highload/services/chats/repository/chat_user"
	"github.com/ravilushqa/highload/services/chats/repository/message"
)

type Api struct {
	chatUserManager *chatuser.Manager
	chatManager     *chat.Manager
	messageManager  *message.Manager
	logger          *zap.Logger
}

func NewApi(chatUserManager *chatuser.Manager, chatManager *chat.Manager, messageManager *message.Manager, logger *zap.Logger) *Api {
	return &Api{chatUserManager: chatUserManager, chatManager: chatManager, messageManager: messageManager, logger: logger}
}

func (a *Api) Run(ctx context.Context) error {
	addr := "50051"
	lis, err := net.Listen("tcp", addr) //@todo
	if err != nil {
		return err
	}

	s := grpc.NewServer(
	//grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	//	grpc_prometheus.StreamServerInterceptor,
	//	grpc_zap.StreamServerInterceptor(a.logger.Named("grpc_stream")),
	//)),
	//grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	//	grpc_prometheus.UnaryServerInterceptor,
	//	grpc_zap.UnaryServerInterceptor(a.logger.Named("grpc_unary")),
	//)),
	)
	RegisterChatsServer(s, a)

	reflection.Register(s)

	a.logger.Info("api started..", zap.String("addr", addr))

	defer s.GracefulStop()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return s.Serve(lis)
}

func (a *Api) GetUserChats(ctx context.Context, req *GetUserChatsRequest) (*GetUserChatsResponse, error) {
	chatIDs, err := a.chatUserManager.GetUserChats(ctx, int(req.UserId))
	if err != nil {
		a.logger.Error("failed get chats", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	chatIDsRes := make([]int64, 0, len(chatIDs))
	for _, v := range chatIDs {
		chatIDsRes = append(chatIDsRes, int64(v))
	}

	return &GetUserChatsResponse{
		ChatIds: chatIDsRes,
	}, nil
}

func (a *Api) GetChatMessages(ctx context.Context, req *GetChatMessagesRequest) (*GetChatMessagesResponse, error) {
	messages, err := a.messageManager.GetChatMessages(ctx, []int{int(req.ChatId)})
	if err != nil {
		a.logger.Error("failed get messages", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	resMessages, err := a.messagesToProto(messages)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &GetChatMessagesResponse{
		Messages: resMessages,
	}, nil
}

func (a *Api) StoreMessage(ctx context.Context, req *StoreMessageRequest) (*empty.Empty, error) {
	err := a.messageManager.Insert(ctx, &message.Message{
		UserID: int(req.UserId),
		ChatID: int(req.ChatId),
		Text:   req.Text,
	})
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &empty.Empty{}, nil
}

func (a *Api) FindOrCreateChat(ctx context.Context, req *FindOrCreateChatRequest) (*FindOrCreateChatResponse, error) {
	chatID, err := a.chatUserManager.GetUsersDialogChat(ctx, int(req.UserId_1), int(req.UserId_2))
	if err != nil {
		a.logger.Error("failed get users dialog chat", zap.Error(err))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if chatID == 0 {
		chatID, err = a.chatManager.Insert(ctx, &chat.Chat{
			Type: "dialog",
		})
		if err != nil {
			a.logger.Error("failed create chat", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		err = a.chatUserManager.Insert(ctx, &chatuser.ChatUser{
			UserID: int(req.UserId_1),
			ChatID: chatID,
		})

		if err != nil {
			a.logger.Error("failed create chat user", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		err = a.chatUserManager.Insert(ctx, &chatuser.ChatUser{
			UserID: int(req.UserId_2),
			ChatID: chatID,
		})

		if err != nil {
			a.logger.Error("failed create chat user", zap.Error(err))
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &FindOrCreateChatResponse{
		ChatId: int64(chatID),
	}, err
}

func (a *Api) messagesToProto(messages []message.Message) ([]*Message, error) {
	resMessages := make([]*Message, 0, len(messages))
	for _, v := range messages {
		ca, err := ptypes.TimestampProto(v.CreatedAt)
		if err != nil {
			return nil, err
		}
		var ua *timestamp.Timestamp
		if v.UpdatedAt.Valid {
			ua, err = ptypes.TimestampProto(v.UpdatedAt.Time)
			if err != nil {
				return nil, err
			}
		}
		var da *timestamp.Timestamp
		if v.DeletedAt.Valid {
			da, err = ptypes.TimestampProto(v.DeletedAt.Time)
			if err != nil {
				return nil, err
			}
		}
		resMessages = append(resMessages, &Message{
			Uuid:      v.UUID,
			UserId:    int64(v.UserID),
			ChatId:    int64(v.ChatID),
			Text:      v.Text,
			CreatedAt: ca,
			UpdatedAt: ua,
			DeletedAt: da,
		})
	}
	return resMessages, nil
}
