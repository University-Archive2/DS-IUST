package kafka

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"pkg/data"
)

type StockTypePartitionBalancer struct {
}

func (b *StockTypePartitionBalancer) Balance(msg kafka.Message, partitions ...int) int {
	partition := 0

	if msg.Key == nil || string(msg.Key) != data.StockDataType {
		return partition
	}

	var stockType data.StockData
	err := json.Unmarshal(msg.Value, &stockType)
	if err != nil {
		logrus.WithError(err).Error("failed to unmarshal message in partition balancer")
		return partition
	}

	partition, exists := StockDataSymbolToPartition[stockType.StockSymbol]
	if exists {
		return partition
	}

	partition = 5

	return partition
}

var StockDataSymbolToPartition = map[string]int{
	"AAPL":  0,
	"GOOGL": 1,
	"AMZN":  2,
	"MSFT":  3,
	"TSLA":  4,
}
