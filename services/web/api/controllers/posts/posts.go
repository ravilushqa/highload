package posts

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
	return r.Route("/posts", func(r chi.Router) {
		r.Get("/", c.index)
		r.Post("/", c.Store)
	})
}

func (c *Controller) Store(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	_ = r.ParseForm()

	text := r.FormValue("text")
	if text == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong text"))
		return
	}

	_, err := c.postsClient.Store(r.Context(), &grpc.StoreRequest{
		UserId: int64(uid),
		Text:   text,
	})
	if err != nil {
		c.l.Error("failed to store post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	resp, err := c.postsClient.GetByUserID(r.Context(), &grpc.GetByUserIDRequest{UserId: int64(uid)})
	if err != nil {
		c.l.Error("failed get users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/posts/nav.html",
		"resources/views/posts/index.html",
	}

	data := struct {
		AuthUserID int
		Posts      []*grpc.Post
	}{uid, resp.Posts}

	err = apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.l.Error("failed to render template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
