package users

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/chat"
	chatuser "github.com/ravilushqa/highload/lib/chat_user"
	"github.com/ravilushqa/highload/lib/friend"
	"github.com/ravilushqa/highload/lib/user"
)

type Controller struct {
	logger *zap.Logger
	u      *user.Manager
	c      *chat.Manager
	f      *friend.Manager
	cu     *chatuser.Manager
}

func NewController(logger *zap.Logger, u *user.Manager, f *friend.Manager, cu *chatuser.Manager, c *chat.Manager) *Controller {
	return &Controller{logger: logger, u: u, f: f, cu: cu, c: c}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/users", func(r chi.Router) {
		r.Get("/", c.index)
		r.Route("/{user_id}", func(r chi.Router) {
			r.HandleFunc("/", c.profile)
			r.Post("/add", c.add)
			r.Post("/approve", c.approve)
			r.Post("/chat", c.chatOpen)
		})
	})
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	users, err := c.u.GetAll(r.Context(), r.URL.Query().Get("query"))
	if err != nil {
		c.logger.Error("failed get users", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
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

	_ = tmpl.ExecuteTemplate(w, "layout", struct {
		AuthUserID int
		Users      []user.User
	}{uid, users})
}

func (c *Controller) profile(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	u, err := c.u.GetByID(r.Context(), userID)
	// @todo check for no results
	if err != nil {
		c.logger.Error("failed get user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	friendIds, err := c.f.GetFriends(r.Context(), userID)
	if err != nil {
		c.logger.Error("failed get friends", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	friends, err := c.u.GetListByIds(r.Context(), friendIds)
	if err != nil {
		c.logger.Error("failed GetListByIds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	status, err := c.f.GetRelation(r.Context(), authUserID, u.ID)
	if err != nil {
		c.logger.Error("failed get relation", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	data := struct {
		AuthUserID int
		*user.User
		Friends []user.User
		Status  friend.Status
	}{authUserID, u, friends, status}

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

	_ = tmpl.ExecuteTemplate(w, "layout", data)
}

func (c *Controller) add(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	err = c.f.FriendRequest(r.Context(), authUserID, userID)
	if err != nil {
		c.logger.Error("failed find friend request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (c *Controller) approve(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	err = c.f.ApproveFriendRequest(r.Context(), authUserID, userID)
	if err != nil {
		c.logger.Error("approve friend request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (c *Controller) chatOpen(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}
	if userID == authUserID {
		_, _ = w.Write([]byte("own chat does not support"))
		return
	}

	chatID, err := c.cu.GetUsersDialogChat(r.Context(), authUserID, userID)
	if err != nil {
		c.logger.Error("failed get users dialog chat", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}
	if chatID == 0 {
		chatID, err = c.c.Insert(r.Context(), &chat.Chat{
			Type: "dialog",
		})
		if err != nil {
			c.logger.Error("failed create chat", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("something was wrong"))
			return
		}
		err = c.cu.Insert(r.Context(), &chatuser.ChatUser{
			UserID: authUserID,
			ChatID: chatID,
		})

		if err != nil {
			c.logger.Error("failed create chat user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("something was wrong"))
			return
		}
		err = c.cu.Insert(r.Context(), &chatuser.ChatUser{
			UserID: userID,
			ChatID: chatID,
		})

		if err != nil {
			c.logger.Error("failed create chat user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("something was wrong"))
			return
		}

	}

	http.Redirect(w, r, "/chats/"+strconv.Itoa(chatID), http.StatusFound)
}
