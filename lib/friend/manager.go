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

func (m *Manager) GetFriendRequests(ctx context.Context, id int) ([]int, error) {
	var friendIDs []int
	err := m.DB.Select(&friendIDs, "select friend_id from friends where user_id = ? and approved = 0", id)
	return friendIDs, err
}

func (m *Manager) FriendRequest(ctx context.Context, requesterUser, addedUser int) error {
	query := `
		insert into friends
		(user_id, friend_id, approved)
		VALUES (:user_id, :friend_id, 0)
	`

	_, err := m.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"user_id":   addedUser,
		"friend_id": requesterUser,
	})

	return err
}

func (m *Manager) ApproveFriendRequest(ctx context.Context, approverUser, requesterUser int) error {
	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	updateQuery := `
		update friends
		set approved = 1
		where user_id = :user_id and friend_id = :friend_id
	`
	insertQuery := `
		insert into friends
		(user_id, friend_id, approved)
		VALUES ( :user_id, :friend_id, 1)
	`

	// approve request
	_, err = tx.NamedExecContext(ctx, updateQuery, map[string]interface{}{
		"user_id":   approverUser,
		"friend_id": requesterUser,
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	//link together
	_, err = tx.NamedExecContext(ctx, insertQuery, map[string]interface{}{
		"user_id":   requesterUser,
		"friend_id": approverUser,
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
