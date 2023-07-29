// nolint
package main

import (
	"context"
	"math/rand"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ravilushqa/highload/providers/mongodb"
)

type config struct {
	MONGO_URL string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	MONGO_DB  string `env:"MONGO_DB" envDefault:"highload"`
}

type User struct {
	Email     string    `faker:"email"`
	Password  string    `faker:"password"`
	FirstName string    `faker:"first_name" bson:"first_name"`
	LastName  string    `faker:"last_name" bson:"last_name"`
	Birthday  time.Time `faker:"birthday"`
	Interests string    `faker:"word"`
	Sex       Sex       `faker:"oneof: male, female"`
	City      string    `faker:"oneof: saint-petersburg, moscow, london, rome, oslo, stockholm, helsinki"`
	CreatedAt time.Time `faker:"created_at" bson:"created_at"`
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)

func CustomGenerator() {
	_ = faker.AddProvider("birthday", func(v reflect.Value) (interface{}, error) {
		min := 1
		max := 365 * 100
		t1 := time.Date(1920, 1, rand.Intn(max-min+1)+min, 0, 0, 0, 0, time.UTC)
		return t1, nil
	})
	_ = faker.AddProvider("created_at", func(v reflect.Value) (interface{}, error) {
		min := 1
		max := 365 * 100
		t1 := time.Date(2022, 1, rand.Intn(max-min+1)+min, 0, 0, 0, 0, time.UTC)
		return t1, nil
	})
}

// seed 1m data
func main() {
	c := new(config)
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}

	database, err := mongodb.New(context.Background(), c.MONGO_URL, c.MONGO_DB)
	if err != nil {
		return
	}
	col := database.Collection("users")
	CustomGenerator()
	u := User{}
	users := make([]any, 1000)

	for i := 0; i < 1000; i++ {
		println("iter: ", i)
		for j := 0; j < 1000; j++ {
			_ = faker.FakeData(&u)
			users[j] = u
		}
		_, err := col.InsertMany(context.Background(), users)
		if err != nil {
			panic(err)
		}
	}

}
