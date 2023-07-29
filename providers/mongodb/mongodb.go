package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ProviderOption func(clientOptions *options.ClientOptions) *options.ClientOptions

// New creates mongodb connection which is ready for use.
func New(ctx context.Context, url, dbName string, providerOptions ...ProviderOption) (*mongo.Database, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := options.Client()

	for _, providerOption := range providerOptions {
		opts = providerOption(opts)
	}

	client, err := mongo.Connect(ctx, opts.ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client.Database(dbName), nil
}

func WithMonitor(monitor *event.CommandMonitor) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetMonitor(monitor)
	}
}

func WithMaxPoolSize(size uint64) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetMaxPoolSize(size)
	}
}

func WithConnectTimeout(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetConnectTimeout(timeout)
	}
}

func WithSocketTimeout(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetSocketTimeout(timeout)
	}
}

func WithServerSelectionTimeout(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetServerSelectionTimeout(timeout)
	}
}

func WithMaxConnIdleTime(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetMaxConnIdleTime(timeout)
	}
}

func WithHeartbeatInterval(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetHeartbeatInterval(timeout)
	}
}

func WithLocalThreshold(timeout time.Duration) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetLocalThreshold(timeout)
	}
}

func WithReplicaSet(replicaSet string) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetReplicaSet(replicaSet)
	}
}

func WithReadConcern(readConcern *readconcern.ReadConcern) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetReadConcern(readConcern)
	}
}

func WithReadPreference(readPreference *readpref.ReadPref) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetReadPreference(readPreference)
	}
}

func WithWriteConcern(writeConcern *writeconcern.WriteConcern) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetWriteConcern(writeConcern)
	}
}

func WithAuth(auth options.Credential) ProviderOption {
	return func(clientOptions *options.ClientOptions) *options.ClientOptions {
		return clientOptions.SetAuth(auth)
	}
}
