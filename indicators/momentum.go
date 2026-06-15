package indicators

import (
	"fmt"
	"math"
)

// CalcMACD calculates the MACD indicator.
func CalcMACD(closes []Value, options MACDOptions) []MACDResult {
	short := options.Short
	if short == 0 {
		short = 12
	}
	long := options.Long
	if long == 0 {
		long = 26
	}
	signal := options.Signal
	if signal == 0 {
		signal = 9
	}

	emaShort := CalcEMA(closes, short)
	emaLong := CalcEMA(closes, long)
	dif := make([]Value, len(closes))
	for i := range closes {
		if emaShort[i] == nil || emaLong[i] == nil {
			dif[i] = Null()
			continue
		}
		dif[i] = Float(*emaShort[i] - *emaLong[i])
	}

	dea := CalcEMA(dif, signal)
	result := make([]MACDResult, len(closes))
	for i := range closes {
		row := MACDResult{DEA: dea[i]}
		if dif[i] != nil {
			row.DIF = Float(round(*dif[i]))
		}
		if dif[i] != nil && dea[i] != nil {
			row.MACD = Float(round((*dif[i] - *dea[i]) * 2))
		}
		result[i] = row
	}
	return result
}

// CalcBOLL calculates Bollinger Bands.
func CalcBOLL(closes []Value, options BOLLOptions) []BOLLResult {
	period := options.Period
	if period == 0 {
		period = 20
	}
	stdDev := options.StdDev
	if stdDev == 0 {
		stdDev = 2
	}

	mid := CalcSMA(closes, period)
	std := calcStdDev(closes, period, mid)

	result := make([]BOLLResult, len(closes))
	for i := range closes {
		if mid[i] == nil || std[i] == nil {
			result[i] = BOLLResult{}
			continue
		}
		upper := *mid[i] + stdDev**std[i]
		lower := *mid[i] - stdDev**std[i]
		row := BOLLResult{
			Mid:   mid[i],
			Upper: Float(round(upper)),
			Lower: Float(round(lower)),
		}
		if *mid[i] != 0 {
			row.Bandwidth = Float(round(((upper - lower) / *mid[i]) * 100))
		}
		result[i] = row
	}
	return result
}

func calcStdDev(data []Value, period int, ma []Value) []Value {
	result := make([]Value, 0, len(data))
	for i := range data {
		if i < period-1 || ma[i] == nil {
			result = append(result, Null())
			continue
		}

		sumSquares := 0.0
		count := 0
		for j := i - period + 1; j <= i; j++ {
			if data[j] != nil {
				diff := *data[j] - *ma[i]
				sumSquares += math.Pow(diff, 2)
				count++
			}
		}
		if count == period {
			result = append(result, Float(math.Sqrt(sumSquares/float64(period))))
		} else {
			result = append(result, Null())
		}
	}
	return result
}

// CalcKDJ calculates KDJ.
func CalcKDJ(data []OHLCV, options KDJOptions) []KDJResult {
	period := options.Period
	if period == 0 {
		period = 9
	}
	kPeriod := options.KPeriod
	if kPeriod == 0 {
		kPeriod = 3
	}
	dPeriod := options.DPeriod
	if dPeriod == 0 {
		dPeriod = 3
	}

	result := make([]KDJResult, 0, len(data))
	k := 50.0
	d := 50.0
	for i := range data {
		if i < period-1 {
			result = append(result, KDJResult{})
			continue
		}

		highN := math.Inf(-1)
		lowN := math.Inf(1)
		valid := true
		for j := i - period + 1; j <= i; j++ {
			if data[j].High == nil || data[j].Low == nil {
				valid = false
				break
			}
			highN = math.Max(highN, *data[j].High)
			lowN = math.Min(lowN, *data[j].Low)
		}

		close := data[i].Close
		if !valid || close == nil || highN == lowN {
			result = append(result, KDJResult{})
			continue
		}

		rsv := ((*close - lowN) / (highN - lowN)) * 100
		k = (float64(kPeriod-1)/float64(kPeriod))*k + (1/float64(kPeriod))*rsv
		d = (float64(dPeriod-1)/float64(dPeriod))*d + (1/float64(dPeriod))*k
		j := 3*k - 2*d
		result = append(result, KDJResult{
			K: Float(round(k)),
			D: Float(round(d)),
			J: Float(round(j)),
		})
	}
	return result
}

// CalcRSI calculates RSI for one or more periods.
func CalcRSI(closes []Value, options RSIOptions) []RSIResult {
	periods := options.Periods
	if len(periods) == 0 {
		periods = []int{6, 12, 24}
	}

	changes := make([]Value, len(closes))
	if len(closes) > 0 {
		changes[0] = Null()
	}
	for i := 1; i < len(closes); i++ {
		if closes[i] == nil || closes[i-1] == nil {
			changes[i] = Null()
			continue
		}
		changes[i] = Float(*closes[i] - *closes[i-1])
	}

	series := make(map[int][]Value, len(periods))
	for _, period := range periods {
		rsi := make([]Value, 0, len(closes))
		avgGain := 0.0
		avgLoss := 0.0

		for i := range closes {
			if i < period {
				rsi = append(rsi, Null())
				if changes[i] != nil {
					if *changes[i] > 0 {
						avgGain += *changes[i]
					} else {
						avgLoss += math.Abs(*changes[i])
					}
				}
				continue
			}

			if i == period {
				if changes[i] != nil {
					if *changes[i] > 0 {
						avgGain += *changes[i]
					} else {
						avgLoss += math.Abs(*changes[i])
					}
				}
				avgGain /= float64(period)
				avgLoss /= float64(period)
			} else {
				change := 0.0
				if changes[i] != nil {
					change = *changes[i]
				}
				gain := 0.0
				loss := 0.0
				if change > 0 {
					gain = change
				}
				if change < 0 {
					loss = math.Abs(change)
				}
				avgGain = (avgGain*float64(period-1) + gain) / float64(period)
				avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
			}

			switch {
			case avgLoss == 0:
				rsi = append(rsi, Float(100))
			case avgGain == 0:
				rsi = append(rsi, Float(0))
			default:
				rs := avgGain / avgLoss
				rsi = append(rsi, Float(round(100-100/(1+rs))))
			}
		}
		series[period] = rsi
	}

	result := make([]RSIResult, len(closes))
	for i := range closes {
		row := make(RSIResult, len(periods))
		for _, period := range periods {
			row[fmt.Sprintf("rsi%d", period)] = series[period][i]
		}
		result[i] = row
	}
	return result
}

// CalcWR calculates Williams %R for one or more periods.
func CalcWR(data []OHLCV, options WROptions) []WRResult {
	periods := options.Periods
	if len(periods) == 0 {
		periods = []int{6, 10}
	}

	series := make(map[int][]Value, len(periods))
	for _, period := range periods {
		wr := make([]Value, 0, len(data))
		for i := range data {
			if i < period-1 {
				wr = append(wr, Null())
				continue
			}

			highN := math.Inf(-1)
			lowN := math.Inf(1)
			valid := true
			for j := i - period + 1; j <= i; j++ {
				if data[j].High == nil || data[j].Low == nil {
					valid = false
					break
				}
				highN = math.Max(highN, *data[j].High)
				lowN = math.Min(lowN, *data[j].Low)
			}

			close := data[i].Close
			if !valid || close == nil || highN == lowN {
				wr = append(wr, Null())
				continue
			}

			wr = append(wr, Float(round(((highN-*close)/(highN-lowN))*100)))
		}
		series[period] = wr
	}

	result := make([]WRResult, len(data))
	for i := range data {
		row := make(WRResult, len(periods))
		for _, period := range periods {
			row[fmt.Sprintf("wr%d", period)] = series[period][i]
		}
		result[i] = row
	}
	return result
}
