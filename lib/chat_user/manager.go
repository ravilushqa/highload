package chatuser

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

func (m *Manager) Insert(ctx context.Context, cu *ChatUser) error {
	query := `insert into chat_users 
		(user_id, chat_id)
		values (:user_id, :chat_id)
	`
	_, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"user_id": cu.UserID,
		"chat_id": cu.ChatID,
	})

	return err
}

func (m *Manager) GetUserChats(ctx context.Context, userID int) ([]int, error) {
	var chatIDs []int
	err := m.DB.SelectContext(ctx, &chatIDs, "select chat_id from chat_users where user_id = ?", userID)
	return chatIDs, err
}
