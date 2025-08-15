package decision

import (
	"context"

	"github.com/leaftree/autotrader/config"
	"github.com/leaftree/autotrader/decision/indicator"
	log "github.com/leaftree/autotrader/logger"
	"github.com/leaftree/autotrader/types"
)

var (
	logger = log.NewLogger("decision")
)

// supper trend strategy
func SuppertrendStrategy(ctx context.Context, candles []types.Candle) types.DecisionType {
	// 计算指标
	indicators := indicator.AggIndicators(candles)

	if len(indicators) < 2 {
		logger.Warning("Not enough indicators for decision")
		return types.DecisionNone
	}

	current := indicators[len(indicators)-1]

	// 多单开仓条件
	if current.SuperTrendTrend == "up" { // 超级趋势向上
		logger.Info("LONG signal detected")
		return types.DecisionLong
	}
	logger.Debug("SHORT signal detected")
	return types.DecisionShort
}

// RMI trend sniper
func RMIStrategy(ctx context.Context, candles []types.Candle) types.DecisionType {
	in := indicator.NewRmiTrendSniper(indicator.Config{
		Length:    14,
		PMom:      66,
		NMon:      30,
		BandMulti: 0.3,
		Window:    20,
	})
	return in.Process(candles)
}

func AggStrategy(ctx context.Context, candles []types.Candle) types.DecisionType {
	// 计算指标
	indicators := indicator.AggIndicators(candles)

	if len(indicators) < 2 {
		logger.Warning("Not enough indicators for decision")
		return types.DecisionNone
	}

	current := indicators[len(indicators)-1]
	prev := indicators[len(indicators)-2]

	logger.Infof("Current Indicators: Trend=%s, RSI=%.2f, Price=%.2f, BollUpper=%.2f, BollLower=%.2f",
		current.SuperTrendTrend, current.RSI, current.Price, current.BollUpper, current.BollLower)

	// 多单开仓条件
	if current.SuperTrendTrend == "up" && // 超级趋势向上
		current.RSI > 30 && prev.RSI <= 30 && // RSI从超卖区回升
		current.Price > current.BollLower { // 价格突破布林下轨
		logger.Info("LONG signal detected")
		return types.DecisionLong
	}

	// 空单开仓条件
	if current.SuperTrendTrend == "down" && // 超级趋势向下
		current.RSI < 70 && prev.RSI >= 70 && // RSI从超买区回落
		current.Price < current.BollUpper { // 价格跌破布林上轨
		logger.Info("SHORT signal detected")
		return types.DecisionShort
	}

	logger.Debug("No trading signal detected")
	return types.DecisionNone
}

// 交易决策函数
func MakeTradingDecision(ctx context.Context, candles []types.Candle) types.DecisionType {
	strategy := config.GetConfig().Strategy
	if strategy.RMI {
		logger.Debug("use rmi strategy")
		return RMIStrategy(ctx, candles)
	}
	if strategy.Suppertrend {
		logger.Debug("use suppertrend strategy")
		return SuppertrendStrategy(ctx, candles)
	}
	logger.Debug("use multi strategy agg")
	return AggStrategy(ctx, candles)
}
