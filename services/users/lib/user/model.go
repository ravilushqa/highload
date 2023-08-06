package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Friend struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	City      string             `bson:"city"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Email         string
	Password      string
	FirstName     string `bson:"first_name"`
	LastName      string `bson:"last_name"`
	Birthday      time.Time
	Interests     string
	Sex           Sex
	City          string
	Subscriptions []*Friend  `bson:"subscriptions"`
	Subscribers   []*Friend  `bson:"subscribers"`
	CreatedAt     time.Time  `bson:"created_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty"`
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)

type Relation string

const (
	None         Relation = ""
	Subscriber   Relation = "subscriber"
	Subscription Relation = "subscription"
	Friends      Relation = "friends"
)
