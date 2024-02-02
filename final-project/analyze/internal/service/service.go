package service

import (
	"analyze/internal/indicator"
	"analyze/internal/repository"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"pkg/broker"
	"pkg/data"
)

type AnalyzeService struct {
	consumers              []broker.Consumer
	stockAnalyzeRepository repository.StockAnalysisRepo
}

func NewAnalyzeService(consumers []broker.Consumer, stockAnalyzeRepository repository.StockAnalysisRepo) *AnalyzeService {
	return &AnalyzeService{
		consumers:              consumers,
		stockAnalyzeRepository: stockAnalyzeRepository,
	}
}

func (s *AnalyzeService) Start() {
	messagesChannels := make([]chan *broker.Message, len(s.consumers))

	doneChan := make(chan bool)

	for i, consumer := range s.consumers {
		messagesChannels[i] = make(chan *broker.Message)
		go consumer.Consume(messagesChannels[i])
		go s.processMessages(messagesChannels[i])
	}

	<-doneChan
}

func (s *AnalyzeService) processMessages(messagesChan chan *broker.Message) {
	period := 3
	values := make([]data.StockData, 0)

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

		logrus.WithFields(logrus.Fields{
			"stock_symbol": stockData.StockSymbol,
			"timestamp":    stockData.Timestamp,
		}).Info("received stock data")

		values = append(values, stockData)

		if len(values) != period {
			continue
		}

		s.calculateIndicators(values, period)

		values = make([]data.StockData, 0)
	}
}

func (s *AnalyzeService) calculateIndicators(values []data.StockData, period int) {
	ma := indicator.CalculateMovingAverage(values, period)

	err := s.stockAnalyzeRepository.Create(&repository.StockAnalysis{
		Symbol:    values[0].StockSymbol,
		Indicator: "moving_average",
		Value:     ma,
		Timestamp: values[period-1].Timestamp,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to save moving average")
	}

	ema := indicator.CalculateEMA(values, period)

	err = s.stockAnalyzeRepository.Create(&repository.StockAnalysis{
		Symbol:    values[0].StockSymbol,
		Indicator: "exponential_moving_average",
		Value:     ema,
		Timestamp: values[period-1].Timestamp,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to save exponential moving average")
	}

	rsi := indicator.CalculateRSI(values, period)

	err = s.stockAnalyzeRepository.Create(&repository.StockAnalysis{
		Symbol:    values[0].StockSymbol,
		Indicator: "relative_strength_index",
		Value:     rsi,
		Timestamp: values[period-1].Timestamp,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to save relative strength index")
	}
}
