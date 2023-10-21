package chatuser

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatType string

const (
	TypeDialog ChatType = "dialog"
	TypeGroup  ChatType = "group"
)

type ChatUser struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	ChatID    string             `json:"chat_id" bson:"chat_id"`
	DeletedAt sql.NullTime       `json:"deleted_at" bson:"deleted_at"`
}
