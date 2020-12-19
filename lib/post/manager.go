package post

import (
	"context"

	"github.com/linxGnu/mssqlx"
)

type Manager struct {
	DB *mssqlx.DBs
}

func NewManager(DB *mssqlx.DBs) *Manager {
	return &Manager{DB: DB}
}

func (m *Manager) Insert(ctx context.Context, p *Post) (int64, error) {
	query := `insert into posts
		(user_id, text)
		values (:user_id, :text)
	`
	res, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"user_id": p.UserID,
		"text":    p.Text,
	})

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (m *Manager) GetAll(ctx context.Context) ([]Post, error) {
	var query string
	res := make([]Post, 0)

	query = `
		select id, user_id, text, created_at
		from posts
		where deleted_at is null
		order by id desc
		limit 100 offset 0
	`

	err := m.DB.SelectContext(ctx, &res, query)
	return res, err
}
