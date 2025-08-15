package risk

import (
	"context"

	"github.com/leaftree/autotrader/types"
)

// TODO 检查出现多少次止损

var (
	PositionType types.SideType = types.SideTypeNone
)

func HasPosition(ctx context.Context, contract string) (bool, types.SideType) {
	return PositionType != types.SideTypeNone, PositionType
}
