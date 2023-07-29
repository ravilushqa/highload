package message

import (
	"time"
)

type Message struct {
	UUID      string     `json:"uuid" bson:"uuid"`
	UserID    string     `json:"user_id" bson:"user_id"`
	ChatID    string     `json:"chat_id" bson:"chat_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}
