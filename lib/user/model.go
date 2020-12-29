package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Birthday  time.Time
	Interests string
	Sex       Sex
	City      string
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)
