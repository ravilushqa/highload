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

func (m *Manager) Insert(ctx context.Context, c *Chat) (int, error) {
	query := `insert into chats 
		(name, type)
		values (:name, :type)
	`

	res, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"name": c.Name,
		"type": c.Type,
	})

	if err != nil {
		return 0, err
	}

	chatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(chatID), nil
}

func (m *Manager) GetByIDs(ctx context.Context, c *Chat) error {
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
