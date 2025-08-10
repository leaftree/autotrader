package indicator

import (
	"math"

	"github.com/leaftree/autotrader/types"
)

// 计算ATR (平均真实波幅)
func CalculateATR(candles []types.Candle, period int) []float64 {
	tr := make([]float64, len(candles))
	atr := make([]float64, len(candles))

	for i := 1; i < len(candles); i++ {
		highLow := candles[i].High - candles[i].Low
		highPrevClose := math.Abs(candles[i].High - candles[i-1].Close)
		lowPrevClose := math.Abs(candles[i].Low - candles[i-1].Close)
		tr[i] = math.Max(highLow, math.Max(highPrevClose, lowPrevClose))
	}

	// 第一个ATR是前period个TR的平均值
	var sum float64
	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	// 后续ATR使用公式: ATR = (前一个ATR * (n-1) + 当前TR) / n
	for i := period + 1; i < len(candles); i++ {
		atr[i] = (atr[i-1]*float64(period-1) + tr[i]) / float64(period)
	}

	return atr
}
