package users

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	chatsGrpc "github.com/ravilushqa/highload/services/chats/api/grpc"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
	apiLib "github.com/ravilushqa/highload/services/web/api/lib"
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
			r.Get("/friends", c.getFriendsList) // New endpoint for HTMX to fetch updated friends list
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

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/users/nav.html",
		"resources/views/users/index.html",
	}

	data := struct {
		AuthUserID int
		Users      []*usersGrpc.User
	}{uid, res.Users}

	err = apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.logger.Error("failed to render template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) profile(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	getUserResponse, err := c.usersClient.GetById(r.Context(), &usersGrpc.GetByIdRequest{UserId: int64(userID)})
	// @todo check for no results
	if err != nil {
		c.logger.Error("failed get user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	getFriendsIdsResponse, err := c.usersClient.GetFriendsIds(r.Context(), &usersGrpc.GetFriendsIdsRequest{UserId: int64(userID)})
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
		FromUserId: int64(authUserID),
		ToUserId:   getUserResponse.User.Id,
	})
	if err != nil {
		c.logger.Error("failed get relation", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	data := struct {
		AuthUserID int
		*usersGrpc.User
		Friends []*usersGrpc.User
		Status  usersGrpc.UserRelation
	}{authUserID, getUserResponse.User, getListByIdsResponse.Users, status.Relation}

	templateFiles := []string{
		"resources/views/base.html",
		"resources/views/users/nav.html",
		"resources/views/users/profile.html",
	}

	err = apiLib.RenderTemplate(w, r, templateFiles, data)
	if err != nil {
		c.logger.Error("failed to render template", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) add(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	_, err = c.usersClient.FriendRequest(r.Context(), &usersGrpc.FriendRequestRequest{
		RequesterUserId: int64(authUserID),
		AddedUserId:     int64(userID),
	})
	if err != nil {
		c.logger.Error("failed find friend request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	// If this is an HTMX request, update just the friend status button
	if apiLib.IsHTMXRequest(r) {
		// Get the current relationship status after the update
		status, err := c.usersClient.GetRelation(r.Context(), &usersGrpc.GetRelationRequest{
			FromUserId: int64(authUserID),
			ToUserId:   int64(userID),
		})
		if err != nil {
			c.logger.Error("failed get relation", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Render just the friend actions section
		data := struct {
			AuthUserID int
			Id         int64
			Status     usersGrpc.UserRelation
		}{authUserID, int64(userID), status.Relation}

		tmpl, err := template.ParseFiles("resources/views/users/friend_actions.html")
		if err != nil {
			c.logger.Error("failed to parse template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			c.logger.Error("failed to render template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	_, err = c.usersClient.ApproveFriendRequest(r.Context(), &usersGrpc.ApproveFriendRequestRequest{
		ApproverUserId:  int64(authUserID),
		RequesterUserId: int64(userID),
	})
	if err != nil {
		c.logger.Error("approve friend request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	// If this is an HTMX request, update just the friend status button
	if apiLib.IsHTMXRequest(r) {
		// Get the current relationship status after the update
		status, err := c.usersClient.GetRelation(r.Context(), &usersGrpc.GetRelationRequest{
			FromUserId: int64(authUserID),
			ToUserId:   int64(userID),
		})
		if err != nil {
			c.logger.Error("failed get relation", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Render just the friend actions section
		data := struct {
			AuthUserID int
			Id         int64
			Status     usersGrpc.UserRelation
		}{authUserID, int64(userID), status.Relation}

		tmpl, err := template.ParseFiles("resources/views/users/friend_actions.html")
		if err != nil {
			c.logger.Error("failed to parse template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			c.logger.Error("failed to render template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// After approval, fetch and update the friends list section
		getFriendsIdsResponse, err := c.usersClient.GetFriendsIds(r.Context(), &usersGrpc.GetFriendsIdsRequest{UserId: int64(userID)})
		if err != nil {
			c.logger.Error("failed get friends", zap.Error(err))
			return
		}

		_, err = c.usersClient.GetListByIds(r.Context(), &usersGrpc.GetListByIdsRequest{UserIds: getFriendsIdsResponse.UserIds})
		if err != nil {
			c.logger.Error("failed GetListByIds", zap.Error(err))
			return
		}

		// Add trigger to update friend list on the page
		w.Header().Add("HX-Trigger", "friendsUpdated")

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

	res, err := c.chatsClient.FindOrCreateChat(r.Context(), &chatsGrpc.FindOrCreateChatRequest{
		UserId_1: int64(authUserID),
		UserId_2: int64(userID),
	})

	if err != nil {
		c.logger.Error("failed find or create chat", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	http.Redirect(w, r, "/chats/"+strconv.Itoa(int(res.ChatId)), http.StatusFound)
}

// getFriendsList renders just the friends list section for HTMX updates
func (c *Controller) getFriendsList(w http.ResponseWriter, r *http.Request) {
	authUserID, _ := lib.GetAuthUserID(r.Context())
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong user id"))
		return
	}

	// Get user info
	getUserResponse, err := c.usersClient.GetById(r.Context(), &usersGrpc.GetByIdRequest{UserId: int64(userID)})
	if err != nil {
		c.logger.Error("failed get user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get friends list
	getFriendsIdsResponse, err := c.usersClient.GetFriendsIds(r.Context(), &usersGrpc.GetFriendsIdsRequest{UserId: int64(userID)})
	if err != nil {
		c.logger.Error("failed get friends", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	getListByIdsResponse, err := c.usersClient.GetListByIds(r.Context(), &usersGrpc.GetListByIdsRequest{UserIds: getFriendsIdsResponse.UserIds})
	if err != nil {
		c.logger.Error("failed GetListByIds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		AuthUserID int
		Id         int64
		Friends    []*usersGrpc.User
	}{authUserID, getUserResponse.User.Id, getListByIdsResponse.Users}

	// Parse the friends list template
	tmpl, err := template.ParseFiles("resources/views/users/friends_list.html")
	if err != nil {
		c.logger.Error("failed to parse template", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render the template
	err = tmpl.Execute(w, data)
	if err != nil {
		c.logger.Error("failed to render template", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
