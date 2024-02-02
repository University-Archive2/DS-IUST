package indicator

import (
	"pkg/data"
)

// CalculateRSI calculates the Relative Strength Index for a slice of StockData
func CalculateRSI(data []data.StockData, period int) float64 {
	gains, losses := 0.0, 0.0

	for i := 1; i < len(data); i++ {
		change := data[i].ClosingPrice - data[i-1].ClosingPrice
		if change > 0 {
			gains += change
		} else {
			losses -= change // Convert to positive
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100 // Prevent division by zero; implies all gains, no losses
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}
