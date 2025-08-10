package decision

import (
	"context"
	"log"

	"github.com/leaftree/autotrader/types"
)

// 交易决策函数
func MakeTradingDecision(ctx context.Context, indicators []types.Indicators) types.DecisionType {
	if len(indicators) < 2 {
		log.Println("Not enough indicators for decision")
		return types.DecisionNothing
	}

	current := indicators[len(indicators)-1]
	prev := indicators[len(indicators)-2]

	log.Printf("Current Indicators: Trend=%s, RSI=%.2f, Price=%.2f, BollUpper=%.2f, BollLower=%.2f",
		current.SuperTrendTrend, current.RSI, current.Price, current.BollUpper, current.BollLower)

	// 多单开仓条件
	if current.SuperTrendTrend == "up" && // 超级趋势向上
		current.RSI > 30 && prev.RSI <= 30 && // RSI从超卖区回升
		current.Price > current.BollLower { // 价格突破布林下轨
		log.Println("LONG signal detected")
		return types.DecisionLong
	}

	// 空单开仓条件
	if current.SuperTrendTrend == "down" && // 超级趋势向下
		current.RSI < 70 && prev.RSI >= 70 && // RSI从超买区回落
		current.Price < current.BollUpper { // 价格跌破布林上轨
		log.Println("SHORT signal detected")
		return types.DecisionShort
	}

	log.Println("No trading signal detected")
	return types.DecisionNothing
}
