package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
)

type Manager struct {
	DB *mssqlx.DBs
}

func New(DB *mssqlx.DBs) *Manager {
	return &Manager{DB: DB}
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
	query := `
		select id, email, password, firstname, lastname, birthday, sex, interests, city
		from users
		where id = ? and deleted_at is null
	`

	res := &User{}

	err := m.DB.GetContext(ctx, res, query, id)

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
