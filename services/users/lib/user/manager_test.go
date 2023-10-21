package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	mongocontainer "github.com/ravilushqa/highload/pkg/testcontainers/mongo"
)

func TestManager_GetFriends(t *testing.T) {
	ctx := context.Background()

	container, err := mongocontainer.StartContainer(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	})

	endpoint, err := container.Endpoint(ctx, "mongodb")
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	db := client.Database("test")

	manager := New(zap.NewNop(), db)
	err = db.Drop(ctx)
	require.NoError(t, err)

	testFriends := []*Friend{
		{
			ID:        primitive.NewObjectID(),
			FirstName: "Jane",
			LastName:  "Doe",
			City:      "New York",
		},
		{
			ID:        primitive.NewObjectID(),
			FirstName: "Jack",
			LastName:  "Doe",
			City:      "New York",
		},
	}

	id, err := manager.Store(ctx, &User{
		FirstName:     "John",
		LastName:      "Doe",
		Subscriptions: testFriends,
		Subscribers: []*Friend{
			{
				ID:        primitive.NewObjectID(),
				FirstName: "Ricky",
				LastName:  "Doe",
				City:      "New York",
			},
			testFriends[0],
		},
	})
	require.NoError(t, err)

	friends, err := manager.GetFriends(ctx, *id)
	require.NoError(t, err)
	require.Len(t, friends, 1)
	require.Equal(t, testFriends[0].ID, friends[0].ID)
}

func TestGetRelation(t *testing.T) {
	ctx := context.Background()

	container, err := mongocontainer.StartContainer(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	})

	endpoint, err := container.Endpoint(ctx, "mongodb")
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)
	db := client.Database("test-db")
	manager := New(zap.NewNop(), db)
	err = db.Drop(ctx)
	require.NoError(t, err)
	t.Run("None", func(t *testing.T) {
		firstID, err := manager.Store(ctx, &User{})
		require.NoError(t, err)
		secondID, err := manager.Store(ctx, &User{})
		require.NoError(t, err)
		relation, err := manager.GetRelations(ctx, *firstID, *secondID)

		require.NoError(t, err)
		require.Equal(t, None, relation)
	})

	t.Run("Subscription and Subscriber", func(t *testing.T) {
		subscriberID, err := manager.Store(ctx, &User{
			ID:            primitive.NewObjectID(),
			Subscriptions: []*Friend{},
			Subscribers:   []*Friend{},
		})
		require.NoError(t, err)

		subscriptionID, err := manager.Store(ctx, &User{
			ID:            primitive.NewObjectID(),
			Subscriptions: []*Friend{},
			Subscribers:   []*Friend{},
		})
		require.NoError(t, err)

		err = manager.Subscribe(ctx, *subscriberID, *subscriptionID)
		require.NoError(t, err)

		relation, err := manager.GetRelations(ctx, *subscriberID, *subscriptionID)
		require.NoError(t, err)
		require.Equal(t, Subscription, relation)

		relation, err = manager.GetRelations(ctx, *subscriptionID, *subscriberID)
		require.NoError(t, err)
		require.Equal(t, Subscriber, relation)
	})
}
