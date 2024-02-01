package data

type FinancialDataType struct {
	Value string `json:"data_type"`
}

type StockData struct {
	StockSymbol  string  `json:"stock_symbol"`
	OpeningPrice float64 `json:"opening_price"`
	ClosingPrice float64 `json:"closing_price"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Volume       float64 `json:"volume"`
	Timestamp    float64 `json:"timestamp"`
}

type OrderBookData struct {
	Timestamp   float64 `json:"timestamp"`
	StockSymbol string  `json:"stock_symbol"`
	OrderType   string  `json:"order_type"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
}

type NewsSentimentData struct {
	Timestamp          float64 `json:"timestamp"`
	StockSymbol        string  `json:"stock_symbol"`
	SentimentScore     float64 `json:"sentiment_score"`
	SentimentMagnitude float64 `json:"sentiment_magnitude"`
}

type MarketData struct {
	Timestamp   float64 `json:"timestamp"`
	StockSymbol string  `json:"stock_symbol"`
	MarketCap   float64 `json:"market_cap"`
	PERatio     float64 `json:"pe_ratio"`
}

type EconomicIndicatorData struct {
	Timestamp     float64 `json:"timestamp"`
	IndicatorName string  `json:"indicator_name"`
	Value         float64 `json:"value"`
}
