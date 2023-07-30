package post

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Manager struct {
	col *mongo.Collection
}

func NewManager(db *mongo.Database) *Manager {
	return &Manager{col: db.Collection("posts")}
}

func (m *Manager) Insert(ctx context.Context, p *Post) (*Post, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	res, err := m.col.InsertOne(ctx, bson.M{
		"user_id":    p.UserID,
		"text":       p.Text,
		"created_at": p.CreatedAt,
	})
	if err != nil {
		return nil, err
	}

	p.ID = res.InsertedID.(primitive.ObjectID)

	return p, nil
}

func (m *Manager) GetOwnPosts(ctx context.Context, uid string) ([]*Post, error) {
	opts := options.Find()
	opts.SetLimit(100)
	opts.SetSort(bson.M{"_id": -1})

	cursor, err := m.col.Find(ctx, bson.M{
		"user_id": uid,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}, opts)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
