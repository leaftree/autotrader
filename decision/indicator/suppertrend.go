package indicator

import (
	"github.com/leaftree/autotrader/types"
	"math"
)

// 计算超级趋势指标
func CalculateSuperTrend(candles []types.Candle, period int, multiplier float64) []types.Indicators {
	indicators := make([]types.Indicators, len(candles))
	atr := CalculateATR(candles, period)

	for i := period; i < len(candles); i++ {
		hl2 := (candles[i].High + candles[i].Low) / 2
		basicUpper := hl2 + multiplier*atr[i]
		basicLower := hl2 - multiplier*atr[i]

		// 初始化
		if i == period {
			indicators[i] = types.Indicators{
				SuperTrendValue: basicUpper,
				SuperTrendTrend: "down",
				Price:           candles[i].Close,
				Timestamp:       candles[i].Timestamp,
			}
			continue
		}

		prev := indicators[i-1]

		// 计算当前超级趋势
		var currentValue float64
		var currentTrend string

		if prev.SuperTrendTrend == "down" {
			if candles[i].Close > prev.SuperTrendValue {
				currentTrend = "up"
				currentValue = basicLower
			} else {
				currentTrend = "down"
				currentValue = math.Max(basicUpper, prev.SuperTrendValue)
			}
		} else {
			if candles[i].Close < prev.SuperTrendValue {
				currentTrend = "down"
				currentValue = basicUpper
			} else {
				currentTrend = "up"
				currentValue = math.Min(basicLower, prev.SuperTrendValue)
			}
		}

		indicators[i] = types.Indicators{
			SuperTrendValue: currentValue,
			SuperTrendTrend: currentTrend,
			Price:           candles[i].Close,
			Timestamp:       candles[i].Timestamp,
		}
	}

	return indicators[period:]
}
