package chat

import "database/sql"

type ChatType string

const (
	TypeDialog ChatType = "dialog"
	TypeGroup  ChatType = "group"
)

type Chat struct {
	ID        int          `json:"id"`
	Type      ChatType     `json:"type"`
	Name      string       `json:"name"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
