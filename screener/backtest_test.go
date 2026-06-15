package screener

import "testing"

type barFixture struct {
	Close *float64
}

func TestBacktestBuysSellsAndReportsPerformance(t *testing.T) {
	klines := []barFixture{
		{Close: Float(10)},
		{Close: Float(12)},
		{Close: Float(9)},
		{Close: Float(11)},
	}

	report := Backtest(BacktestOptions[barFixture]{
		Klines:         klines,
		InitialCapital: 1000,
		Fee:            0.001,
		GetClose:       func(item barFixture) *float64 { return item.Close },
		Strategy: func(_ barFixture, index int, _ []barFixture) StrategySignal {
			switch index {
			case 0:
				return Buy
			case 1:
				return Sell
			case 2:
				return Buy
			default:
				return Hold
			}
		},
	})

	assertFloat64(t, report.FinalEquity, 1460.81)
	assertFloat64(t, report.TotalReturn, 46.08)
	assertFloat64(t, report.WinRate, 100)
	assertFloat64(t, report.MaxDrawdown, 0.10)
	if report.TradeCount != 2 || len(report.Trades) != 2 {
		t.Fatalf("trades = %+v", report.Trades)
	}
	assertFloat64(t, report.Trades[0].ReturnPercent, 19.76)
	assertFloat64(t, report.Trades[1].ReturnPercent, 21.98)
}

func TestBacktestMarksInvalidCloseWithLastPrice(t *testing.T) {
	klines := []barFixture{
		{Close: Float(10)},
		{Close: nil},
		{Close: Float(11)},
	}

	report := Backtest(BacktestOptions[barFixture]{
		Klines:         klines,
		InitialCapital: 1000,
		GetClose:       func(item barFixture) *float64 { return item.Close },
		Strategy: func(_ barFixture, index int, _ []barFixture) StrategySignal {
			if index == 0 {
				return Buy
			}
			return Hold
		},
	})

	if len(report.EquityCurve) != 3 {
		t.Fatalf("equity curve = %+v", report.EquityCurve)
	}
	assertFloat64(t, report.EquityCurve[1], 1000)
	assertFloat64(t, report.FinalEquity, 1100)
}

func assertFloat64(t *testing.T, got float64, want float64) {
	t.Helper()
	rounded := round(got)
	if rounded != want {
		t.Fatalf("value = %.2f, want %.2f", rounded, want)
	}
}
