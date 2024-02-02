package indicator

import (
	"pkg/data"
)

// CalculateMovingAverage calculates the moving average for a slice of StockData
func CalculateMovingAverage(data []data.StockData, windowSize int) float64 {
	sum := 0.0

	for _, d := range data[len(data)-windowSize:] {
		sum += d.ClosingPrice
	}

	return sum / float64(windowSize)
}
