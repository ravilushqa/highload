package message

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	ChatID    int          `json:"chat_id" db:"chat_id"`
	Text      string       `json:"text"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}
