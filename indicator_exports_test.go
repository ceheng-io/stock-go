package stock

import "testing"

func TestRootReExportsIndicatorUtilities(t *testing.T) {
	rows := []KlineInput{
		{Date: "2024-06-06", Open: Float(9), High: Float(10), Low: Float(8), Close: Float(9), Volume: Float(100)},
		{Date: "2024-06-07", Open: Float(10), High: Float(11), Low: Float(9), Close: Float(10), Volume: Float(110)},
		{Date: "2024-06-10", Open: Float(11), High: Float(12), Low: Float(10), Close: Float(11), Volume: Float(120)},
	}
	options := IndicatorOptions{
		MA:  &MAOptions{Periods: []int{3}},
		RSI: &RSIOptions{Periods: []int{2}},
		KC:  &KCOptions{EMAPeriod: 2, ATRPeriod: 2},
	}

	keys := GetEnabledIndicatorKeys(options)
	want := []IndicatorKey{IndicatorMA, IndicatorRSI, IndicatorKC}
	if len(keys) != len(want) {
		t.Fatalf("keys = %#v, want %#v", keys, want)
	}
	registry := IndicatorRegistry()
	if len(registry) < len(want) || registry[IndicatorMA].Key != IndicatorMA {
		t.Fatalf("IndicatorRegistry = %#v", registry)
	}
	if len(INDICATOR_REGISTRY) != len(registry) {
		t.Fatalf("INDICATOR_REGISTRY len = %d, registry len = %d", len(INDICATOR_REGISTRY), len(registry))
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("keys[%d] = %s, want %s", i, keys[i], want[i])
		}
	}

	context := BuildIndicatorContext(rows)
	if len(context.Closes) != len(rows) || context.Closes[2] == nil || *context.Closes[2] != 11 {
		t.Fatalf("context = %+v", context)
	}

	enriched := AddIndicators(rows, options)
	if len(enriched) != len(rows) || enriched[2].MA["ma3"] == nil || *enriched[2].MA["ma3"] != 10 {
		t.Fatalf("enriched = %+v", enriched)
	}

	lookback := EstimateIndicatorLookback(options)
	if lookback.MaxLookback == 0 || !lookback.HasEMABasedIndicator || lookback.RequiredBars < lookback.MaxLookback {
		t.Fatalf("lookback = %+v", lookback)
	}
}

func TestRootReExportsIndicatorCalculators(t *testing.T) {
	closes := Values(1, 2, 3, 4)

	if got := CalcSMA(closes, 3); got[2] == nil || *got[2] != 2 {
		t.Fatalf("CalcSMA = %#v", got)
	}
	if got := CalcMA(closes, MAOptions{Periods: []int{3}}); got[2]["ma3"] == nil || *got[2]["ma3"] != 2 {
		t.Fatalf("CalcMA = %#v", got)
	}
	if got := CalcRSI(closes, RSIOptions{Periods: []int{2}}); got[2]["rsi2"] == nil {
		t.Fatalf("CalcRSI = %#v", got)
	}

	ohlcv := []OHLCV{
		{Open: Float(1), High: Float(2), Low: Float(1), Close: Float(1), Volume: Float(10)},
		{Open: Float(2), High: Float(3), Low: Float(2), Close: Float(2), Volume: Float(20)},
		{Open: Float(3), High: Float(4), Low: Float(3), Close: Float(3), Volume: Float(30)},
	}
	if got := CalcATR(ohlcv, ATROptions{Period: 2}); got[1].ATR == nil {
		t.Fatalf("CalcATR = %#v", got)
	}
	if got := CalcOBV(ohlcv, OBVOptions{MAPeriod: 2}); got[2].OBV == nil {
		t.Fatalf("CalcOBV = %#v", got)
	}
}
