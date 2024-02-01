package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	data3 "pkg/data"

	"github.com/sirupsen/logrus"

	"ingestion/internal/service"
)

type IngestionAPI struct {
	serviceModule *service.IngestionService
}

func NewIngestionAPI(serviceModule *service.IngestionService) *IngestionAPI {
	ingestAPI := &IngestionAPI{
		serviceModule: serviceModule,
	}

	http.HandleFunc("/ingest", ingestAPI.ingestHandler)

	return ingestAPI
}

func (a *IngestionAPI) ingestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dataType data3.FinancialDataType
	if err = json.Unmarshal(body, &dataType); err != nil {
		logrus.WithError(err).Error("Failed to decode data type")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data any

	switch dataType.Value {
	case data3.StockDataType:
		data, err = a.handleStockData(body)
		if err != nil {
			logrus.WithError(err).Error("Failed to handle stock data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case data3.OrderBookDataType:
		data, err = a.handleOrderBookData(body)
		if err != nil {
			logrus.WithError(err).Error("Failed to handle order book data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case data3.NewsSentimentDataType:
		data, err = a.handleNewsSentimentData(body)
		if err != nil {
			logrus.WithError(err).Error("Failed to handle news sentiment data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case data3.MarketDataType:
		data, err = a.handleMarketData(body)
		if err != nil {
			logrus.WithError(err).Error("Failed to handle market data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case data3.EconomicIndicatorDataType:
		data, err = a.handleEconomicIndicatorData(body)
		if err != nil {
			logrus.WithError(err).Error("Failed to handle economic indicator data")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		logrus.Error("Unknown data type")
		http.Error(w, "Unknown data type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received Stock Data: %+v\n", data)

	err = a.serviceModule.ProduceFinancialData(dataType.Value, data)
	if err != nil {
		logrus.WithError(err).Error("Failed to produce financial data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"data": data,
		"type": dataType.Value,
	}).Info("Received financial data")
}

func (a *IngestionAPI) handleStockData(body []byte) (any, error) {
	var data data3.StockData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *IngestionAPI) Start() {
	port := 8080
	fmt.Printf("Server running at port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
