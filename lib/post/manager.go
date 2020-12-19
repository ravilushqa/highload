package post

import (
	"context"
	"time"

	"github.com/linxGnu/mssqlx"
)

type Manager struct {
	DB *mssqlx.DBs
}

func NewManager(DB *mssqlx.DBs) *Manager {
	return &Manager{DB: DB}
}

func (m *Manager) Insert(ctx context.Context, p *Post) (*Post, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	query := `insert into posts
		(user_id, text, created_at)
		values (:user_id, :text, :created_at)
	`
	res, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"user_id":    p.UserID,
		"text":       p.Text,
		"created_at": p.CreatedAt,
	})

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	p.ID = int(id)

	return p, nil
}

func (m *Manager) GetOwnPosts(ctx context.Context, uid int) ([]Post, error) {
	var query string
	res := make([]Post, 0)

	query = `
		select id, user_id, text, created_at
		from posts
		where deleted_at is null and user_id = ?
		order by id desc
		limit 100 offset 0
	`

	err := m.DB.SelectContext(ctx, &res, query, uid)
	return res, err
}
