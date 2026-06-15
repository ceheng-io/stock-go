package indicators

import "testing"

func TestCalcDMI(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(8), Close: Float(9)},
		{High: Float(12), Low: Float(9), Close: Float(11)},
		{High: Float(13), Low: Float(10), Close: Float(12)},
		{High: Float(12), Low: Float(8), Close: Float(9)},
		{High: Float(14), Low: Float(9), Close: Float(13)},
		{High: Float(15), Low: Float(11), Close: Float(14)},
	}

	dmi := CalcDMI(data, DMIOptions{Period: 2, ADXPeriod: 2})

	if len(dmi) != 6 {
		t.Fatalf("len(dmi) = %d, want 6", len(dmi))
	}
	assertNil(t, dmi[1].PDI)
	assertFloat(t, dmi[2].PDI, 50)
	assertFloat(t, dmi[2].MDI, 0)
	assertNil(t, dmi[2].ADX)
	assertFloat(t, dmi[3].ADX, 57.14)
	assertFloat(t, dmi[5].ADXR, 57.85)
}

func TestCalcSAR(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(8), Close: Float(9)},
		{High: Float(11), Low: Float(9), Close: Float(10)},
		{High: Float(12), Low: Float(10), Close: Float(11)},
		{High: Float(13), Low: Float(11), Close: Float(12)},
		{High: Float(9), Low: Float(7), Close: Float(8)},
	}

	sar := CalcSAR(data, SAROptions{})

	if len(sar) != 5 {
		t.Fatalf("len(sar) = %d, want 5", len(sar))
	}
	assertNil(t, sar[0].SAR)
	assertFloat(t, sar[1].SAR, 8)
	assertInt(t, sar[1].Trend, 1)
	assertFloat(t, sar[3].SAR, 8.24)
	assertInt(t, sar[4].Trend, -1)
	assertFloat(t, sar[4].SAR, 13)
	assertFloat(t, sar[4].EP, 7)
	assertFloat(t, sar[4].AF, 0.02)
}

func TestCalcKC(t *testing.T) {
	data := []OHLCV{
		{High: Float(10), Low: Float(8), Close: Float(9)},
		{High: Float(12), Low: Float(9), Close: Float(11)},
		{High: Float(13), Low: Float(10), Close: Float(12)},
		{High: Float(16), Low: Float(12), Close: Float(15)},
	}

	kc := CalcKC(data, KCOptions{EMAPeriod: 3, ATRPeriod: 3, Multiplier: 2})

	if len(kc) != 4 {
		t.Fatalf("len(kc) = %d, want 4", len(kc))
	}
	assertNil(t, kc[1].Mid)
	assertFloat(t, kc[2].Mid, 10.67)
	assertFloat(t, kc[2].Upper, 16.01)
	assertFloat(t, kc[2].Lower, 5.33)
	assertFloat(t, kc[2].Width, 100.09)
	assertFloat(t, kc[3].Mid, 12.83)
	assertFloat(t, kc[3].Upper, 19.05)
	assertFloat(t, kc[3].Lower, 6.61)
	assertFloat(t, kc[3].Width, 96.96)
}

func assertInt(t *testing.T, value *int, want int) {
	t.Helper()
	if value == nil {
		t.Fatalf("value is nil, want %d", want)
	}
	if *value != want {
		t.Fatalf("value = %d, want %d", *value, want)
	}
}
