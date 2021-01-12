package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/linxGnu/mssqlx"
	"github.com/neonxp/rutina"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	chatsGrpc "github.com/ravilushqa/highload/services/chats/api/grpc"
	postsGrpc "github.com/ravilushqa/highload/services/posts/api/grpc"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/web/api/controllers/auth"
	"github.com/ravilushqa/highload/services/web/api/controllers/chats"
	"github.com/ravilushqa/highload/services/web/api/controllers/feed"
	"github.com/ravilushqa/highload/services/web/api/controllers/posts"
	"github.com/ravilushqa/highload/services/web/api/controllers/users"
	"github.com/ravilushqa/highload/services/web/lib"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*lib.Auth, error) {
			return lib.NewAuth(c.JwtSecret), nil
		},
		func(c *config) (chatsGrpc.ChatsClient, error) {
			conn, err := grpc.Dial(c.ChatsURL, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
			return chatsGrpc.NewChatsClient(conn), nil
		},
		func(c *config) (postsGrpc.PostsClient, error) {
			conn, err := grpc.Dial(c.PostsURL, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
			return postsGrpc.NewPostsClient(conn), nil
		},
		func(c *config) (usersGrpc.UsersClient, error) {
			conn, err := grpc.Dial(c.UsersURL, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
			return usersGrpc.NewUsersClient(conn), nil
		},
		NewAPI,
		auth.NewController,
		users.NewController,
		chats.NewController,
		posts.NewController,
		feed.NewController,
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, container.Invoke(func(a *API) {})
}

func main() {
	tl, _ := zap.NewDevelopment()
	container, err := buildContainer()
	if err != nil {
		tl.Fatal("cannot build depends", zap.Error(err))
	}

	r := rutina.New(rutina.WithErrChan())
	go func() {
		for err := range r.Errors() {
			tl.Error("runtime error", zap.Error(err))
		}
	}()

	err = container.Invoke(func(api *API) {
		r.Go(api.run)
		r.ListenOsSignals()
	})
	if err != nil {
		tl.Fatal("invoke failed", zap.Error(err))
	}

	if err := r.Wait(); err != nil {
		tl.Error("run failed", zap.Error(err))
	}

	err = container.Invoke(func(l *zap.Logger, db *mssqlx.DBs) error {
		if errs := db.Destroy(); len(errs) > 0 {
			l.Error("failed to close db", zap.Errors("errors", errs))
		}
		l.Info("gracefully shutdown...")
		return l.Sync()
	})
	if err != nil {
		tl.Error("shutdown failed", zap.Error(err))
	}

}
