package friend

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Status string

const (
	None      Status = ""
	Added     Status = "added"
	Requested Status = "requested"
	Friends   Status = "friends"
)

type Manager struct {
	col *mongo.Collection
}

func New(db *mongo.Database) *Manager {
	return &Manager{col: db.Collection("friends")}
}

func (m *Manager) GetFriends(ctx context.Context, id string) ([]string, error) {
	filter := bson.M{"user_id": id, "approved": true}
	projection := bson.M{"friend_id": 1}

	cursor, err := m.col.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendIDs []string
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, result["friend_id"].(string))
	}

	return friendIDs, nil
}

func (m *Manager) GetFriendRequests(ctx context.Context, id string) ([]string, error) {
	filter := bson.M{"user_id": id, "approved": false}
	projection := bson.M{"friend_id": 1}

	cursor, err := m.col.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendIDs []string
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, result["friend_id"].(string))
	}

	return friendIDs, nil
}

func (m *Manager) FriendRequest(ctx context.Context, requesterUser, addedUser string) error {
	document := bson.M{
		"user_id":   addedUser,
		"friend_id": requesterUser,
		"approved":  false,
	}

	_, err := m.col.InsertOne(ctx, document)
	return err
}

func (m *Manager) ApproveFriendRequest(ctx context.Context, approverUser, requesterUser string) error {
	filter := bson.M{
		"user_id":   approverUser,
		"friend_id": requesterUser,
		"approved":  false,
	}

	update := bson.M{"$set": bson.M{"approved": true}}

	_, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	document := bson.M{
		"user_id":   requesterUser,
		"friend_id": approverUser,
		"approved":  true,
	}

	_, err = m.col.InsertOne(ctx, document)
	return err
}

// @todo refactor
func (m *Manager) GetRelation(ctx context.Context, authUser, user string) (Status, error) {
	filter := bson.M{
		"user_id":   authUser,
		"friend_id": user,
	}

	var result Friend
	err := m.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			filter = bson.M{
				"user_id":   user,
				"friend_id": authUser,
			}

			var result Friend
			err = m.col.FindOne(ctx, filter).Decode(&result)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return None, nil
				}
				return None, err
			}

			if result.Approved {
				return Friends, nil
			} else {
				return Added, nil
			}
		}
		return None, err
	}

	if result.Approved {
		return Friends, nil
	} else {
		return Requested, nil
	}
}
