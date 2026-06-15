package indicators

import (
	"fmt"
	"math"
)

// CalcBIAS calculates BIAS for one or more periods.
func CalcBIAS(closes []Value, options BIASOptions) []BIASResult {
	periods := options.Periods
	if len(periods) == 0 {
		periods = []int{6, 12, 24}
	}

	series := make(map[int][]Value, len(periods))
	for _, period := range periods {
		ma := CalcSMA(closes, period)
		bias := make([]Value, len(closes))
		for i := range closes {
			if closes[i] == nil || ma[i] == nil || *ma[i] == 0 {
				bias[i] = Null()
				continue
			}
			bias[i] = Float(round(((*closes[i] - *ma[i]) / *ma[i]) * 100))
		}
		series[period] = bias
	}

	result := make([]BIASResult, len(closes))
	for i := range closes {
		row := make(BIASResult, len(periods))
		for _, period := range periods {
			row[fmt.Sprintf("bias%d", period)] = series[period][i]
		}
		result[i] = row
	}
	return result
}

// CalcCCI calculates Commodity Channel Index.
func CalcCCI(data []OHLCV, options CCIOptions) []CCIResult {
	period := options.Period
	if period == 0 {
		period = 14
	}

	tp := make([]Value, len(data))
	for i, item := range data {
		if item.High == nil || item.Low == nil || item.Close == nil {
			continue
		}
		tp[i] = Float((*item.High + *item.Low + *item.Close) / 3)
	}

	result := make([]CCIResult, 0, len(data))
	for i := range data {
		if i < period-1 {
			result = append(result, CCIResult{})
			continue
		}

		sum := 0.0
		count := 0
		for j := i - period + 1; j <= i; j++ {
			if tp[j] != nil {
				sum += *tp[j]
				count++
			}
		}
		if count != period || tp[i] == nil {
			result = append(result, CCIResult{})
			continue
		}

		ma := sum / float64(period)
		mdSum := 0.0
		for j := i - period + 1; j <= i; j++ {
			mdSum += math.Abs(*tp[j] - ma)
		}
		md := mdSum / float64(period)
		if md == 0 {
			result = append(result, CCIResult{CCI: Float(0)})
			continue
		}
		result = append(result, CCIResult{CCI: Float(round((*tp[i] - ma) / (0.015 * md)))})
	}
	return result
}

// CalcATR calculates True Range and Average True Range.
func CalcATR(data []OHLCV, options ATROptions) []ATRResult {
	period := options.Period
	if period == 0 {
		period = 14
	}

	tr := make([]Value, len(data))
	for i, item := range data {
		if item.High == nil || item.Low == nil || item.Close == nil {
			continue
		}
		if i == 0 || data[i-1].Close == nil {
			tr[i] = Float(*item.High - *item.Low)
			continue
		}
		hl := *item.High - *item.Low
		hpc := math.Abs(*item.High - *data[i-1].Close)
		lpc := math.Abs(*item.Low - *data[i-1].Close)
		tr[i] = Float(math.Max(hl, math.Max(hpc, lpc)))
	}

	result := make([]ATRResult, 0, len(data))
	var atr *float64
	for i := range data {
		if i < period-1 {
			result = append(result, ATRResult{TR: roundedValue(tr[i])})
			continue
		}

		switch {
		case i == period-1:
			sum := 0.0
			count := 0
			for j := 0; j < period; j++ {
				if tr[j] != nil {
					sum += *tr[j]
					count++
				}
			}
			if count == period {
				atr = Float(sum / float64(period))
			}
		case atr != nil && tr[i] != nil:
			atr = Float((*atr*float64(period-1) + *tr[i]) / float64(period))
		}

		result = append(result, ATRResult{
			TR:  roundedValue(tr[i]),
			ATR: roundedValue(atr),
		})
	}
	return result
}

func roundedValue(value Value) Value {
	if value == nil {
		return nil
	}
	return Float(round(*value))
}
