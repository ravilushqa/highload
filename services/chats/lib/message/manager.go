package message

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/neonxp/rutina"
	"github.com/satori/go.uuid"
)

type Manager struct {
	Shards []*sqlx.DB
}

func NewManager(shards []*sqlx.DB) *Manager {
	return &Manager{Shards: shards}
}

func (m *Manager) Insert(ctx context.Context, message *Message) (string, error) {
	shard := m.getShardByChatID(message.ChatID)
	query := `insert into messages 
		(uuid, user_id, chat_id, text)
		values (:uuid, :user_id, :chat_id, :text)
	`

	id := uuid.NewV4().String()

	_, err := shard.NamedExecContext(ctx, query, map[string]interface{}{
		"uuid":    id,
		"user_id": message.UserID,
		"chat_id": message.ChatID,
		"text":    message.Text,
	})
	return id, err
}

func (m *Manager) HardDelete(ctx context.Context, chatID int, uuid string) error {
	shard := m.getShardByChatID(chatID)
	query := `delete from messages where uuid = ?`

	_, err := shard.NamedExecContext(ctx, query, map[string]interface{}{
		"uuid": uuid,
	})
	return err
}

func (m *Manager) GetChatMessages(ctx context.Context, chatIDs []int) ([]Message, error) {
	if len(chatIDs) == 0 {
		return nil, nil
	}
	firstChatShard := m.getShardByChatID(chatIDs[0])

	query := `
		select uuid, user_id, chat_id, text, created_at, updated_at, deleted_at 
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

	r := rutina.New()

	for _, shard := range m.Shards {
		r.Go(func(ctx context.Context) error {
			shardData := make([]Message, 0, 1024)
			err := shard.SelectContext(ctx, &shardData, query, args...)

			if err != nil {
				return err
			}
			res = append(res, shardData...)
			return nil
		})
	}
	err = r.Wait()
	return res, err
}

func (m *Manager) getShardByChatID(chatID int) *sqlx.DB {
	fmt.Println(len(m.Shards))
	return m.Shards[chatID%len(m.Shards)]
}
