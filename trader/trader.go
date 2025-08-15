package trader

import (
	"context"

	"github.com/leaftree/autotrader/notification"
	"github.com/leaftree/autotrader/types"
)

type Trader interface {
	QueryCandles(ctx context.Context, contract string, timeframe string, limit int) ([]types.Candle, error)
	CalcPositionSize(ctx context.Context) (int64, error)
	CreateOrder(ctx context.Context, order *types.Order) error
	ClosePosition(ctx context.Context) error
	//CreatePriceTriggerOrder(ctx context.Context) error
}

func NewTrader(exchange string) Trader {
	var t Trader
	switch exchange {
	case "gateio":
		t = NewGateIOTrader()
	}
	return &trader{tr: t}
}

type trader struct {
	tr Trader
}

func (t *trader) QueryCandles(ctx context.Context, contract string, timeframe string, limit int) ([]types.Candle, error) {
	return t.tr.QueryCandles(ctx, contract, timeframe, limit)
}

func (t *trader) CalcPositionSize(ctx context.Context) (int64, error) {
	// TODO: currently only used for ETH, return 1 by default
	return 1, nil
}

func (t *trader) CreateOrder(ctx context.Context, order *types.Order) error {
	//t.tr.CreateOrder(ctx, order)
	// TODO record
	// TODO notify
	notification.SendCreateOrderNotification(ctx, map[string]any{
		"contract":        order.Contract,
		"side":            order.Side.String(),
		"price":           order.Price,
		"size":            order.Size,
		"stop_loss_price": 0,
	})
	return nil
}

func (t *trader) ClosePosition(ctx context.Context) error { return nil }
