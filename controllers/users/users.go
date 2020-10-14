package users

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/friend"
	"github.com/ravilushqa/highload/lib/user"
)

type Controller struct {
	logger *zap.Logger
	u      *user.Manager
	f      *friend.Manager
}

func NewController(logger *zap.Logger, u *user.Manager, f *friend.Manager) *Controller {
	return &Controller{logger: logger, u: u, f: f}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/users", func(r chi.Router) {
		r.Get("/", c.index)
		r.HandleFunc("/{user_id}", c.profile)
	})
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetUsedIDFromCtx(r.Context())
	users, err := c.u.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something was wrong"))
		return
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/users/nav.html",
		"resources/views/users/index.html",
	)
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", struct {
		ID    int
		Users []user.User
	}{uid, users})
	return
}

func (c *Controller) profile(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("wrong user id"))
		return
	}

	u, err := c.u.GetByID(r.Context(), userID)
	// @todo check for no results
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something was wrong"))
		return
	}

	friendIds, err := c.f.GetFriends(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something was wrong"))
		return
	}

	friends, err := c.u.GetListByIds(r.Context(), friendIds)
	if err != nil {
		c.logger.Error("failed GetListByIds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something was wrong"))
		return
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/users/nav.html",
		"resources/views/users/profile.html",
	)
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", struct {
		*user.User
		Friends []user.User
	}{u, friends})
	return
}
