package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Text      string             `json:"text"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
