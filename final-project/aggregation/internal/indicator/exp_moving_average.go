package indicator

import (
	"pkg/data"
)

// CalculateEMA calculates the Exponential Moving Average for a slice of StockData
func CalculateEMA(data []data.StockData, windowSize int) float64 {
	k := 2.0 / (float64(windowSize) + 1.0) // Smoothing factor

	ema := CalculateMovingAverage(data, windowSize) // Initial EMA (SMA of first windowSize elements)

	for _, d := range data[1:] {
		ema = d.ClosingPrice*k + ema*(1-k) // equals to EMA = (Close - EMA(previous)) * k + EMA(previous)
	}

	return ema
}
