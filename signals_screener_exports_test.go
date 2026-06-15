package stock

import "testing"

func TestRootReExportsSignals(t *testing.T) {
	klines := []SignalKline{
		{
			Timestamp: SignalTime(1000),
			Close:     SignalFloat(10),
			MA:        map[string]*float64{"ma5": SignalFloat(9), "ma10": SignalFloat(10)},
			RSI:       map[string]*float64{"rsi6": SignalFloat(50)},
		},
		{
			Timestamp: SignalTime(2000),
			Close:     SignalFloat(11),
			MA:        map[string]*float64{"ma5": SignalFloat(11), "ma10": SignalFloat(10)},
			RSI:       map[string]*float64{"rsi6": SignalFloat(80)},
		},
	}

	result, err := CalcSignals(klines, SignalOptions{
		MA:  &SignalMAOptions{Fast: 5, Slow: 10},
		RSI: &SignalRSIOptions{Period: 6},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 || result[0].Type != SignalMAGoldenCross || result[1].Type != SignalRSIOverbought {
		t.Fatalf("signals = %+v", result)
	}
}

func TestRootCalcSignalsInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	klines := []SignalKline{
		{
			Timestamp: SignalTime(1000),
			MA:        map[string]*float64{"ma5": SignalFloat(1), "ma10": SignalFloat(2)},
			RSI:       map[string]*float64{"rsi6": SignalFloat(50)},
		},
	}
	tests := []struct {
		name    string
		options SignalOptions
	}{
		{
			name:    "missing ma key",
			options: SignalOptions{MA: &SignalMAOptions{Fast: 3, Slow: 10}},
		},
		{
			name:    "missing rsi key",
			options: SignalOptions{RSI: &SignalRSIOptions{Period: 12}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CalcSignals(klines, tt.options)
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestRootReExportsScreenerAndBacktest(t *testing.T) {
	items := []struct {
		Code   string
		Amount *float64
	}{
		{Code: "a", Amount: ScreenerFloat(1)},
		{Code: "b", Amount: ScreenerFloat(3)},
		{Code: "c", Amount: ScreenerFloat(2)},
	}

	picks, err := Screen(items).
		Where(func(item struct {
			Code   string
			Amount *float64
		}) bool {
			return item.Amount != nil && *item.Amount >= 2
		}).
		SortBy(func(item struct {
			Code   string
			Amount *float64
		}) *float64 {
			return item.Amount
		}, Desc).
		Top(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(picks) != 1 || picks[0].Code != "b" {
		t.Fatalf("picks = %+v", picks)
	}

	report := Backtest(BacktestOptions[struct{ Close *float64 }]{
		Klines: []struct{ Close *float64 }{
			{Close: ScreenerFloat(10)},
			{Close: ScreenerFloat(12)},
		},
		InitialCapital: 1000,
		GetClose:       func(item struct{ Close *float64 }) *float64 { return item.Close },
		Strategy: func(_ struct{ Close *float64 }, index int, _ []struct{ Close *float64 }) StrategySignal {
			if index == 0 {
				return Buy
			}
			return Hold
		},
	})
	if report.TradeCount != 1 || report.TotalReturn <= 0 {
		t.Fatalf("report = %+v", report)
	}
}

func TestRootScreenerTopInvalidArgumentErrorCode(t *testing.T) {
	_, err := Screen([]int{1, 2, 3}).Top(-1)
	if err == nil {
		t.Fatal("expected invalid argument error")
	}
	if code := GetErrorCode(err); code != CodeInvalidArgument {
		t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
	}
}
