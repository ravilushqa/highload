package chat

import (
	"context"

	"github.com/linxGnu/mssqlx"
)

type Manager struct {
	DB *mssqlx.DBs
}

func NewManager(db *mssqlx.DBs) *Manager {
	return &Manager{DB: db}
}

func (m *Manager) Insert(ctx context.Context, c *Chat) error {
	query := `insert into chats 
		(name, type)
		values (:name, :type)
	`
	_, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"name": c.Name,
		"type": c.Type,
	})

	return err
}
