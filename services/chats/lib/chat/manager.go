package chat

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Manager struct {
	DB *mongo.Database
}

func NewManager(db *mongo.Database) *Manager {
	return &Manager{DB: db}
}

func (m *Manager) Insert(ctx context.Context, c *Chat) (string, error) {
	collection := m.DB.Collection("chats")

	res, err := collection.InsertOne(ctx, bson.M{
		"name": c.Name,
		"type": c.Type,
	})

	if err != nil {
		return "", err
	}

	// In MongoDB, the inserted ID is typically an ObjectID and not a number.
	// You'll need to adjust your models to work with that, or use a different strategy.
	// This code assumes that an ObjectID is fine and that it's okay to return it as a string.
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
