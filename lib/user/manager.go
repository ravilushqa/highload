package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"github.com/tarantool/go-tarantool"
)

type Manager struct {
	DB        *mssqlx.DBs
	tarantool *tarantool.Connection
}

func New(DB *mssqlx.DBs, tarantool *tarantool.Connection) *Manager {
	return &Manager{DB: DB, tarantool: tarantool}
}

func (m *Manager) Store(ctx context.Context, user *User) (int, error) {
	query := `
		insert into users
		(email, password, firstname, lastname, birthday, sex, interests, city)
		VALUES (:email, :password, :firstname, :lastname, :birthday, :sex, :interests, :city)
	`

	res, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"email":     user.Email,
		"password":  user.Password,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"birthday":  user.Birthday,
		"sex":       user.Sex,
		"interests": user.Interests,
		"city":      user.City,
	})

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (m *Manager) GetByID(ctx context.Context, id int) (*User, error) {
	//query := `
	//	select id, email, password, firstname, lastname, birthday, sex, interests, city
	//	from users
	//	where id = ? and deleted_at is null
	//`
	//res := &User{}
	//err := m.DB.GetContext(ctx, res, query, id)
	resp, err := m.tarantool.Select("mysqldata", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(id)})
	if err != nil {
		return nil, err
	}
	data := resp.Data[0].([]interface{})
	bd := data[5].(string)
	bdTime, err := time.Parse("2006-01-02", bd)
	if err != nil {
		return nil, fmt.Errorf("bd parse err %w", err)
	}

	res := &User{
		ID:        int(data[0].(uint64)),
		Email:     data[1].(string),
		Password:  data[2].(string),
		FirstName: data[3].(string),
		LastName:  data[4].(string),
		Birthday:  bdTime,
		Interests: data[6].(string),
		Sex:       Sex(data[7].(string)),
		City:      data[8].(string),
	}

	return res, err
}

func (m *Manager) GetListByIds(ctx context.Context, ids []int) ([]User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := `
		select id, email, password, firstname, lastname, birthday, sex, interests, city
		from users
		where id in (?) and deleted_at is null
	`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}

	query = m.DB.Rebind(query)

	res := make([]User, 0, len(ids))
	err = m.DB.SelectContext(ctx, &res, query, args...)

	return res, err
}

func (m *Manager) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		select id, email, password, firstname, lastname, birthday, sex, interests, city
		from users
		where email = ? and deleted_at is null
	`

	res := &User{}

	err := m.DB.GetContext(ctx, res, query, email)

	return res, err
}

func (m *Manager) GetAll(ctx context.Context, filter string) ([]User, error) {
	var query string
	res := make([]User, 0)

	if filter != "" {
		query = `
			select id, email, password, firstname, lastname, birthday, sex, interests, city
			from users
			where firstname like ? and lastname like ? and deleted_at is null
			order by id
		`
		err := m.DB.SelectContext(ctx, &res, query, fmt.Sprintf("%s%%", filter), fmt.Sprintf("%s%%", filter))
		return res, err
	}
	query = `
		select id, email, password, firstname, lastname, birthday, sex, interests, city
		from users
		where deleted_at is null
		limit 100 offset 0
	`

	err := m.DB.SelectContext(ctx, &res, query)
	return res, err
}
