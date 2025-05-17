package feed

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/services/posts/api/grpc"
	apiLib "github.com/ravilushqa/highload/services/web/api/lib"
	"github.com/ravilushqa/highload/services/web/lib"
)

type Controller struct {
	l           *zap.Logger
	postsClient grpc.PostsClient
}

func NewController(l *zap.Logger, postsClient grpc.PostsClient) *Controller {
	return &Controller{l: l, postsClient: postsClient}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/feed", func(r chi.Router) {
		r.Get("/", c.feed)
	})
}

func (c *Controller) feed(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())

	res, err := c.postsClient.GetFeed(r.Context(), &grpc.GetFeedRequest{UserId: int64(uid)})
	if err != nil {
		c.l.Error("failed get feed", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/feed/nav.html",
		"resources/views/feed/index.html",
	}

	data := struct {
		AuthUserID int
		Posts      []*grpc.Post
	}{uid, res.Posts}

	err = apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.l.Error("failed to render template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
