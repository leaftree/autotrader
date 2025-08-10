package indicator

import (
	"github.com/leaftree/autotrader/types"
	"math"
)

// 计算布林带
func CalculateBollingerBands(candles []types.Candle, period int, stdDev float64) ([]float64, []float64, []float64) {
	middle := make([]float64, len(candles))
	upper := make([]float64, len(candles))
	lower := make([]float64, len(candles))

	for i := period - 1; i < len(candles); i++ {
		// 计算SMA
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += candles[j].Close
		}
		sma := sum / float64(period)
		middle[i] = sma

		// 计算标准差
		var variance float64
		for j := i - period + 1; j <= i; j++ {
			diff := candles[j].Close - sma
			variance += diff * diff
		}
		std := math.Sqrt(variance / float64(period))

		upper[i] = sma + stdDev*std
		lower[i] = sma - stdDev*std
	}

	return upper, middle, lower
}
