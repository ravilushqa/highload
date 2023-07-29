package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Manager struct {
	l   *zap.Logger
	col *mongo.Collection
}

func New(l *zap.Logger, database *mongo.Database) *Manager {
	return &Manager{l: l, col: database.Collection("users")}
}

func (m *Manager) Store(ctx context.Context, user *User) (*primitive.ObjectID, error) {
	res, err := m.col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return &oid, nil
	}

	return nil, fmt.Errorf("failed to retrieve inserted ID")
}

func (m *Manager) GetByID(ctx context.Context, id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	result := m.col.FindOne(ctx, filter)

	user := &User{}
	err = result.Decode(user)

	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, return nil
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Manager) GetListByIds(ctx context.Context, ids []string) ([]User, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs[i] = objectID
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := m.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, cursor.Err()
}

func (m *Manager) GetByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}
	result := m.col.FindOne(ctx, filter)

	user := &User{}
	err := result.Decode(user)

	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, return nil
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Manager) GetAll(ctx context.Context, filter string) ([]User, error) {
	var queryFilter bson.M
	if filter != "" {
		queryFilter = bson.M{
			"$or": []bson.M{
				{"firstname": primitive.Regex{Pattern: "^" + filter, Options: "i"}},
				{"lastname": primitive.Regex{Pattern: "^" + filter, Options: "i"}},
			},
		}
	} else {
		queryFilter = bson.M{}
	}

	findOptions := options.Find()
	findOptions.SetLimit(100)

	cursor, err := m.col.Find(ctx, queryFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
