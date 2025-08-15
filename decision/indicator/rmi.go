package indicator

import (
	"math"

	"github.com/leaftree/autotrader/types"
)

// 批量处理结果
type BatchResult struct {
	BuySignal    bool    // 买入信号
	SellSignal   bool    // 卖出信号
	FinalRmiMfi  float64 // 最终的RMI_MFI值
	FinalRWMA    float64 // 最终的RWMA值
	FinalBand    float64 // 最终的Band值
	FinalChannel float64 // 最终的通道位置
}

// 指标配置参数
type Config struct {
	Length    int     // RMI长度 (默认14)
	PMom      float64 // 正动量阈值 (默认66)
	NMon      float64 // 负动量阈值 (默认30)
	BandMulti float64 // 通道乘数 (默认0.3)
	// 小窗口(如5-10)：对价格变化更敏感，反应更快, 适合短线交易，捕捉短期趋势
	// 大窗口(如20-50)：更平滑，过滤市场噪音, 长窗口：适合长线投资，识别主要趋势
	Window int // RWMA窗口大小 (默认20)
}

// 指标计算状态
type RmiTrendSniper struct {
	config          Config
	prevRmiMfi      float64
	prevClose       float64
	prevTypical     float64
	positive        bool
	negative        bool
	prevPositive    bool
	prevNegative    bool
	ema5            float64
	prevEma5        float64
	atrValues       []float64
	barRangeHistory []float64
	closeHistory    []float64
	prevPositiveMF  float64
	prevNegativeMF  float64
	rwmaUp          float64
	rwmaDown        float64
}

// 初始化指标计算器
func NewRmiTrendSniper(config Config) *RmiTrendSniper {
	if config.Window == 0 {
		config.Window = 20 // 默认窗口大小
	}

	return &RmiTrendSniper{
		config:          config,
		barRangeHistory: make([]float64, 0, config.Window+10),
		closeHistory:    make([]float64, 0, config.Window+10),
		atrValues:       make([]float64, 0, 50),
	}
}

// 批量处理K线数据，返回综合信号
func (r *RmiTrendSniper) Process(candles []types.Candle) types.DecisionType {
	result := BatchResult{}

	for _, candle := range candles {
		r.update(candle)
	}

	// 基于最后的状态生成综合信号
	if r.positive {
		result.BuySignal = true
	} else if r.negative {
		result.SellSignal = true
	}

	// 获取最终指标值
	if len(r.barRangeHistory) > 0 {
		result.FinalRWMA = r.calculateRWMA(r.barRangeHistory[len(r.barRangeHistory)-1], r.closeHistory[len(r.closeHistory)-1])
	}
	if len(r.atrValues) > 0 {
		result.FinalBand = r.atrValues[len(r.atrValues)-1] * r.config.BandMulti * 4
	}

	result.FinalRmiMfi = r.prevRmiMfi

	if r.positive {
		result.FinalChannel = result.FinalRWMA - result.FinalBand
	} else if r.negative {
		result.FinalChannel = result.FinalRWMA + result.FinalBand
	}

	if result.BuySignal {
		return types.DecisionLong
	}
	return types.DecisionShort
}

// 处理单个K线（内部方法，不直接返回结果）
func (r *RmiTrendSniper) update(candle types.Candle) {
	// 1. 计算基础值
	barRange := candle.High - candle.Low
	hlc3 := (candle.High + candle.Low + candle.Close) / 3

	// 2. 更新EMA5
	r.prevEma5 = r.ema5
	if r.ema5 == 0 {
		r.ema5 = candle.Close
	} else {
		alpha := 2.0 / (5 + 1)
		r.ema5 = alpha*candle.Close + (1-alpha)*r.ema5
	}
	ema5Change := r.ema5 - r.prevEma5

	// 3. 计算RMI
	change := candle.Close - r.prevClose

	// 跳过第一根K线（没有前一收盘价）
	if r.prevClose == 0 {
		r.prevClose = candle.Close
		r.prevTypical = hlc3
		return
	}

	up := math.Max(change, 0)
	down := math.Max(-change, 0)

	// 使用指数平滑更新RMA
	r.rwmaUp = rma(r.rwmaUp, up, r.config.Length)
	r.rwmaDown = rma(r.rwmaDown, down, r.config.Length)

	var rsi float64
	if r.rwmaDown == 0 {
		rsi = 100
	} else if r.rwmaUp == 0 {
		rsi = 0
	} else {
		rsi = 100 - (100 / (1 + r.rwmaUp/r.rwmaDown))
	}

	// 4. 计算MFI
	mfi := r.calculateMFI(candle, hlc3)

	// 5. 组合指标
	rsiMfi := (rsi + mfi) / 2

	// 6. 动量信号检测
	pMom := r.prevRmiMfi < r.config.PMom &&
		rsiMfi > r.config.PMom &&
		rsiMfi > r.config.NMon &&
		ema5Change > 0

	nMom := rsiMfi < r.config.NMon && ema5Change < 0

	// 7. 更新状态机
	r.prevPositive = r.positive
	r.prevNegative = r.negative

	if pMom {
		r.positive = true
		r.negative = false
	} else if nMom {
		r.positive = false
		r.negative = true
	}

	// 8. 计算带状通道
	r.calculateBand(candle)
	r.calculateRWMA(barRange, candle.Close)

	// 9. 保存当前值供下次计算使用
	r.prevRmiMfi = rsiMfi
	r.prevClose = candle.Close
}

// 计算MFI（资金流量指数）
func (r *RmiTrendSniper) calculateMFI(candle types.Candle, hlc3 float64) float64 {
	// 1. 计算原始资金流
	rawMoneyFlow := hlc3 * candle.Volume

	// 2. 确定正负资金流
	var positiveRaw, negativeRaw float64
	if r.prevTypical > 0 { // 确保有前值
		if hlc3 > r.prevTypical {
			positiveRaw = rawMoneyFlow
		} else if hlc3 < r.prevTypical {
			negativeRaw = rawMoneyFlow
		}
	}

	// 3. 计算平滑资金流
	positiveMF := rma(r.prevPositiveMF, positiveRaw, r.config.Length)
	negativeMF := rma(r.prevNegativeMF, negativeRaw, r.config.Length)

	// 4. 计算MFI
	var mfi float64
	if negativeMF == 0 {
		mfi = 100
	} else if positiveMF == 0 {
		mfi = 0
	} else {
		moneyRatio := positiveMF / negativeMF
		mfi = 100 - (100 / (1 + moneyRatio))
	}

	// 5. 更新状态
	r.prevPositiveMF = positiveMF
	r.prevNegativeMF = negativeMF
	r.prevTypical = hlc3

	return mfi
}

// 计算带状通道宽度
func (r *RmiTrendSniper) calculateBand(candle types.Candle) {
	// 计算真实波幅(TR)
	highLow := candle.High - candle.Low
	highPrevClose := math.Abs(candle.High - r.prevClose)
	lowPrevClose := math.Abs(candle.Low - r.prevClose)
	tr := math.Max(highLow, math.Max(highPrevClose, lowPrevClose))

	// 初始ATR计算
	if len(r.atrValues) == 0 {
		r.atrValues = append(r.atrValues, tr)
	} else {
		// 使用指数平滑更新ATR
		atr := rma(r.atrValues[len(r.atrValues)-1], tr, 30)
		r.atrValues = append(r.atrValues, atr)
	}
}

// 计算范围加权移动平均线
func (r *RmiTrendSniper) calculateRWMA(barRange, closePrice float64) float64 {
	// 保存历史数据
	r.barRangeHistory = append(r.barRangeHistory, barRange)
	r.closeHistory = append(r.closeHistory, closePrice)

	// 维护固定长度的历史窗口
	if len(r.barRangeHistory) > r.config.Window {
		r.barRangeHistory = r.barRangeHistory[1:]
	}
	if len(r.closeHistory) > r.config.Window {
		r.closeHistory = r.closeHistory[1:]
	}

	// 计算总范围
	var sumRange float64
	for _, r := range r.barRangeHistory {
		sumRange += r
	}

	// 计算范围加权移动平均
	var weightSum, priceSum float64
	for i, br := range r.barRangeHistory {
		if i >= len(r.closeHistory) {
			break
		}
		weight := br / sumRange
		weightSum += weight
		priceSum += weight * r.closeHistory[i]
	}

	// 避免除以零
	if weightSum == 0 {
		return closePrice
	}
	return priceSum / weightSum
}

// RMA计算辅助函数（Wilder平滑）
func rma(prev, current float64, length int) float64 {
	if prev == 0 {
		return current
	}
	alpha := 1.0 / float64(length)
	return alpha*current + (1-alpha)*prev
}
