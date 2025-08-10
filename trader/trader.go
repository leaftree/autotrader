package trader

import (
	"context"

	"github.com/leaftree/autotrader/notification"
	"github.com/leaftree/autotrader/types"
)

type Trader interface {
	QueryCandles(ctx context.Context, contract string, limit int) ([]types.Candle, error)
	CreateOrder(ctx context.Context, order *types.Order) error
	ClosePosition(ctx context.Context) error
}

func NewTrader(exchange string) Trader {
	var t Trader
	switch exchange {
	case "gateio":
		t = NewGateIOTrader()
	}
	return t
}

type trader struct {
	tr Trader
}

func (t *trader) QueryCandles(ctx context.Context, contract string, limit int) ([]types.Candle, error) {
	return t.tr.QueryCandles(ctx, contract, limit)
}

func (t *trader) CreateOrder(ctx context.Context, order *types.Order) error {
	t.tr.CreateOrder(ctx, order)
	// TODO record
	// TODO notify
	notification.SendCreateOrderNotification(ctx, map[string]any{
		"side":            order.Side.String(),
		"price":           order.Price,
		"size":            order.Size,
		"stop_loss_price": 0,
	})
}
