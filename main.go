package main

import (
	"context"
	"os"
	"time"

	"github.com/leaftree/autotrader/config"
	"github.com/leaftree/autotrader/decision"
	"github.com/leaftree/autotrader/decision/risk"
	log "github.com/leaftree/autotrader/logger"
	tr "github.com/leaftree/autotrader/trader"
	"github.com/leaftree/autotrader/types"
)

// TODO 多个逻辑需要使用 K line 做分析，结构修改成任务式，独立进程获取 K 线，异步通知分析逻辑

var (
	trader tr.Trader
	logger = log.NewLoggerW("main", os.Stdout, os.Stdout, os.Stderr, os.Stderr)
)

func init() {
	// TODO init config
	config.Init()
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
	MaxLines   = 100  // K线数量
	KLineSize  = 50
)

// 开多单
func openLongPosition(ctx context.Context, price float64) {
	// 设置订单参数
	size, _ := trader.CalcPositionSize(ctx) // 计算仓位大小

	err := trader.CreateOrder(ctx, &types.Order{
		Contract: CONTRACT,
		Size:     size,
		Price:    price,
		Side:     types.SideTypeLong,
	})

	if err != nil {
		logger.Errorf("Error opening long position: %v", err)
		return
	}
	logger.Infof("Opened long position at %.2f, size: %d", price, size)
}

// 开空单
func openShortPosition(ctx context.Context, price float64) {

	// 设置订单参数 (负值表示空单)
	size, _ := trader.CalcPositionSize(ctx)

	// 发送订单
	err := trader.CreateOrder(ctx, &types.Order{
		Contract: CONTRACT,
		Size:     size,
		Price:    price,
		Side:     types.SideTypeShort,
	})
	if err != nil {
		logger.Errorf("Error opening short position: %v", err)
		return
	}
	logger.Infof("Opened short position at %.2f, size: %d", price, size)
}

// 计算仓位大小 (简化版)
func calculatePositionSize(price float64) int64 {
	// 实际应用中应根据账户余额和风险模型计算
	// 这里固定为1个合约
	return 1
}

func mainLoop() {
	ctx := context.Background()
	for {
		// 获取最新100根K线
		candles, err := trader.QueryCandles(ctx, CONTRACT, INTERVAL, KLineSize)
		if err != nil {
			logger.Errorf("Error fetching candles: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		dec := decision.MakeTradingDecision(ctx, candles)
		price := candles[len(candles)-1].Low

		switch dec {
		case types.DecisionLong:
			if has, _ := risk.HasPosition(ctx, CONTRACT); has {
				continue
			}
			trader.CreateOrder(ctx, &types.Order{Contract: CONTRACT, Price: price, Side: types.SideTypeLong, Size: 1})
			risk.PositionType = types.SideTypeLong
		case types.DecisionShort:
			if has, _ := risk.HasPosition(ctx, CONTRACT); has {
				continue
			}
			trader.CreateOrder(ctx, &types.Order{Contract: CONTRACT, Price: price, Side: types.SideTypeShort, Size: 1})
			risk.PositionType = types.SideTypeShort
		case types.DecisionClose:
			trader.ClosePosition(ctx)
			risk.PositionType = types.SideTypeNone
		default:
			logger.Info("decision nothing todo")
		}

		// 等待5second(与K线间隔匹配)
		time.Sleep(5 * time.Second)
	}
}

// XXX
func sendNotify(ctx context.Context) {
	trader.CreateOrder(ctx, &types.Order{Price: 100, Side: types.SideTypeLong, Size: 1})
}

// 主循环
func main() {
	// 设置日志
	//logFile, err := os.OpenFile("eth_trading.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatal("Failed to open log file:", err)
	//}
	//defer logFile.Close()
	//log.SetOutput(logFile)

	logger.Info("Starting ETH trading bot...")

	// 主循环
	mainLoop()
	//sendNotify(context.Background())

	time.Sleep(2 * time.Second)
}
