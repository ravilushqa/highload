package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
)

var (
	cacheKey      = "feed:user_id:%s"
	centrifugoKey = "feed_user_id_%s"
)

type postMessage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type daemon struct {
	logger      *zap.Logger
	redis       *redis.Client
	consumer    *cluster.Consumer
	usersClient usersGrpc.UsersClient
	centrifugo  *gocent.Client
}

func newDaemon(logger *zap.Logger, redis *redis.Client, consumer *cluster.Consumer, usersClient usersGrpc.UsersClient, centrifugo *gocent.Client) *daemon {
	return &daemon{logger: logger, redis: redis, consumer: consumer, usersClient: usersClient, centrifugo: centrifugo}
}

func (d *daemon) run(ctx context.Context) error {
	nackErrors := make(chan error, 1)
	for {
		select {
		case part, ok := <-d.consumer.Partitions():
			if !ok {
				return nil
			}
			d.logger.Info(fmt.Sprintf("start listening %s:%d", part.Topic(), part.Partition()))
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					if err := d.handle(msg); err != nil {
						d.logger.Error("cannot handle event", zap.Error(err))
						nackErrors <- err
						return
					}
					d.consumer.MarkOffset(msg, "")
				}
			}(part)
		case err := <-d.consumer.Errors():
			d.logger.Error("listen error", zap.Error(err))
		case ntf := <-d.consumer.Notifications():
			d.logger.Debug("consumer rebalanced", zap.String("notification", fmt.Sprintf("%+v", ntf)))
		case nackError := <-nackErrors:
			return nackError
		case <-ctx.Done():
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

	subscribersResponse, err := d.usersClient.GetFriendsIds(context.Background(), &usersGrpc.GetFriendsIdsRequest{UserId: m.UserID})
	if err != nil {
		d.logger.Error(
			"failed get subscribers",
			zap.Error(err),
		)
		return nil
	}

	d.logger.Info("subscribers", zap.Any("subscribers", subscribersResponse.UserIds))

	for _, id := range subscribersResponse.UserIds {
		key := fmt.Sprintf(cacheKey, id)
		llen, err := d.redis.LLen(key).Result()
		if err != nil {
			d.logger.Error(
				"failed llen",
				zap.Error(err),
				zap.String("user_id", id),
				zap.String("post_id", m.ID),
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
				zap.String("user_id", id),
				zap.String("post_id", m.ID),
			)
			return nil
		}

		d.logger.Info("set message to centrifugo", zap.String("feed_user_id", fmt.Sprintf(centrifugoKey, id)))
		res, err := d.centrifugo.PresenceStats(context.Background(), fmt.Sprintf(centrifugoKey, id))
		if err != nil {
			d.logger.Error(
				"failed check presence centrifugo",
				zap.Error(err),
				zap.String("user_id", id),
				zap.String("post_id", m.ID),
				zap.String("channel", fmt.Sprintf(centrifugoKey, id)),
				zap.String("msg", string(msg.Value)),
			)
			return nil
		}

		if res.NumClients == 0 {
			d.logger.Info("empty centrifugo clients")
			continue
		}

		err = d.centrifugo.Publish(context.Background(), fmt.Sprintf(centrifugoKey, id), msg.Value)
		if err != nil {
			d.logger.Error(
				"failed set message to centrifugo",
				zap.Error(err),
				zap.String("user_id", id),
				zap.String("post_id", m.ID),
				zap.String("channel", fmt.Sprintf(centrifugoKey, id)),
				zap.String("msg", string(msg.Value)),
			)
		}
	}

	return nil
}
