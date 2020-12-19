package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/lib/friend"
)

var cacheKey = "feed:user_id:%d"

type postMessage struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type daemon struct {
	logger   *zap.Logger
	redis    *redis.Client
	consumer *cluster.Consumer
	fm       *friend.Manager
}

func newDaemon(logger *zap.Logger, redis *redis.Client, consumer *cluster.Consumer, fm *friend.Manager) *daemon {
	return &daemon{logger: logger, redis: redis, consumer: consumer, fm: fm}
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

	subscribers, err := d.fm.GetFriends(context.Background(), m.UserID)
	if err != nil {
		d.logger.Error(
			"failed get subscribers",
			zap.Error(err),
		)
		return nil
	}

	for _, id := range subscribers {
		key := fmt.Sprintf(cacheKey, id)
		llen, err := d.redis.LLen(key).Result()
		if err != nil {
			d.logger.Error(
				"failed llen",
				zap.Error(err),
				zap.Int("user_id", id),
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
				zap.Int("user_id", id),
				zap.Int("post_id", m.ID),
			)
			return nil
		}
	}

	return nil
}
