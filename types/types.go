package types

import "time"

// K线数据结构
type Candle struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// 指标计算结果
type Indicators struct {
	Timestamp       time.Time
	SuperTrendValue float64
	SuperTrendTrend string // "up" 或 "down"
	RSI             float64
	BollUpper       float64
	BollMiddle      float64
	BollLower       float64
	Price           float64 // 当前价格
}

type Order struct {
	Contract string
	Size     int64
	Price    float64
	Side     SideType
	Close    bool
}

type SideType int

const (
	SideTypeLong SideType = iota + 1
	SideTypeShort
)

func (st SideType) IsLong() bool {
	return st == SideTypeLong
}

func (st SideType) IsShort() bool {
	return st == SideTypeShort
}

func (st SideType) String() string {
	if st == SideTypeLong {
		return "long"
	}
	return "short"
}

type DecisionType int64 // 决策类型

const (
	DecisionLong    DecisionType = iota + 1 // 多仓
	DecisionShort                           // 空仓
	DecisionClose                           // 平仓
	DecisionNothing                         // 无操作
)
