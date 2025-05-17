package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
)

var cacheKey = "feed:user_id:%d"
var centrifugoKey = "feed_user_id_%d"

type postMessage struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type daemon struct {
	logger      *zap.Logger
	redis       *redis.Client
	consumer    sarama.ConsumerGroup
	usersClient usersGrpc.UsersClient
	centrifugo  *gocent.Client
}

func newDaemon(logger *zap.Logger, redis *redis.Client, consumer sarama.ConsumerGroup, usersClient usersGrpc.UsersClient, centrifugo *gocent.Client) *daemon {
	return &daemon{logger: logger, redis: redis, consumer: consumer, usersClient: usersClient, centrifugo: centrifugo}
}

func (d *daemon) run(ctx context.Context) error {
	topics := []string{"posts_feed"}
	for {
		err := d.consumer.Consume(ctx, topics, d)
		if err != nil {
			d.logger.Error("consume error", zap.Error(err))
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (d *daemon) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (d *daemon) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (d *daemon) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if err := d.handle(message); err != nil {
				d.logger.Error("cannot handle event", zap.Error(err))
				return err
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (d *daemon) handle(msg *sarama.ConsumerMessage) error {
	d.logger.Debug(
		"new message",
		zap.String("value", string(msg.Value)),
		zap.Int64("offset", msg.Offset),
		zap.Int32("partition", msg.Partition),
	)

	m := postMessage{}
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		d.logger.Error(
			"failed read message",
			zap.Error(err),
		)
		return nil
	}

	subscribersResponse, err := d.usersClient.GetFriendsIds(context.Background(), &usersGrpc.GetFriendsIdsRequest{UserId: int64(m.UserID)})
	if err != nil {
		d.logger.Error(
			"failed get subscribers",
			zap.Error(err),
		)
		return nil
	}

	for _, id := range subscribersResponse.UserIds {
		key := fmt.Sprintf(cacheKey, id)
		llen, err := d.redis.LLen(key).Result()
		if err != nil {
			d.logger.Error(
				"failed llen",
				zap.Error(err),
				zap.Int64("user_id", id),
				zap.Int("post_id", m.ID),
			)
			return nil
		}

		if llen >= 1000 {
			d.redis.RPop(key)
		}
		err = d.redis.LPush(key, msg.Value).Err()
		if err != nil {
			d.logger.Error(
				"failed set message to user's cache",
				zap.Error(err),
				zap.Int64("user_id", id),
				zap.Int("post_id", m.ID),
			)
			return nil
		}

		res, err := d.centrifugo.PresenceStats(context.Background(), fmt.Sprintf(centrifugoKey, id))
		if err != nil {
			d.logger.Error(
				"failed check presence centrifugo",
				zap.Error(err),
				zap.Int64("user_id", id),
				zap.Int("post_id", m.ID),
				zap.String("channel", fmt.Sprintf(centrifugoKey, id)),
				zap.String("msg", string(msg.Value)),
			)
			return nil
		}

		if res.NumClients == 0 {
			d.logger.Info("empty centrifugo clients")
			return nil
		}

		err = d.centrifugo.Publish(context.Background(), fmt.Sprintf(centrifugoKey, id), msg.Value)
		if err != nil {
			d.logger.Error(
				"failed set message to centrifugo",
				zap.Error(err),
				zap.Int64("user_id", id),
				zap.Int("post_id", m.ID),
				zap.String("channel", fmt.Sprintf(centrifugoKey, id)),
				zap.String("msg", string(msg.Value)),
			)
		}
	}

	return nil
}
