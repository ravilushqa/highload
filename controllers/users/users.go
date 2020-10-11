package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib/user"
)

type Controller struct {
	logger *zap.Logger
	m      *user.Manager
}

func NewController(logger *zap.Logger, m *user.Manager) *Controller {
	return &Controller{logger: logger, m: m}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/users", func(r chi.Router) {
		r.HandleFunc("/{user_id}", c.profile)
	})
}

func (c *Controller) profile(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("wrong user id"))
		return
	}

	u, err := c.m.GetByID(r.Context(), userID)
	// @todo check for no results
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something was wrong"))
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", *u)))
	return
}
