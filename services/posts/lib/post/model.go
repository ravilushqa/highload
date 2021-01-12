package post

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	Text      string       `json:"text"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}
