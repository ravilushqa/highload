package feed

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/post"
)

var cacheKey = "feed:user_id:%d"

type Controller struct {
	l     *zap.Logger
	redis *redis.Client
}

func NewController(l *zap.Logger, redis *redis.Client) *Controller {
	return &Controller{l: l, redis: redis}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/feed", func(r chi.Router) {
		r.Get("/", c.feed)
	})
}

func (c *Controller) feed(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())

	list, err := c.redis.LRange(fmt.Sprintf(cacheKey, uid), 0, 1000).Result()
	if err != nil {
		c.l.Error("failed to get feed from cache", zap.Error(err), zap.Int("user_id", uid))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}
	posts := make([]post.Post, 0, len(list))
	for _, jsonPost := range list {
		var p post.Post
		err = json.Unmarshal([]byte(jsonPost), &p)
		if err != nil {
			c.l.Error("failed unmarshal post", zap.Error(err), zap.Int("user_id", uid))
			continue
		}
		posts = append(posts, p)
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/feed/nav.html",
		"resources/views/feed/index.html",
	)
	if err != nil {
		c.l.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = tmpl.ExecuteTemplate(w, "layout", struct {
		AuthUserID int
		Posts      []post.Post
	}{uid, posts})
}
