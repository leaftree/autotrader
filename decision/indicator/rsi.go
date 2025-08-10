package indicator

import (
	"github.com/leaftree/autotrader/types"
	"math"
)

// 计算RSI指标
func CalculateRSI(candles []types.Candle, period int) []float64 {
	rsi := make([]float64, len(candles))
	gains := make([]float64, len(candles))
	losses := make([]float64, len(candles))

	// 计算价格变化
	for i := 1; i < len(candles); i++ {
		change := candles[i].Close - candles[i-1].Close
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = math.Abs(change)
		}
	}

	// 计算初始平均增益和平均损失
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// 计算后续RSI值
	for i := period + 1; i < len(candles); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}
