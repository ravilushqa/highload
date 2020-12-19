package chats

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/chat"
	chatuser "github.com/ravilushqa/highload/lib/chat_user"
	"github.com/ravilushqa/highload/lib/message"
)

type Controller struct {
	logger          *zap.Logger
	chatManager     *chat.Manager
	chatUserManager *chatuser.Manager
	messageManager  *message.Manager
}

func NewController(logger *zap.Logger, chatManager *chat.Manager, chatUserManager *chatuser.Manager, messageManager *message.Manager) *Controller {
	return &Controller{logger: logger, chatManager: chatManager, chatUserManager: chatUserManager, messageManager: messageManager}
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
	chatIDs, err := c.chatUserManager.GetUserChats(r.Context(), uid)
	if err != nil {
		c.logger.Error("failed get chats", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}
	if len(chatIDs) == 0 {
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
		ChatIDs    []int
	}{uid, chatIDs})
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

	messages, err := c.messageManager.GetChatMessages(r.Context(), []int{chatID})
	if err != nil {
		c.logger.Error("failed get messages", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}
	_ = messages

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
		Messages   []message.Message
		ChatID     int
	}{uid, messages, chatID})
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

	err = c.messageManager.Insert(r.Context(), &message.Message{
		UserID: uid,
		ChatID: chatID,
		Text:   r.FormValue("text"),
	})
	if err != nil {
		c.logger.Error("failed insert message", zap.NamedError("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
