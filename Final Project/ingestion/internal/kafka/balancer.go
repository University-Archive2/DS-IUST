package kafka

import (
	"github.com/segmentio/kafka-go"
	"ingestion/internal/data"
)

type StockTypePartitionBalancer struct {
}

func (b *StockTypePartitionBalancer) Balance(msg kafka.Message, partitions ...int) int {
	partition := 0

	if msg.Key == nil {
		return partition
	}

	switch string(msg.Key) {
	case data.StockDataType:
		partition = 1
	case data.OrderBookDataType:
		partition = 2
	case data.NewsSentimentDataType:
		partition = 3
	case data.MarketDataType:
		partition = 4
	case data.EconomicIndicatorDataType:
		partition = 5
	default:
		partition = 6
	}

	return partition
}
