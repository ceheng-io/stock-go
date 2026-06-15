package indicators

import "testing"

func TestCalcMACD(t *testing.T) {
	data := make([]Value, 50)
	for i := range data {
		data[i] = Float(100 + float64(i))
	}

	macd := CalcMACD(data, MACDOptions{Short: 12, Long: 26, Signal: 9})

	if len(macd) != 50 {
		t.Fatalf("len(macd) = %d, want 50", len(macd))
	}
	assertNil(t, macd[0].DIF)
	if macd[49].DIF == nil || macd[49].DEA == nil || macd[49].MACD == nil {
		t.Fatalf("macd[49] has nil fields: %+v", macd[49])
	}
}

func TestCalcBOLL(t *testing.T) {
	data := make([]Value, 30)
	for i := range data {
		data[i] = Float(100 + float64(i%5) - 2)
	}

	boll := CalcBOLL(data, BOLLOptions{Period: 20, StdDev: 2})

	if len(boll) != 30 {
		t.Fatalf("len(boll) = %d, want 30", len(boll))
	}
	if boll[19].Mid == nil || boll[19].Upper == nil || boll[19].Lower == nil {
		t.Fatalf("boll[19] has nil bands: %+v", boll[19])
	}
	if *boll[19].Upper <= *boll[19].Mid {
		t.Fatalf("upper %.2f should be greater than mid %.2f", *boll[19].Upper, *boll[19].Mid)
	}
	if *boll[19].Lower >= *boll[19].Mid {
		t.Fatalf("lower %.2f should be less than mid %.2f", *boll[19].Lower, *boll[19].Mid)
	}
}

func TestCalcKDJ(t *testing.T) {
	data := make([]OHLCV, 20)
	for i := range data {
		x := float64(i)
		data[i] = OHLCV{
			Open:   Float(100 + x),
			High:   Float(105 + x),
			Low:    Float(95 + x),
			Close:  Float(102 + x),
			Volume: Float(1000),
		}
	}

	kdj := CalcKDJ(data, KDJOptions{Period: 9})

	if len(kdj) != 20 {
		t.Fatalf("len(kdj) = %d, want 20", len(kdj))
	}
	if kdj[8].K == nil || kdj[8].D == nil || kdj[8].J == nil {
		t.Fatalf("kdj[8] has nil fields: %+v", kdj[8])
	}
}

func TestCalcKDJHandlesNullAndFlatRanges(t *testing.T) {
	withNull := []OHLCV{
		{Open: Float(100), High: Float(110), Low: Float(90), Close: Float(100)},
		{Open: Float(100), High: Null(), Low: Float(95), Close: Float(102)},
		{Open: Float(100), High: Float(115), Low: Null(), Close: Float(105)},
		{Open: Float(100), High: Float(120), Low: Float(100), Close: Null()},
		{Open: Float(100), High: Float(125), Low: Float(105), Close: Float(115)},
		{Open: Float(100), High: Float(130), Low: Float(110), Close: Float(120)},
		{Open: Float(100), High: Float(135), Low: Float(115), Close: Float(125)},
		{Open: Float(100), High: Float(140), Low: Float(120), Close: Float(130)},
		{Open: Float(100), High: Float(145), Low: Float(125), Close: Float(135)},
		{Open: Float(100), High: Float(150), Low: Float(130), Close: Float(140)},
	}
	kdj := CalcKDJ(withNull, KDJOptions{Period: 5})
	assertNil(t, kdj[4].K)

	flat := make([]OHLCV, 10)
	for i := range flat {
		flat[i] = OHLCV{Open: Float(100), High: Float(100), Low: Float(100), Close: Float(100)}
	}
	kdj = CalcKDJ(flat, KDJOptions{Period: 5})
	assertNil(t, kdj[4].K)
}

func TestCalcRSI(t *testing.T) {
	upData := make([]Value, 20)
	downData := make([]Value, 20)
	for i := range upData {
		upData[i] = Float(100 + float64(i))
		downData[i] = Float(100 - float64(i))
	}

	rsiUp := CalcRSI(upData, RSIOptions{Periods: []int{6}})
	assertFloat(t, rsiUp[6]["rsi6"], 100)

	rsiDown := CalcRSI(downData, RSIOptions{Periods: []int{6}})
	assertFloat(t, rsiDown[6]["rsi6"], 0)
}

func TestCalcRSISeedsFirstValueWithWilderWindow(t *testing.T) {
	rsi := CalcRSI(Values(10, 12, 11, 13), RSIOptions{Periods: []int{2}})

	assertNil(t, rsi[0]["rsi2"])
	assertNil(t, rsi[1]["rsi2"])
	assertFloat(t, rsi[2]["rsi2"], 66.67)
	assertFloat(t, rsi[3]["rsi2"], 85.71)
}

func TestCalcWR(t *testing.T) {
	data := make([]OHLCV, 15)
	for i := range data {
		data[i] = OHLCV{
			Open:  Float(100),
			High:  Float(110),
			Low:   Float(90),
			Close: Float(100 + float64(i%3)*5 - 5),
		}
	}

	wr := CalcWR(data, WROptions{Periods: []int{6}})

	if len(wr) != 15 {
		t.Fatalf("len(wr) = %d, want 15", len(wr))
	}
	if wr[5]["wr6"] == nil {
		t.Fatal("wr[5].wr6 is nil")
	}
}

func TestCalcWRHandlesNullAndFlatRanges(t *testing.T) {
	withNull := []OHLCV{
		{Open: Float(100), High: Float(110), Low: Float(90), Close: Float(100)},
		{Open: Float(100), High: Null(), Low: Float(95), Close: Float(102)},
		{Open: Float(100), High: Float(115), Low: Null(), Close: Float(105)},
		{Open: Float(100), High: Float(120), Low: Float(100), Close: Null()},
		{Open: Float(100), High: Float(125), Low: Float(105), Close: Float(115)},
		{Open: Float(100), High: Float(130), Low: Float(110), Close: Float(120)},
		{Open: Float(100), High: Float(135), Low: Float(115), Close: Float(125)},
		{Open: Float(100), High: Float(140), Low: Float(120), Close: Float(130)},
		{Open: Float(100), High: Float(145), Low: Float(125), Close: Float(135)},
		{Open: Float(100), High: Float(150), Low: Float(130), Close: Float(140)},
	}
	wr := CalcWR(withNull, WROptions{Periods: []int{5}})
	assertNil(t, wr[4]["wr5"])

	flat := make([]OHLCV, 10)
	for i := range flat {
		flat[i] = OHLCV{Open: Float(100), High: Float(100), Low: Float(100), Close: Float(100)}
	}
	wr = CalcWR(flat, WROptions{Periods: []int{5}})
	assertNil(t, wr[4]["wr5"])
}
