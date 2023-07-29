package main

import (
	"context"
	"net"

	_ "github.com/go-sql-driver/mysql"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ravilushqa/highload/providers/mongodb"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/users/lib/friend"
	"github.com/ravilushqa/highload/services/users/lib/user"
)

func main() {
	fx.New(
		fx.Provide(
			newConfig,
			zap.NewDevelopment,
			func(c *config) (*mongo.Database, error) {
				return mongodb.New(context.Background(), c.MongoURL, c.MongoDB)
			},
			NewApi,
			user.New,
			friend.New,
		),
		fx.Invoke(
			registerHooks,
		),
	).Run()
}

func registerHooks(lifecycle fx.Lifecycle, l *zap.Logger, db *mongo.Database, a *Api) {
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcprometheus.StreamServerInterceptor,
			grpczap.StreamServerInterceptor(l.Named("grpc_stream")),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcprometheus.UnaryServerInterceptor,
			grpczap.UnaryServerInterceptor(l.Named("grpc_unary")),
		)),
	)

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				addr := "0.0.0.0:50051"
				lis, err := net.Listen("tcp", addr) //@todo
				if err != nil {
					return err
				}

				usersGrpc.RegisterUsersServer(s, a)

				reflection.Register(s)

				l.Info("api started..", zap.String("addr", addr))

				go func() {
					l.Info("starting grpc server", zap.String("addr", addr))
					err := s.Serve(lis)
					if err != nil {
						l.Error("failed to serve grpc", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.GracefulStop()
				if err := db.Client().Disconnect(ctx); err != nil {
					l.Error("failed to disconnect from mongodb", zap.Error(err))
				}
				return l.Sync()
			},
		},
	)
}
