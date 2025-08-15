package trader

import (
	"context"
	"fmt"
	"time"

	"github.com/antihax/optional"
	gateapi "github.com/gateio/gateapi-go/v6"
	"github.com/leaftree/autotrader/config"
	log "github.com/leaftree/autotrader/logger"
	"github.com/leaftree/autotrader/types"
	"github.com/leaftree/autotrader/util"
)

const (
	settle = "usdt"
)

var (
	logger = log.NewLogger("gateio")
)

type gateioTrader struct {
	client *gateapi.APIClient
}

func NewGateIOTrader() Trader {
	return newGateIOTrader()
}

func newGateIOTrader() *gateioTrader {
	cfg := gateapi.NewConfiguration()
	cfg.BasePath = "https://api.gateio.ws/api/v4"
	//cfg.Debug = true
	cfg.Key = config.GetConfig().GateIO.Key
	cfg.Secret = config.GetConfig().GateIO.Secret
	client := gateapi.NewAPIClient(cfg)
	return &gateioTrader{
		client: client,
	}
}

// 获取K线数据
func (gt *gateioTrader) QueryCandles(ctx context.Context, contract string, timeframe string, limit int) ([]types.Candle, error) {
	opts := gateapi.ListFuturesCandlesticksOpts{
		Limit:    optional.NewInt32(int32(limit)),
		Interval: optional.NewString(timeframe),
	}

	// 获取K线数据
	candles, _, err := gt.client.FuturesApi.ListFuturesCandlesticks(ctx, settle, contract, &opts)
	if err != nil {
		return nil, err
	}

	result := make([]types.Candle, len(candles))
	for i, c := range candles {
		timestamp := time.Unix(int64(c.T), 0)
		result[i] = types.Candle{
			Timestamp: timestamp,
			Open:      util.Str2Float64(c.O),
			High:      util.Str2Float64(c.H),
			Low:       util.Str2Float64(c.L),
			Close:     util.Str2Float64(c.C),
			Volume:    float64(c.V),
		}
		if err := checkCandle(result[i].Open, result[i].High, result[i].Low, result[i].Close); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func checkCandle(data ...float64) error {
	for _, val := range data {
		if int64(val) == 0 {
			return fmt.Errorf("candle error, %v", data)
		}
	}
	return nil
}

func (gt *gateioTrader) CalcPositionSize(ctx context.Context) (int64, error) {
	return 1, nil
}

func (gt *gateioTrader) ClosePosition(ctx context.Context) error { return nil }

func (gt *gateioTrader) CreateOrder(ctx context.Context, order *types.Order) error {
	forder := gateapi.FuturesOrder{
		Contract:     order.Contract,
		Size:         order.Size,
		Price:        fmt.Sprintf("%.2f", order.Price),
		Tif:          "ioc",
		IsReduceOnly: true,
	}

	if order.Side.IsShort() {
		forder.Size = -order.Size
	}

	_, _, err := gt.client.FuturesApi.CreateFuturesOrder(ctx, settle, forder, nil)
	if err != nil {
		logger.Errorf("create order failed: %v, price=%v, size=%v", err, order.Price, order.Size)
	} else {
		logger.Infof("create order successful, price=%v, size=%v", order.Price, order.Size)
	}
	return err
}
