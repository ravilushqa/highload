package posts

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/post"
	kafkaproducerprovider "github.com/ravilushqa/highload/providers/kafka-producer"
)

type Controller struct {
	l  *zap.Logger
	pm *post.Manager
	kp *kafkaproducerprovider.KafkaProducer
}

func NewController(l *zap.Logger, pm *post.Manager, kp *kafkaproducerprovider.KafkaProducer) *Controller {
	return &Controller{l: l, pm: pm, kp: kp}
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

	p := &post.Post{
		UserID: uid,
		Text:   text,
	}
	p, err := c.pm.Insert(r.Context(), p)
	if err != nil {
		c.l.Error("failed to insert post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	message, err := json.Marshal(p)
	if err != nil {
		c.l.Error("failed to marshal post", zap.Error(err))
	} else if err = c.kp.SendMessage(message, nil); err != nil {
		c.l.Error("failed to send message to kafka", zap.Error(err))
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	posts, err := c.pm.GetOwnPosts(r.Context(), uid)
	if err != nil {
		c.l.Error("failed get users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/posts/nav.html",
		"resources/views/posts/index.html",
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
