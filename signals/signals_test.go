package signals

import "testing"

func TestCalcSignalsDetectsCrossesAndThresholds(t *testing.T) {
	klines := []Kline{
		{
			Timestamp: FloatTime(1000),
			Close:     Float(10),
			MA:        map[string]*float64{"ma5": Float(9), "ma10": Float(10)},
			MACD:      &MACD{DIF: Float(0.1), DEA: Float(0.2)},
			KDJ:       &KDJ{K: Float(15), D: Float(20)},
			RSI:       map[string]*float64{"rsi6": Float(35)},
			BOLL:      &BOLL{Upper: Float(12), Lower: Float(8)},
			SAR:       &SAR{Trend: Int(1)},
		},
		{
			Timestamp: FloatTime(2000),
			Close:     Float(13),
			MA:        map[string]*float64{"ma5": Float(11), "ma10": Float(10)},
			MACD:      &MACD{DIF: Float(0.3), DEA: Float(0.2)},
			KDJ:       &KDJ{K: Float(85), D: Float(80)},
			RSI:       map[string]*float64{"rsi6": Float(75)},
			BOLL:      &BOLL{Upper: Float(12), Lower: Float(8)},
			SAR:       &SAR{Trend: Int(-1)},
		},
	}

	result, err := CalcSignals(klines, SignalOptions{
		MA:   &MAOptions{Fast: 5, Slow: 10},
		MACD: true,
		KDJ:  &KDJOptions{},
		RSI:  &RSIOptions{Period: 6},
		BOLL: true,
		SAR:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	types := signalTypes(result)
	want := []SignalType{
		SignalMAGoldenCross,
		SignalMACDGoldenCross,
		SignalKDJGoldenCross,
		SignalKDJOverbought,
		SignalRSIOverbought,
		SignalBOLLBreakUpper,
		SignalSARReversalDown,
	}
	if len(types) != len(want) {
		t.Fatalf("types = %#v, want %#v", types, want)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("types[%d] = %q, want %q; all=%#v", i, types[i], want[i], types)
		}
		if result[i].At != 2000 || result[i].Index != 1 {
			t.Fatalf("signal[%d] = %+v", i, result[i])
		}
	}
}

func TestCalcSignalsDetectsDeathAndOversold(t *testing.T) {
	klines := []Kline{
		{
			Timestamp: FloatTime(1000),
			Close:     Float(10),
			MA:        map[string]*float64{"ma5": Float(11), "ma10": Float(10)},
			MACD:      &MACD{DIF: Float(0.3), DEA: Float(0.2)},
			KDJ:       &KDJ{K: Float(25), D: Float(20)},
			RSI:       map[string]*float64{"rsi6": Float(35)},
			BOLL:      &BOLL{Upper: Float(12), Lower: Float(8)},
			SAR:       &SAR{Trend: Int(-1)},
		},
		{
			Timestamp: FloatTime(2000),
			Close:     Float(7),
			MA:        map[string]*float64{"ma5": Float(9), "ma10": Float(10)},
			MACD:      &MACD{DIF: Float(0.1), DEA: Float(0.2)},
			KDJ:       &KDJ{K: Float(15), D: Float(20)},
			RSI:       map[string]*float64{"rsi6": Float(25)},
			BOLL:      &BOLL{Upper: Float(12), Lower: Float(8)},
			SAR:       &SAR{Trend: Int(1)},
		},
	}

	result, err := CalcSignals(klines, SignalOptions{
		MA:   &MAOptions{Fast: 5, Slow: 10},
		MACD: true,
		KDJ:  &KDJOptions{},
		RSI:  &RSIOptions{Period: 6},
		BOLL: true,
		SAR:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	types := signalTypes(result)
	want := []SignalType{
		SignalMADeathCross,
		SignalMACDDeathCross,
		SignalKDJDeathCross,
		SignalKDJOversold,
		SignalRSIOversold,
		SignalBOLLBreakLower,
		SignalSARReversalUp,
	}
	if len(types) != len(want) {
		t.Fatalf("types = %#v, want %#v", types, want)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("types[%d] = %q, want %q; all=%#v", i, types[i], want[i], types)
		}
	}
}

func TestCalcSignalsSkipsMissingTimestampAndValidatesIndicatorKeys(t *testing.T) {
	klines := []Kline{
		{Timestamp: FloatTime(1000), MA: map[string]*float64{"ma5": Float(1), "ma10": Float(2)}, RSI: map[string]*float64{"rsi6": Float(50)}},
		{Timestamp: nil, MA: map[string]*float64{"ma5": Float(3), "ma10": Float(2)}, RSI: map[string]*float64{"rsi6": Float(80)}},
	}

	result, err := CalcSignals(klines, SignalOptions{
		MA:  &MAOptions{Fast: 5, Slow: 10},
		RSI: &RSIOptions{Period: 6},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("result = %+v, want empty", result)
	}

	_, err = CalcSignals(klines, SignalOptions{MA: &MAOptions{Fast: 3, Slow: 10}})
	if err == nil {
		t.Fatal("expected missing MA key error")
	}
	_, err = CalcSignals(klines, SignalOptions{RSI: &RSIOptions{Period: 12}})
	if err == nil {
		t.Fatal("expected missing RSI key error")
	}
}

func signalTypes(signals []Signal) []SignalType {
	types := make([]SignalType, len(signals))
	for i, signal := range signals {
		types[i] = signal.Type
	}
	return types
}
