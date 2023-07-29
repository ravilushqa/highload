package message

import (
	"context"

	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Manager struct {
	col *mongo.Collection
}

func NewManager(db *mongo.Database) *Manager {
	return &Manager{col: db.Collection("messages")}
}

func (m *Manager) Insert(ctx context.Context, message *Message) (string, error) {
	id := uuid.NewV4().String()

	_, err := m.col.InsertOne(ctx, bson.M{
		"uuid":    id,
		"user_id": message.UserID,
		"chat_id": message.ChatID,
		"text":    message.Text,
	})
	return id, err
}

func (m *Manager) HardDeleteLastMessage(ctx context.Context, chatID, userID, text string) error {
	filter := bson.M{
		"chat_id": chatID,
		"user_id": userID,
		"text":    text,
	}

	opts := options.FindOneAndDelete().SetSort(bson.M{"_id": -1}) // sorting by _id as a proxy for created_at
	err := m.col.FindOneAndDelete(ctx, filter, opts).Err()

	return err
}

func (m *Manager) GetChatMessages(ctx context.Context, chatIDs []string) ([]Message, error) {
	if len(chatIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"chat_id": bson.M{"$in": chatIDs},
	}

	cursor, err := m.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []Message
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
