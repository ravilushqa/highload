//nolint
package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type config struct {
	DatabaseUrl string `env:"DATABASE_URL" envDefault:"user:password@(localhost:3306)/app"`
}

type User struct {
	ID        int       `faker:"-"`
	Email     string    `faker:"email"`
	Password  string    `faker:"password"`
	FirstName string    `faker:"first_name"`
	LastName  string    `faker:"last_name"`
	Birthday  time.Time `faker:"birthday"`
	Interests string    `faker:"word"`
	Sex       Sex       `faker:"oneof: male, female"`
	City      string    `faker:"oneof: saint-petersburg, moscow, london, rome, oslo, stockholm, helsinki"`
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
}

// seed 1m data
func main() {
	c := new(config)
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("mysql", fmt.Sprint(c.DatabaseUrl, "?parseTime=true"))
	if err != nil {
		panic(err)
	}
	CustomGenerator()
	u := User{}

	for i := 0; i < 1000; i++ {
		println("iter: ", i)

		sqlStr := `insert into users (email, password, firstname, lastname, birthday, sex, interests, city) VALUES`
		var vals []interface{}
		for j := 0; j < 1000; j++ {

			_ = faker.FakeData(&u)
			sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?),"
			vals = append(vals, u.Email, u.Password, u.FirstName, u.LastName, u.Birthday, u.Sex, u.Interests, u.City)
		}
		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//prepare the statement
		stmt, _ := db.Prepare(sqlStr)

		//format all vals at once
		res, _ := stmt.Exec(vals...)
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		println("count: ", count)

	}

}
