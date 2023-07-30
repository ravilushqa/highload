package friend

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define your Manager type and GetRelation function here (copy the code provided in the question).

func TestGetRelation(t *testing.T) {
	// Start a MongoDB container for testing.
	mongoC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer mongoC.Terminate(context.Background())

	// Get the MongoDB connection string from the container.
	endpoint, err := mongoC.Endpoint(context.Background(), "mongodb")
	if err != nil {
		t.Error(fmt.Errorf("failed to get endpoint: %w", err))
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatal(fmt.Errorf("error creating mongo client: %w", err))
	}
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	db := client.Database("test-db")
	collection := db.Collection("test-collection")

	// Insert test data into the collection.
	testData := []interface{}{
		bson.M{"user_id": "user1", "friend_id": "user2", "approved": true},
		bson.M{"user_id": "user2", "friend_id": "user1", "approved": true},
		bson.M{"user_id": "user3", "friend_id": "user1", "approved": false},
	}
	_, err = collection.InsertMany(context.Background(), testData)
	require.NoError(t, err)

	// Create an instance of your Manager type.
	manager := &Manager{
		col: collection,
	}

	// Test cases for GetRelation function.
	testCases := []struct {
		authUser string
		user     string
		expected Status
	}{
		{"user1", "user2", Friends},
		{"user1", "user3", Added},
		{"user2", "user1", Friends},
		{"user3", "user1", Requested},
		{"user2", "user3", ""},
		{"non-existent-user", "user1", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.authUser+"-"+tc.user, func(t *testing.T) {
			status, err := manager.GetRelation(context.Background(), tc.authUser, tc.user)
			require.NoError(t, err)
			require.Equal(t, tc.expected, status)
		})
	}
}
