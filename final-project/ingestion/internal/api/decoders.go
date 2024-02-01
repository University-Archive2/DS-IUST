package api

import (
	"encoding/json"
	data2 "pkg/data"
)

func (a *IngestionAPI) handleOrderBookData(body []byte) (any, error) {
	var data data2.OrderBookData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *IngestionAPI) handleNewsSentimentData(body []byte) (any, error) {
	var data data2.NewsSentimentData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *IngestionAPI) handleMarketData(body []byte) (any, error) {
	var data data2.MarketData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *IngestionAPI) handleEconomicIndicatorData(body []byte) (any, error) {
	var data data2.EconomicIndicatorData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
