package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/leaftree/autotrader/decision"
	"github.com/leaftree/autotrader/decision/indicator"
	tr "github.com/leaftree/autotrader/trader"
	"github.com/leaftree/autotrader/types"
)

var trader tr.Trader

func init() {
	// TODO init config
	// TODO init logger
	trader = tr.NewTrader("gateio")
}

func finally() {
	// TODO close log
}

// 配置参数
const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
	SYMBOL     = "ETH_USDT"
	CURRENCY   = "ETH"
	SETTLE     = "usdt" // 结算货币
	CONTRACT   = "ETH_USDT"
	INTERVAL   = "5m" // K线间隔
)

// 开多单
func openLongPosition(ctx context.Context, price float64) {
	// 设置订单参数
	size := calculatePositionSize(price) // 计算仓位大小

	err := trader.CreateOrder(ctx, &types.Order{
		Contract: CONTRACT,
		Size:     size,
		Price:    price,
		Side:     types.SideTypeLong,
	})

	if err != nil {
		log.Printf("Error opening long position: %v", err)
		return
	}
	log.Printf("Opened long position at %.2f, size: %d", price, size)
}

// 开空单
func openShortPosition(ctx context.Context, price float64) {

	// 设置订单参数 (负值表示空单)
	size := calculatePositionSize(price)

	// 发送订单
	err := trader.CreateOrder(ctx, &types.Order{
		Contract: CONTRACT,
		Size:     size,
		Price:    price,
		Side:     types.SideTypeShort,
	})
	if err != nil {
		log.Printf("Error opening short position: %v", err)
		return
	}
	log.Printf("Opened short position at %.2f, size: %d", price, size)
}

// 计算仓位大小 (简化版)
func calculatePositionSize(price float64) int64 {
	// 实际应用中应根据账户余额和风险模型计算
	// 这里固定为1个合约
	return 1
}

// 交易决策函数
func makeTradingDecision(ctx context.Context, indicators []types.Indicators) {
	if len(indicators) < 2 {
		log.Println("Not enough indicators for decision")
		return
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
		openLongPosition(ctx, current.Price)
		return
	}

	// 空单开仓条件
	if current.SuperTrendTrend == "down" && // 超级趋势向下
		current.RSI < 70 && prev.RSI >= 70 && // RSI从超买区回落
		current.Price < current.BollUpper { // 价格跌破布林上轨
		log.Println("SHORT signal detected")
		openShortPosition(ctx, current.Price)
		return
	}

	log.Println("No trading signal detected")
}

// 主循环
func main() {
	ctx := context.Background()
	// 设置日志
	logFile, err := os.OpenFile("eth_trading.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Starting ETH trading bot...")

	// 主循环
	for {
		// 获取最新100根K线
		candles, err := trader.QueryCandles(ctx, "ETH_USDT", 100)
		if err != nil {
			log.Printf("Error fetching candles: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		// 计算指标
		indicators := indicator.AggIndicators(candles)

		// 执行交易决策
		dec := decision.MakeTradingDecision(ctx, indicators)

		switch dec {
		case types.DecisionLong:
			trader.CreateOrder(ctx, &types.Order{})
		case types.DecisionShort:
			trader.CreateOrder(ctx, &types.Order{})
		case types.DecisionClose:
			trader.ClosePosition(ctx)
		default:
			log.Printf("decision nothing todo")
		}

		// 等待5second(与K线间隔匹配)
		time.Sleep(5 * time.Second)
	}
}
