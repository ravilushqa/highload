package chatuser

import (
	"context"
	"database/sql"

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

func (m *Manager) GetChatMembers(ctx context.Context, chatID int) ([]int64, error) {
	var userIDs []int64
	err := m.DB.SelectContext(ctx, &userIDs, "select user_id from chat_users where chat_id = ?", chatID)
	return userIDs, err
}

func (m Manager) GetUsersDialogChat(ctx context.Context, uID1, uID2 int) (int, error) {
	q := `select cu1.chat_id from chat_users cu1
		join chat_users as cu2 on  cu1.chat_id = cu2.chat_id and cu1.id != cu2.id
		join chats c on c.id = cu1.chat_id
		where cu1.user_id = ? and cu2.user_id = ? and c.type = 'dialog'
	`

	var dialogChat int

	err := m.DB.SelectContext(ctx, &dialogChat, q, uID1, uID2)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return dialogChat, nil
}
