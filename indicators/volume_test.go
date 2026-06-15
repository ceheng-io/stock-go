package indicators

import "testing"

func TestCalcOBV(t *testing.T) {
	data := []OHLCV{
		{Close: Float(10), Volume: Float(100)},
		{Close: Float(12), Volume: Float(50)},
		{Close: Float(11), Volume: Float(30)},
		{Close: Float(11), Volume: Float(20)},
	}

	obv := CalcOBV(data, OBVOptions{})

	if len(obv) != 4 {
		t.Fatalf("len(obv) = %d, want 4", len(obv))
	}
	assertFloat(t, obv[0].OBV, 100)
	assertFloat(t, obv[1].OBV, 150)
	assertFloat(t, obv[2].OBV, 120)
	assertFloat(t, obv[3].OBV, 120)
	assertNil(t, obv[3].OBVMA)
}

func TestCalcOBVWithMA(t *testing.T) {
	data := []OHLCV{
		{Close: Float(10), Volume: Float(100)},
		{Close: Float(12), Volume: Float(50)},
		{Close: Float(11), Volume: Float(30)},
	}

	obv := CalcOBV(data, OBVOptions{MAPeriod: 2})

	assertNil(t, obv[0].OBVMA)
	assertFloat(t, obv[1].OBVMA, 125)
	assertFloat(t, obv[2].OBVMA, 135)
}

func TestCalcROC(t *testing.T) {
	data := []OHLCV{
		{Close: Float(10)},
		{Close: Float(12)},
		{Close: Float(15)},
		{Close: Float(18)},
	}

	roc := CalcROC(data, ROCOptions{Period: 2, SignalPeriod: 2})

	if len(roc) != 4 {
		t.Fatalf("len(roc) = %d, want 4", len(roc))
	}
	assertNil(t, roc[1].ROC)
	assertFloat(t, roc[2].ROC, 50)
	assertFloat(t, roc[3].ROC, 50)
	assertFloat(t, roc[3].Signal, 50)
}

func TestCalcROCHandlesMissingOrZeroPreviousClose(t *testing.T) {
	data := []OHLCV{
		{Close: Float(0)},
		{Close: Float(10)},
		{Close: Null()},
	}

	roc := CalcROC(data, ROCOptions{Period: 2})

	assertNil(t, roc[2].ROC)
}
