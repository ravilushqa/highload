package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string
	Password  string
	FirstName string
	LastName  string
	Birthday  time.Time
	Interests string
	Sex       Sex
	City      string
	CreatedAt time.Time
	DeletedAt *time.Time
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)
