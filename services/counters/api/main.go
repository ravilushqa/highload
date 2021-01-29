package main

import (
	"context"
	"net"

	"github.com/go-redis/redis"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	redisprovider "github.com/ravilushqa/highload/providers/redis"
	countersGrpc "github.com/ravilushqa/highload/services/counters/api/grpc"
)

func main() {
	_ = fx.New(
		fx.Provide(
			newConfig,
			zap.NewDevelopment,
			func(c *config) (*redis.Client, error) {
				return redisprovider.New(c.RedisURL)
			},
			NewApi,
		),
		fx.Invoke(
			registerHooks,
		),
	).Start(context.Background())
}

func registerHooks(lifecycle fx.Lifecycle, l *zap.Logger, r *redis.Client, a *Api) {
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
				countersGrpc.RegisterCountersServer(s, a)

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
				if err := r.ShutdownSave().Err(); err != nil {
					l.Error("failed to shutdown redis", zap.Error(err))
				}
				return l.Sync()
			},
		},
	)
}
