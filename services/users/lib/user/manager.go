package user

import (
	"context"
	"errors"
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
	if user.Subscriptions == nil {
		user.Subscriptions = make([]*Friend, 0)
	}
	if user.Subscribers == nil {
		user.Subscribers = make([]*Friend, 0)
	}
	res, err := m.col.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return &oid, nil
	}

	return nil, fmt.Errorf("failed to retrieve inserted ID")
}

func (m *Manager) GetByID(ctx context.Context, id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	result := m.col.FindOne(ctx, filter)

	user := &User{}
	err = result.Decode(user)

	if !errors.Is(err, mongo.ErrNoDocuments) {
		if err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
	} else {
		return nil, nil // Not found, return nil
	}

	return user, nil
}

func (m *Manager) GetListByIds(ctx context.Context, ids []string) ([]User, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid object ID: %w", err)
		}
		objectIDs[i] = objectID
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := m.col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, cursor.Err()
}

func (m *Manager) GetByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}
	result := m.col.FindOne(ctx, filter)

	user := &User{}
	err := result.Decode(user)

	if errors.Is(err, mongo.ErrNoDocuments) {
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
		return nil, fmt.Errorf("failed to find: %w", err)
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate: %w", err)
	}

	return users, nil
}

func (m *Manager) GetFriends(ctx context.Context, id primitive.ObjectID) ([]*Friend, error) {
	cursor, err := m.col.Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"_id": id}},
		bson.M{"$project": bson.M{"friends": bson.M{
			"$setIntersection": bson.A{"$subscriptions", "$subscribers"},
		}}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var res struct {
		Friends []*Friend `bson:"friends"`
	}
	for cursor.Next(ctx) {
		if err := cursor.Decode(&res); err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate: %w", err)
	}

	return res.Friends, nil
}

func (m *Manager) Subscribe(ctx context.Context, userID, subscriberID primitive.ObjectID) error {
	cur, err := m.col.Find(ctx, bson.M{"_id": bson.M{"$in": []primitive.ObjectID{userID, subscriberID}}})
	if err != nil {
		return fmt.Errorf("failed to find users: %w", err)
	}
	defer cur.Close(ctx)

	var users []User
	for cur.Next(ctx) {
		var user User
		if err := cur.Decode(&user); err != nil {
			return fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return fmt.Errorf("failed to iterate users: %w", err)
	}

	if len(users) != 2 {
		return fmt.Errorf("unexpected number of users: %d", len(users))
	}

	// @TODO: optimize?!
	var user, subscriber *Friend
	if users[0].ID == userID {
		user = &Friend{ID: userID, FirstName: users[0].FirstName, LastName: users[0].LastName, City: users[0].City}
		subscriber = &Friend{ID: subscriberID, FirstName: users[1].FirstName, LastName: users[1].LastName, City: users[1].City}
	} else {
		user = &Friend{ID: userID, FirstName: users[1].FirstName, LastName: users[1].LastName, City: users[1].City}
		subscriber = &Friend{ID: subscriberID, FirstName: users[0].FirstName, LastName: users[0].LastName, City: users[0].City}
	}

	_, err = m.col.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"subscriptions": subscriber}})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	_, err = m.col.UpdateOne(ctx, bson.M{"_id": subscriberID}, bson.M{"$addToSet": bson.M{"subscribers": user}})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetRelations returns the relation between two users. Possible values are: "none", "follower", "following", "friend"
func (m *Manager) GetRelations(ctx context.Context, userID, otherUserID primitive.ObjectID) (Relation, error) {
	user, err := m.GetByID(ctx, userID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if contains(user.Subscribers, otherUserID) && contains(user.Subscriptions, otherUserID) {
		return Friends, nil
	}

	if contains(user.Subscribers, otherUserID) {
		return Subscriber, nil
	}

	if contains(user.Subscriptions, otherUserID) {
		return Subscription, nil
	}

	return None, nil
}

func contains(s []*Friend, e primitive.ObjectID) bool {
	for _, a := range s {
		if a.ID == e {
			return true
		}
	}
	return false
}
