package indicators

import "testing"

func TestAddIndicatorsAttachesEnabledResults(t *testing.T) {
	rows := []KlineInput{
		{Date: "2024-06-06", Open: Float(9), High: Float(10), Low: Float(8), Close: Float(9), Volume: Float(100)},
		{Date: "2024-06-07", Open: Float(10), High: Float(11), Low: Float(9), Close: Float(10), Volume: Float(110)},
		{Date: "2024-06-10", Open: Float(11), High: Float(12), Low: Float(10), Close: Float(11), Volume: Float(120)},
		{Date: "2024-06-11", Open: Float(12), High: Float(13), Low: Float(11), Close: Float(12), Volume: Float(130)},
	}

	got := AddIndicators(rows, Options{
		MA:  &MAOptions{Periods: []int{3}},
		RSI: &RSIOptions{Periods: []int{2}},
	})

	if len(got) != len(rows) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(rows))
	}
	if got[0].Date != rows[0].Date || got[3].Close != rows[3].Close {
		t.Fatalf("kline fields were not preserved: %+v", got[3])
	}
	assertFloat(t, got[2].MA["ma3"], 10)
	assertFloat(t, got[3].MA["ma3"], 11)
	if got[2].RSI["rsi2"] == nil {
		t.Fatal("expected enabled RSI result on the third row")
	}
	if got[3].MACD != nil {
		t.Fatalf("disabled MACD result = %+v, want nil", got[3].MACD)
	}
}

func TestEstimateLookbackAppliesIndicatorBuffers(t *testing.T) {
	ma := EstimateLookback(Options{
		MA: &MAOptions{Periods: []int{3}},
	})
	if ma.MaxLookback != 3 || ma.RequiredBars != 4 || ma.HasEMABasedIndicator {
		t.Fatalf("MA lookback = %+v, want max=3 required=4 ema=false", ma)
	}

	macd := EstimateLookback(Options{
		MACD: &MACDOptions{Long: 26, Signal: 9},
	})
	if macd.MaxLookback != 87 || macd.RequiredBars != 131 || !macd.HasEMABasedIndicator {
		t.Fatalf("MACD lookback = %+v, want max=87 required=131 ema=true", macd)
	}
}

func TestGetEnabledKeysReturnsRegistryOrder(t *testing.T) {
	keys := GetEnabledKeys(Options{
		KC:  &KCOptions{},
		MA:  &MAOptions{},
		RSI: &RSIOptions{},
	})
	want := []IndicatorKey{IndicatorMA, IndicatorRSI, IndicatorKC}
	if len(keys) != len(want) {
		t.Fatalf("keys = %#v, want %#v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("keys[%d] = %s, want %s", i, keys[i], want[i])
		}
	}
}

func TestBuildContextExtractsClosesAndOHLCV(t *testing.T) {
	rows := []KlineInput{
		{Date: "2024-06-06", Open: Float(9), High: Float(10), Low: Float(8), Close: Float(9.5), Volume: Float(100)},
		{Date: "2024-06-07", Open: Float(10), High: Float(11), Low: Float(9), Close: nil, Volume: Float(110)},
	}

	context := BuildContext(rows)
	if len(context.Closes) != 2 || context.Closes[0] == nil || *context.Closes[0] != 9.5 || context.Closes[1] != nil {
		t.Fatalf("closes = %#v", context.Closes)
	}
	if len(context.OHLCV) != 2 || context.OHLCV[1].Close != nil || context.OHLCV[1].Volume == nil || *context.OHLCV[1].Volume != 110 {
		t.Fatalf("ohlcv = %#v", context.OHLCV)
	}
}

func TestTypeScriptStyleIndicatorNames(t *testing.T) {
	rows := []KlineInput{
		{Date: "2024-06-06", Open: Float(9), High: Float(10), Low: Float(8), Close: Float(9.5), Volume: Float(100)},
		{Date: "2024-06-07", Open: Float(10), High: Float(11), Low: Float(9), Close: Float(10.5), Volume: Float(110)},
	}
	options := Options{
		MA:   &MAOptions{Periods: []int{2}},
		MACD: &MACDOptions{},
	}

	context := BuildIndicatorContext(rows)
	if len(context.Closes) != 2 || context.Closes[1] == nil || *context.Closes[1] != 10.5 {
		t.Fatalf("context closes = %#v", context.Closes)
	}

	keys := GetEnabledIndicatorKeys(options)
	if len(keys) != 2 || keys[0] != IndicatorMA || keys[1] != IndicatorMACD {
		t.Fatalf("keys = %#v, want [ma macd]", keys)
	}

	lookback := EstimateIndicatorLookback(options)
	if lookback.RequiredBars == 0 || !lookback.HasEMABasedIndicator {
		t.Fatalf("lookback = %+v, want required bars with EMA flag", lookback)
	}

	registry := IndicatorRegistry()
	descriptor, ok := registry[IndicatorMA]
	if !ok || descriptor.Key != IndicatorMA {
		t.Fatalf("registry[ma] = %+v/%v", descriptor, ok)
	}
	if _, ok := INDICATOR_REGISTRY[IndicatorMACD]; !ok {
		t.Fatal("expected INDICATOR_REGISTRY to include macd")
	}
}
