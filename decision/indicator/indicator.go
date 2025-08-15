package indicator

import (
	log "github.com/leaftree/autotrader/logger"
	"github.com/leaftree/autotrader/types"
)

var (
	logger = log.NewLogger("indicator")
)

type Indicator interface {
	Process([]types.Candle) types.DecisionType
}

// 组合所有指标计算
func AggIndicators(candles []types.Candle) []types.Indicators {
	// 指标参数
	superTrendPeriod := 10
	superTrendMultiplier := 3.0
	rsiPeriod := 14
	bollPeriod := 20
	bollStdDev := 2.0

	// 计算最大所需数据长度
	maxPeriod := superTrendPeriod
	if rsiPeriod > maxPeriod {
		maxPeriod = rsiPeriod
	}
	if bollPeriod > maxPeriod {
		maxPeriod = bollPeriod
	}

	// 确保有足够数据
	if len(candles) < maxPeriod*2 {
		logger.Debugf("candles: %v\n len(candels)=%d maxperiod*2=%d", candles, len(candles), maxPeriod*2)
		logger.Error("Insufficient data for indicator calculation")
		return nil
	}

	// 计算各项指标
	superTrend := CalculateSuperTrend(candles, superTrendPeriod, superTrendMultiplier)
	rsi := CalculateRSI(candles, rsiPeriod)
	bollUpper, bollMiddle, bollLower := CalculateBollingerBands(candles, bollPeriod, bollStdDev)

	// 组合指标结果
	startIndex := maxPeriod
	results := make([]types.Indicators, len(candles)-startIndex)
	for i := startIndex; i < len(candles); i++ {
		stIndex := i - superTrendPeriod
		results[i-startIndex] = types.Indicators{
			Timestamp:       candles[i].Timestamp,
			SuperTrendValue: superTrend[stIndex].SuperTrendValue,
			SuperTrendTrend: superTrend[stIndex].SuperTrendTrend,
			RSI:             rsi[i],
			BollUpper:       bollUpper[i],
			BollMiddle:      bollMiddle[i],
			BollLower:       bollLower[i],
			Price:           candles[i].Close,
		}
	}

	return results
}
