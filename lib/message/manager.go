package message

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	Shard1 *sqlx.DB
	Shard2 *sqlx.DB
}

func NewManager(shard1 *sqlx.DB, shard2 *sqlx.DB) *Manager {
	return &Manager{Shard1: shard1, Shard2: shard2}
}

func (m *Manager) Insert(ctx context.Context, message *Message) error {
	shard := m.getShardByChatID(message.ChatID)
	query := `insert into messages 
		(user_id, chat_id, text)
		values (:user_id, :chat_id, :text)
	`
	_, err := shard.NamedExecContext(ctx, query, map[string]interface{}{
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
	firstChatShard := m.getShardByChatID(chatIDs[0])

	query := `
		select id, user_id, chat_id, text, created_at, updated_at, deleted_at 
		from messages where chat_id in (?)
	`

	query, args, err := sqlx.In(query, chatIDs)
	if err != nil {
		return nil, err
	}

	query = firstChatShard.Rebind(query)

	res := make([]Message, 0, 1024)

	if len(chatIDs) == 1 {
		err = firstChatShard.SelectContext(ctx, &res, query, args...)
		return res, err
	}
	err = m.Shard1.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	res2 := make([]Message, 0, 1024)
	err = m.Shard1.SelectContext(ctx, &res2, query, args...)
	if err != nil {
		return nil, err
	}
	res = append(res, res2...)

	return res, err
}

func (m *Manager) getShardByChatID(chatID int) *sqlx.DB {
	switch chatID % 2 {
	case 0:
		return m.Shard1
	}
	return m.Shard2
}
