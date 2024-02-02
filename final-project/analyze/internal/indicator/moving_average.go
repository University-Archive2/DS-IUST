package indicator

import (
	"pkg/data"
)

// CalculateEMA calculates the Exponential Moving Average for a slice of StockData
func CalculateEMA(data []data.StockData, period int) float64 {
	alpha := 2.0 / (float64(period) + 1.0)
	ema := data[0].ClosingPrice // Starting with the first data point as initial EMA

	for _, d := range data[1:] {
		ema = (d.ClosingPrice-alpha)*ema + alpha*d.ClosingPrice
	}

	return ema
}
