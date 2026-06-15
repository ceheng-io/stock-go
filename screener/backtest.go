package screener

// StrategySignal is a single-bar trading action.
type StrategySignal string

const (
	Buy  StrategySignal = "buy"
	Sell StrategySignal = "sell"
	Hold StrategySignal = "hold"
)

// Strategy returns one trading signal for the current bar.
type Strategy[T any] func(bar T, index int, history []T) StrategySignal

// Trade is one completed long trade.
type Trade struct {
	EntryIndex    int
	ExitIndex     int
	EntryPrice    float64
	ExitPrice     float64
	ReturnPercent float64
}

// BacktestOptions configures a simple all-in long-only backtest.
type BacktestOptions[T any] struct {
	Klines         []T
	Strategy       Strategy[T]
	InitialCapital float64
	Fee            float64
	GetClose       func(T) *float64
}

// BacktestReport contains simple performance metrics.
type BacktestReport struct {
	InitialCapital float64
	FinalEquity    float64
	TotalReturn    float64
	WinRate        float64
	MaxDrawdown    float64
	TradeCount     int
	Trades         []Trade
	EquityCurve    []float64
}

// Backtest runs a simple all-in long-only backtest.
func Backtest[T any](options BacktestOptions[T]) BacktestReport {
	initialCapital := options.InitialCapital
	if initialCapital == 0 {
		initialCapital = 100000
	}
	getClose := options.GetClose
	if getClose == nil {
		getClose = func(T) *float64 { return nil }
	}
	strategy := options.Strategy
	if strategy == nil {
		strategy = func(T, int, []T) StrategySignal { return Hold }
	}

	cash := initialCapital
	position := 0.0
	entryPrice := 0.0
	entryIndex := -1
	lastPrice := 0.0
	trades := make([]Trade, 0)
	equityCurve := make([]float64, 0, len(options.Klines))

	recordTrade := func(exitIndex int, exitPrice float64) {
		returnPercent := ((exitPrice/entryPrice)*(1-options.Fee)*(1-options.Fee) - 1) * 100
		trades = append(trades, Trade{
			EntryIndex:    entryIndex,
			ExitIndex:     exitIndex,
			EntryPrice:    entryPrice,
			ExitPrice:     exitPrice,
			ReturnPercent: returnPercent,
		})
	}

	for i, bar := range options.Klines {
		price, ok := validPrice(getClose(bar))
		signal := strategy(bar, i, options.Klines)
		if ok {
			lastPrice = price
			if signal == Buy && position == 0 {
				position = cash * (1 - options.Fee) / price
				entryPrice = price
				entryIndex = i
				cash = 0
			} else if signal == Sell && position > 0 {
				cash = position * price * (1 - options.Fee)
				recordTrade(i, price)
				position = 0
			}
		}

		mark := lastPrice
		if ok {
			mark = price
		}
		equityCurve = append(equityCurve, cash+position*mark)
	}

	if position > 0 && lastPrice > 0 {
		cash = position * lastPrice * (1 - options.Fee)
		recordTrade(len(options.Klines)-1, lastPrice)
		position = 0
		if len(equityCurve) > 0 {
			equityCurve[len(equityCurve)-1] = cash
		}
	}

	finalEquity := initialCapital
	if len(equityCurve) > 0 {
		finalEquity = equityCurve[len(equityCurve)-1]
	}
	totalReturn := (finalEquity/initialCapital - 1) * 100
	return BacktestReport{
		InitialCapital: initialCapital,
		FinalEquity:    finalEquity,
		TotalReturn:    totalReturn,
		WinRate:        winRate(trades),
		MaxDrawdown:    maxDrawdown(equityCurve),
		TradeCount:     len(trades),
		Trades:         trades,
		EquityCurve:    equityCurve,
	}
}

func validPrice(value *float64) (float64, bool) {
	number, ok := finiteValue(value)
	return number, ok && number > 0
}

func winRate(trades []Trade) float64 {
	if len(trades) == 0 {
		return 0
	}
	wins := 0
	for _, trade := range trades {
		if trade.ReturnPercent > 0 {
			wins++
		}
	}
	return float64(wins) / float64(len(trades)) * 100
}

func maxDrawdown(equityCurve []float64) float64 {
	peak := 0.0
	drawdown := 0.0
	for _, equity := range equityCurve {
		if equity > peak {
			peak = equity
		}
		if peak > 0 {
			current := (peak - equity) / peak * 100
			if current > drawdown {
				drawdown = current
			}
		}
	}
	return drawdown
}
