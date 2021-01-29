package chats

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/services/chats/api/grpc"
	"github.com/ravilushqa/highload/services/web/lib"
)

type Controller struct {
	logger      *zap.Logger
	chatsClient grpc.ChatsClient
}

func NewController(logger *zap.Logger, chatsClient grpc.ChatsClient) *Controller {
	return &Controller{logger: logger, chatsClient: chatsClient}
}

func (c *Controller) Router(r chi.Router) chi.Router {
	return r.Route("/chats", func(r chi.Router) {
		r.Get("/", c.index)
		r.Get("/{chat_id}", c.show)
		r.Post("/{chat_id}/message", c.postMessage)
	})
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	res, err := c.chatsClient.GetUserChats(r.Context(), &grpc.GetUserChatsRequest{UserId: int64(uid)})
	if err != nil {
		c.logger.Error("failed get chats", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}
	if len(res.ChatIds) == 0 {
		_, _ = w.Write([]byte("no chats"))
		return
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/chat/nav.html",
		"resources/views/chat/index.html",
	)
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", struct {
		AuthUserID int
		ChatIDs    []int64
	}{uid, res.ChatIds})
	if err != nil {
		c.logger.Error("failed execute templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) show(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	chatID, err := strconv.Atoi(chi.URLParam(r, "chat_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong chat_id"))
		return
	}

	res, err := c.chatsClient.GetChatMessages(r.Context(), &grpc.GetChatMessagesRequest{ChatId: int64(chatID), UserId: int64(uid)})
	if err != nil {
		c.logger.Error("failed get messages", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	tmpl, err := template.ParseFiles(
		"resources/views/base.html",
		"resources/views/chat/nav.html",
		"resources/views/chat/show.html",
	)
	if err != nil {
		c.logger.Error("failed parse templates", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = tmpl.ExecuteTemplate(w, "layout", struct {
		AuthUserID int
		Messages   []*grpc.Message
		ChatID     int
	}{uid, res.Messages, chatID})
}

func (c *Controller) postMessage(w http.ResponseWriter, r *http.Request) {
	uid, _ := lib.GetAuthUserID(r.Context())
	chatID, err := strconv.Atoi(chi.URLParam(r, "chat_id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("wrong chat_id"))
		return
	}
	_ = r.ParseForm()

	_, err = c.chatsClient.StoreMessage(r.Context(), &grpc.StoreMessageRequest{
		UserId: int64(uid),
		ChatId: int64(chatID),
		Text:   r.FormValue("text"),
	})
	if err != nil {
		c.logger.Error("failed insert message", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
