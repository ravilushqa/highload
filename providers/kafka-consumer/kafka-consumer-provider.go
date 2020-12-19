package kafkaconsumerprovider

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

// New creates kafka consumer
func New(addrs []string, groupID string, topics []string, kfkCfg *cluster.Config) (*cluster.Consumer, error) {
	if kfkCfg == nil {
		kfkCfg = cluster.NewConfig()
		kfkCfg.Consumer.Return.Errors = true
		kfkCfg.Group.Return.Notifications = true
		kfkCfg.Group.Mode = cluster.ConsumerModePartitions
		kfkCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	return cluster.NewConsumer(addrs, groupID, topics, kfkCfg)
}
