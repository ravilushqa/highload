package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	DB *sqlx.DB
}

func New(DB *sqlx.DB) *Manager {
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
