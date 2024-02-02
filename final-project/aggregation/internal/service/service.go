package service

import (
	"aggregation/internal/cache"
	"aggregation/internal/repository"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"pkg/broker"
	"pkg/data"
	"time"
)

type AggregationService struct {
	consumers             []broker.Consumer
	metricsCache          *cache.MetricsCache
	aggregatedMetricsRepo repository.AggregatedMetricsRepo
}

func NewAggregationService(consumers []broker.Consumer, aggregatedMetricsRepo repository.AggregatedMetricsRepo, metricsCache *cache.MetricsCache) *AggregationService {
	return &AggregationService{
		consumers:             consumers,
		aggregatedMetricsRepo: aggregatedMetricsRepo,
		metricsCache:          metricsCache,
	}
}

func (s *AggregationService) Start() {
	messagesChannels := make([]chan *broker.Message, len(s.consumers))

	doneChan := make(chan bool)

	for i, consumer := range s.consumers {
		messagesChannels[i] = make(chan *broker.Message)
		go consumer.Consume(messagesChannels[i])
		go s.processMessages(messagesChannels[i])
	}

	<-doneChan
}

func (s *AggregationService) processMessages(messagesChan chan *broker.Message) {
	count := 0

	totalVolume := 0.0
	totalOpeningPrice := 0.0
	totalClosingPrice := 0.0

	highestPrice := 0.0
	lowestPrice := 0.0

	for message := range messagesChan {
		logrus.WithFields(logrus.Fields{
			"key": message.Key,
		}).Info("received message")

		var stockData data.StockData
		err := json.Unmarshal(message.Value, &stockData)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal message from kafka")
			continue
		}

		if count == 0 {
			count, totalVolume, totalOpeningPrice, totalClosingPrice, highestPrice, lowestPrice = s.getMiddleValues(stockData)
		}

		count++

		logrus.WithFields(logrus.Fields{
			"stock_symbol": stockData.StockSymbol,
			"timestamp":    time.Unix(int64(stockData.Timestamp), 0),
		}).Info("received stock data")

		totalVolume += stockData.Volume

		totalOpeningPrice += stockData.OpeningPrice

		totalClosingPrice += stockData.ClosingPrice

		if stockData.High > highestPrice {
			highestPrice = stockData.High
		}

		if stockData.Low < lowestPrice {
			lowestPrice = stockData.Low
		}

		go s.persistMiddleValues(stockData.StockSymbol, count, totalVolume, totalOpeningPrice, totalClosingPrice, highestPrice, lowestPrice)

		if count%10 == 0 {
			s.StoreAggregatedMetrics(stockData.StockSymbol, count, totalVolume, totalOpeningPrice, totalClosingPrice, highestPrice, lowestPrice)
		}
	}
}

func (s *AggregationService) getMiddleValues(stockData data.StockData) (count int, totalVolume float64, totalOpeningPrice float64, totalClosingPrice float64, highestPrice float64, lowestPrice float64) {
	countData, err := s.metricsCache.Get("count" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get count from cache")
	}

	if countData != nil {
		err = json.Unmarshal(countData, &count)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal count from cache")
		}
	}

	totalVolumeData, err := s.metricsCache.Get("total_volume" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get total_volume from cache")
	}

	if totalVolumeData != nil {
		err = json.Unmarshal(totalVolumeData, &totalVolume)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal total_volume from cache")
		}
	}

	totalOpeningPriceData, err := s.metricsCache.Get("total_opening_price" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get total_opening_price from cache")
	}

	if totalOpeningPriceData != nil {
		err = json.Unmarshal(totalOpeningPriceData, &totalOpeningPrice)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal total_opening_price from cache")
		}
	}

	totalClosingPriceData, err := s.metricsCache.Get("total_closing_price" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get total_closing_price from cache")
	}

	if totalClosingPriceData != nil {
		err = json.Unmarshal(totalClosingPriceData, &totalClosingPrice)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal total_closing_price from cache")
		}
	}

	highestPriceData, err := s.metricsCache.Get("highest_price" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get highest_price from cache")
	}

	if highestPriceData != nil {
		err = json.Unmarshal(highestPriceData, &highestPrice)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal highest_price from cache")
		}
	}

	lowestPriceData, err := s.metricsCache.Get("lowest_price" + stockData.StockSymbol)
	if err != nil {
		logrus.WithError(err).Error("failed to get lowest_price from cache")
	}

	if lowestPriceData != nil {
		err = json.Unmarshal(lowestPriceData, &lowestPrice)
		if err != nil {
			logrus.WithError(err).Error("failed to unmarshal lowest_price from cache")
		}
	}

	return
}

func (s *AggregationService) persistMiddleValues(stockSymbol string, count int, totalVolume float64, totalOpeningPrice float64, totalClosingPrice float64, highestPrice float64, lowestPrice float64) {
	countData, err := json.Marshal(count)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal count")
	}

	err = s.metricsCache.Set("count"+stockSymbol, countData)
	if err != nil {
		logrus.WithError(err).Error("failed to set count in cache")
	}

	totalVolumeData, err := json.Marshal(totalVolume)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal total_volume")
	}

	err = s.metricsCache.Set("total_volume"+stockSymbol, totalVolumeData)
	if err != nil {
		logrus.WithError(err).Error("failed to set total_volume in cache")
	}

	totalOpeningPriceData, err := json.Marshal(totalOpeningPrice)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal total_opening_price")
	}

	err = s.metricsCache.Set("total_opening_price"+stockSymbol, totalOpeningPriceData)
	if err != nil {
		logrus.WithError(err).Error("failed to set total_opening_price in cache")
	}

	totalClosingPriceData, err := json.Marshal(totalClosingPrice)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal total_closing_price")
	}

	err = s.metricsCache.Set("total_closing_price"+stockSymbol, totalClosingPriceData)
	if err != nil {
		logrus.WithError(err).Error("failed to set total_closing_price in cache")
	}

	highestPriceData, err := json.Marshal(highestPrice)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal highest_price")
	}

	err = s.metricsCache.Set("highest_price"+stockSymbol, highestPriceData)
	if err != nil {
		logrus.WithError(err).Error("failed to set highest_price in cache")
	}

	lowestPriceData, err := json.Marshal(lowestPrice)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal lowest_price")
	}

	err = s.metricsCache.Set("lowest_price"+stockSymbol, lowestPriceData)
	if err != nil {
		logrus.WithError(err).Error("failed to set lowest_price in cache")
	}
}

func (s *AggregationService) StoreAggregatedMetrics(stockSymbol string, count int, totalVolume float64, totalOpeningPrice float64, totalClosingPrice float64, highestPrice float64, lowestPrice float64) {
	metrics := []*repository.AggregatedMetric{
		{
			Symbol: stockSymbol,
			Metric: "total_volume",
			Value:  totalVolume,
		},
		{
			Symbol: stockSymbol,
			Metric: "average_volume",
			Value:  totalVolume / float64(count),
		},
		{
			Symbol: stockSymbol,
			Metric: "average_opening_price",
			Value:  totalOpeningPrice / float64(count),
		},
		{
			Symbol: stockSymbol,
			Metric: "average_closing_price",
			Value:  totalClosingPrice / float64(count),
		},
		{
			Symbol: stockSymbol,
			Metric: "highest_price",
			Value:  highestPrice,
		},
		{
			Symbol: stockSymbol,
			Metric: "lowest_price",
			Value:  lowestPrice,
		},
	}

	err := s.aggregatedMetricsRepo.BatchCreate(metrics)
	if err != nil {
		logrus.WithError(err).Error("failed to create aggregated metrics")
	}
}
