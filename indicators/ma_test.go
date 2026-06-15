package indicators

import "testing"

func TestCalcSMA(t *testing.T) {
	data := Values(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	sma := CalcSMA(data, 3)

	assertNil(t, sma[0])
	assertNil(t, sma[1])
	assertFloat(t, sma[2], 2)
	assertFloat(t, sma[3], 3)
	assertFloat(t, sma[9], 9)
}

func TestCalcSMAHandlesNullValues(t *testing.T) {
	data := []Value{Float(1), Float(2), Null(), Float(4), Float(5)}

	sma := CalcSMA(data, 3)

	assertNil(t, sma[2])
}

func TestCalcEMA(t *testing.T) {
	data := Values(10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

	ema := CalcEMA(data, 3)

	assertNil(t, ema[0])
	assertNil(t, ema[1])
	assertFloat(t, ema[2], 11)
	if ema[3] == nil {
		t.Fatal("ema[3] is nil")
	}
}

func TestCalcWMA(t *testing.T) {
	data := Values(1, 2, 3, 4, 5)

	wma := CalcWMA(data, 3)

	assertNil(t, wma[0])
	assertNil(t, wma[1])
	assertFloat(t, wma[2], 2.33)
}

func TestCalcMA(t *testing.T) {
	data := Values(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	ma := CalcMA(data, MAOptions{Periods: []int{3, 5}})

	if len(ma) != 10 {
		t.Fatalf("len(ma) = %d, want 10", len(ma))
	}
	assertFloat(t, ma[2]["ma3"], 2)
	assertFloat(t, ma[4]["ma5"], 3)
}

func TestCalcMASupportsEMA(t *testing.T) {
	data := Values(10, 11, 12, 13, 14)

	ma := CalcMA(data, MAOptions{Periods: []int{3}, Type: MATypeEMA})

	assertFloat(t, ma[2]["ma3"], 11)
}

func assertNil(t *testing.T, value *float64) {
	t.Helper()
	if value != nil {
		t.Fatalf("value = %v, want nil", *value)
	}
}

func assertFloat(t *testing.T, value *float64, want float64) {
	t.Helper()
	if value == nil {
		t.Fatalf("value is nil, want %.2f", want)
	}
	if *value != want {
		t.Fatalf("value = %.2f, want %.2f", *value, want)
	}
}
