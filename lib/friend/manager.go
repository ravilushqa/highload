package friend

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

func (m *Manager) GetFriends(ctx context.Context, id int) ([]int, error) {
	var friendIDs []int
	err := m.DB.SelectContext(ctx, &friendIDs, "select friend_id from friends where user_id = ? and approved = 1", id)
	return friendIDs, err
}

func (m *Manager) GetFriendRequests(id int) ([]int, error) {
	var friendIDs []int
	err := m.DB.Select(&friendIDs, "select friend_id from friends where user_id = ? and approved = 0", id)
	return friendIDs, err
}
