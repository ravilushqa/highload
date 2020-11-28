package message

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
)

type Manager struct {
	DB *mssqlx.DBs
}

func NewManager(db *mssqlx.DBs) *Manager {
	return &Manager{DB: db}
}

func (m *Manager) Insert(ctx context.Context, message *Message) error {
	query := `insert into messages 
		(user_id, chat_id, text)
		values (:user_id, :chat_id, :text)
	`
	_, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"user_id": message.UserID,
		"chat_id": message.ChatID,
		"text":    message.Text,
	})

	return err
}

func (m *Manager) GetChatMessages(ctx context.Context, chatIDs []int) ([]Message, error) {
	if len(chatIDs) == 0 {
		return nil, nil
	}
	query := `
		select id, user_id, chat_id, text, created_at, updated_at, deleted_at 
		from messages where chat_id in (?)
	`

	query, args, err := sqlx.In(query, chatIDs)
	if err != nil {
		return nil, err
	}

	query = m.DB.Rebind(query)

	res := make([]Message, 0, len(chatIDs))
	err = m.DB.SelectContext(ctx, &res, query, args...)

	return res, err
}
