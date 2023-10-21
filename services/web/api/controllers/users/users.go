package users

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	chatsGrpc "github.com/ravilushqa/highload/services/chats/api/grpc"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	"github.com/ravilushqa/highload/services/web/lib"
)

type Controller struct {
	logger      *zap.Logger
	chatsClient chatsGrpc.ChatsClient
	usersClient usersGrpc.UsersClient
}

func NewController(logger *zap.Logger, chatsClient chatsGrpc.ChatsClient, usersClient usersGrpc.UsersClient) *Controller {
	return &Controller{logger: logger, chatsClient: chatsClient, usersClient: usersClient}
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
	res, err := c.usersClient.GetAll(r.Context(), &usersGrpc.GetUsersRequest{Filter: r.URL.Query().Get("query")})
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
		AuthUserID string
		Users      []*usersGrpc.User
	}{uid, res.Users})
}

func (c *Controller) profile(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID := chi.URLParam(r, "user_id")

	getUserResponse, err := c.usersClient.GetById(r.Context(), &usersGrpc.GetByIdRequest{UserId: userID})
	// @todo check for no results
	if err != nil {
		c.logger.Error("failed get user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	getFriendsIdsResponse, err := c.usersClient.GetFriendsIds(r.Context(), &usersGrpc.GetFriendsIdsRequest{UserId: userID})
	if err != nil {
		c.logger.Error("failed get friends", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	getListByIdsResponse, err := c.usersClient.GetListByIds(r.Context(), &usersGrpc.GetListByIdsRequest{UserIds: getFriendsIdsResponse.UserIds})
	if err != nil {
		c.logger.Error("failed GetListByIds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	status, err := c.usersClient.GetRelation(r.Context(), &usersGrpc.GetRelationRequest{
		FromUserId: authUserID,
		ToUserId:   getUserResponse.User.Id,
	})
	if err != nil {
		c.logger.Error("failed get relation", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	data := struct {
		AuthUserID string
		*usersGrpc.User
		Friends []*usersGrpc.User
		Status  usersGrpc.UserRelation
	}{authUserID, getUserResponse.User, getListByIdsResponse.Users, status.Relation}

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
	userID := chi.URLParam(r, "user_id")

	_, err := c.usersClient.FriendRequest(r.Context(), &usersGrpc.FriendRequestRequest{
		RequesterUserId: authUserID,
		AddedUserId:     userID,
	})
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
	userID := chi.URLParam(r, "user_id")
	_, err := c.usersClient.ApproveFriendRequest(r.Context(), &usersGrpc.ApproveFriendRequestRequest{
		ApproverUserId:  authUserID,
		RequesterUserId: userID,
	})
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
	userID := chi.URLParam(r, "user_id")

	if userID == authUserID {
		_, _ = w.Write([]byte("own chat does not support"))
		return
	}

	res, err := c.chatsClient.FindOrCreateChat(r.Context(), &chatsGrpc.FindOrCreateChatRequest{
		UserId_1: authUserID,
		UserId_2: userID,
	})
	if err != nil {
		c.logger.Error("failed find or create chat", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/chats/%s", res.ChatId), http.StatusFound)
}
