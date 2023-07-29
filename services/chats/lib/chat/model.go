package chat

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatType string

const (
	TypeDialog ChatType = "dialog"
	TypeGroup  ChatType = "group"
)

type Chat struct {
	ID        primitive.ObjectID `json:"id" bson:",omitempty"`
	Type      ChatType           `json:"type"`
	Name      string             `json:"name"`
	DeletedAt sql.NullTime       `json:"deleted_at" bson:"deleted_at"`
}
