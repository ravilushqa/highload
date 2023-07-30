package chatuser

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Manager struct {
	col *mongo.Collection
}

func NewManager(db *mongo.Database) *Manager {
	return &Manager{col: db.Collection("chat_users")}
}

func (m *Manager) Insert(ctx context.Context, cu *ChatUser) error {
	_, err := m.col.InsertOne(ctx, bson.M{
		"user_id": cu.UserID,
		"chat_id": cu.ChatID,
	})

	return err
}

func (m *Manager) GetUserChats(ctx context.Context, userID string) ([]string, error) {
	cursor, err := m.col.Find(ctx, bson.M{"user_id": userID})

	var chatIDs []string
	if err != nil {
		return chatIDs, err
	}

	for cursor.Next(ctx) {
		var elem ChatUser
		err := cursor.Decode(&elem)
		if err != nil {
			return chatIDs, err
		}

		chatIDs = append(chatIDs, elem.ChatID)
	}

	if err = cursor.Err(); err != nil {
		return chatIDs, err
	}

	cursor.Close(ctx)

	return chatIDs, nil
}

func (m *Manager) GetChatMembers(ctx context.Context, chatID string) ([]string, error) {
	cursor, err := m.col.Find(ctx, bson.M{"chat_id": chatID})

	var userIDs []string
	if err != nil {
		return userIDs, err
	}

	for cursor.Next(ctx) {
		var elem ChatUser
		err := cursor.Decode(&elem)
		if err != nil {
			return userIDs, err
		}

		userIDs = append(userIDs, elem.UserID)
	}

	if err = cursor.Err(); err != nil {
		return userIDs, err
	}

	cursor.Close(ctx)

	return userIDs, nil
}

// @todo: optimize
func (m *Manager) GetUsersDialogChat(ctx context.Context, uID1, uID2 string) (string, error) {
	cursor, err := m.col.Find(ctx, bson.M{"user_id": uID1})
	if err != nil {
		return "", err
	}

	var chatID string
	for cursor.Next(ctx) {
		var elem ChatUser
		err := cursor.Decode(&elem)
		if err != nil {
			return "", err
		}

		chatIDs, err := m.GetUserChats(ctx, uID2)
		if err != nil {
			return "", err
		}

		for _, chatID := range chatIDs {
			if chatID == elem.ChatID {
				return chatID, nil
			}
		}
	}

	if err = cursor.Err(); err != nil {
		return "", err
	}

	cursor.Close(ctx)

	return chatID, nil
}
