package chatuser

import "database/sql"

type ChatType string

const (
	TypeDialog ChatType = "dialog"
	TypeGroup  ChatType = "group"
)

type ChatUser struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	ChatID    int          `json:"chat_id"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
