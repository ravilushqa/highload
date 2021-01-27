package saga

import (
	"time"

	"github.com/lysu/go-saga"
)

func New(zkAddrs, brokerAddrs []string) *saga.ExecutionCoordinator {
	saga.StorageConfig.Kafka.ZkAddrs = zkAddrs
	saga.StorageConfig.Kafka.BrokerAddrs = brokerAddrs
	saga.StorageConfig.Kafka.Partitions = 1
	saga.StorageConfig.Kafka.Replicas = 1
	saga.StorageConfig.Kafka.ReturnDuration = 50 * time.Millisecond
	ec := saga.NewSEC()

	return &ec
}
