package stock

import "github.com/ceheng-io/stock-go/screener"

type SortDirection = screener.SortDirection

const (
	Asc  = screener.Asc
	Desc = screener.Desc
)

type ScreenerBuilder[T any] struct {
	inner *screener.Builder[T]
}

type StrategySignal = screener.StrategySignal

const (
	Buy  = screener.Buy
	Sell = screener.Sell
	Hold = screener.Hold
)

type Strategy[T any] func(bar T, index int, history []T) StrategySignal
type Trade = screener.Trade

type BacktestOptions[T any] struct {
	Klines         []T
	Strategy       Strategy[T]
	InitialCapital float64
	Fee            float64
	GetClose       func(T) *float64
}

type BacktestReport = screener.BacktestReport

func Screen[T any](items []T) *ScreenerBuilder[T] {
	return &ScreenerBuilder[T]{inner: screener.Screen(items)}
}

func (b *ScreenerBuilder[T]) Where(predicate func(T) bool) *ScreenerBuilder[T] {
	b.inner.Where(predicate)
	return b
}

func (b *ScreenerBuilder[T]) SortBy(selector func(T) *float64, direction SortDirection) *ScreenerBuilder[T] {
	b.inner.SortBy(selector, direction)
	return b
}

func (b *ScreenerBuilder[T]) Top(n int) ([]T, error) {
	return b.inner.Top(n)
}

func (b *ScreenerBuilder[T]) ToArray() ([]T, error) {
	return b.inner.ToArray()
}

func Backtest[T any](options BacktestOptions[T]) BacktestReport {
	var strategy screener.Strategy[T]
	if options.Strategy != nil {
		strategy = func(bar T, index int, history []T) screener.StrategySignal {
			return options.Strategy(bar, index, history)
		}
	}
	return screener.Backtest(screener.BacktestOptions[T]{
		Klines:         options.Klines,
		Strategy:       strategy,
		InitialCapital: options.InitialCapital,
		Fee:            options.Fee,
		GetClose:       options.GetClose,
	})
}

func ScreenerFloat(value float64) *float64 {
	return screener.Float(value)
}
