package redisprovider

import "github.com/go-redis/redis"

// New creates redis client
func New(redisURL string) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(redisOpts), nil
}
