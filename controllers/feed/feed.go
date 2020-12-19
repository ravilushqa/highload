package feed

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
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

	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("feed is empty"))
	}

	_, _ = w.Write([]byte(strings.Join(list, "\n")))
}
