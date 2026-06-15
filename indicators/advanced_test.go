package indicators

import "testing"

func TestCalcBIAS(t *testing.T) {
	bias := CalcBIAS(Values(10, 12, 14, 16), BIASOptions{Periods: []int{3}})

	if len(bias) != 4 {
		t.Fatalf("len(bias) = %d, want 4", len(bias))
	}
	assertNil(t, bias[0]["bias3"])
	assertNil(t, bias[1]["bias3"])
	assertFloat(t, bias[2]["bias3"], 16.67)
	assertFloat(t, bias[3]["bias3"], 14.29)
}

func TestCalcCCI(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(10), Close: Float(10)},
		{High: Float(11), Low: Float(11), Close: Float(11)},
		{High: Float(12), Low: Float(12), Close: Float(12)},
	}

	cci := CalcCCI(data, CCIOptions{Period: 3})

	if len(cci) != 3 {
		t.Fatalf("len(cci) = %d, want 3", len(cci))
	}
	assertNil(t, cci[0].CCI)
	assertNil(t, cci[1].CCI)
	assertFloat(t, cci[2].CCI, 100)
}

func TestCalcCCIHandlesFlatTypicalPrice(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(10), Close: Float(10)},
		{High: Float(10), Low: Float(10), Close: Float(10)},
		{High: Float(10), Low: Float(10), Close: Float(10)},
	}

	cci := CalcCCI(data, CCIOptions{Period: 3})

	assertFloat(t, cci[2].CCI, 0)
}

func TestCalcATR(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(8), Close: Float(9)},
		{High: Float(12), Low: Float(9), Close: Float(11)},
		{High: Float(13), Low: Float(10), Close: Float(12)},
		{High: Float(16), Low: Float(12), Close: Float(15)},
	}

	atr := CalcATR(data, ATROptions{Period: 3})

	if len(atr) != 4 {
		t.Fatalf("len(atr) = %d, want 4", len(atr))
	}
	assertFloat(t, atr[0].TR, 2)
	assertNil(t, atr[0].ATR)
	assertFloat(t, atr[2].TR, 3)
	assertFloat(t, atr[2].ATR, 2.67)
	assertFloat(t, atr[3].TR, 4)
	assertFloat(t, atr[3].ATR, 3.11)
}
