package user

import "time"

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
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)
