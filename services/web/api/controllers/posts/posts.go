package posts

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golang/protobuf/ptypes"
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

	if apiLib.IsHTMXRequest(r) {
		respPosts, err := c.postsClient.GetByUserID(r.Context(), &grpc.GetByUserIDRequest{UserId: int64(uid)})
		if err != nil {
			c.l.Error("failed to get posts after store", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		type postView struct {
			ID        int64
			Text      string
			CreatedAt string
		}
		list := make([]postView, len(respPosts.Posts))
		for i, p := range respPosts.Posts {
			tm, err := ptypes.Timestamp(p.CreatedAt)
			if err != nil {
				c.l.Error("failed to parse post timestamp", zap.Error(err))
			}
			list[i] = postView{ID: p.Id, Text: p.Text, CreatedAt: tm.Format("Jan 2 2006 15:04:05")}
		}
		listData := struct{ Posts []postView }{list}
		tpl, err := template.ParseFiles("resources/views/posts/list.html")
		if err != nil {
			c.l.Error("failed to parse posts list template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := tpl.ExecuteTemplate(w, "posts", listData); err != nil {
			c.l.Error("failed to render posts list", zap.NamedError("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
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

	type postView struct {
		ID        int64
		Text      string
		CreatedAt string
	}
	viewPosts := make([]postView, len(resp.Posts))
	for i, p := range resp.Posts {
		tm, err := ptypes.Timestamp(p.CreatedAt)
		if err != nil {
			c.l.Error("failed to parse post timestamp", zap.Error(err))
		}
		viewPosts[i] = postView{ID: p.Id, Text: p.Text, CreatedAt: tm.Format("Jan 2 2006 15:04:05")}
	}
	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/posts/nav.html",
		"resources/views/posts/index.html",
		"resources/views/posts/list.html",
	}
	data := struct {
		AuthUserID int
		Posts      []postView
	}{uid, viewPosts}
	err = apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.l.Error("failed to render template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
