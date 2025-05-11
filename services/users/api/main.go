package main

import (
	"context"
	"net"

	_ "github.com/go-sql-driver/mysql"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ravilushqa/highload/providers/db"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/users/lib/friend"
	"github.com/ravilushqa/highload/services/users/lib/user"
)

func main() {
	_ = fx.New(
		fx.Provide(
			newConfig,
			zap.NewDevelopment,
			func(c *config) (*mssqlx.DBs, error) {
				return db.New(c.DatabaseURL, c.SlavesUrls)
			},
			NewApi,
			user.New,
			friend.New,
		),
		fx.Invoke(
			registerHooks,
		),
	).Start(context.Background())
}

func registerHooks(lifecycle fx.Lifecycle, l *zap.Logger, db *mssqlx.DBs, a *Api) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				addr := ":50051"
				lis, err := net.Listen("tcp", addr) //@todo
				if err != nil {
					return err
				}

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
				usersGrpc.RegisterUsersServer(s, a)

				reflection.Register(s)

				l.Info("api started..", zap.String("addr", addr))

				defer s.GracefulStop()

				go func() {
					<-ctx.Done()
					s.Stop()
				}()

				return s.Serve(lis)
			},
			OnStop: func(context.Context) error {
				if errs := db.Destroy(); len(errs) > 0 {
					l.Error("failed to close db", zap.Errors("errors", errs))
				}
				return l.Sync()
			},
		},
	)
}
