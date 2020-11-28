package chats

import (
	"fmt"
	"net/http"

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
		w.Write([]byte("no chats"))
		return
	}
	messages, err := c.messageManager.GetChatMessages(r.Context(), chatIDs)
	if err != nil {
		c.logger.Error("failed get messages", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("something was wrong"))
		return
	}

	for _, m := range messages {
		w.Write([]byte(fmt.Sprintf("%d: %s", m.UserID, m.Text)))
	}

	//tmpl, err := template.ParseFiles(
	//	"resources/views/base.html",
	//	"resources/views/users/nav.html",
	//	"resources/views/users/index.html",
	//)
	//if err != nil {
	//	c.logger.Error("failed parse templates", zap.NamedError("error", err))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//_ = tmpl.ExecuteTemplate(w, "layout", struct {
	//	AuthUserID int
	//	Users      []user.User
	//}{uid, users})
}
