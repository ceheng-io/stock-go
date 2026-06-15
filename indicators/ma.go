package indicators

import (
	"fmt"
	"math"
)

func round(value float64) float64 {
	factor := math.Pow(10, 2)
	return math.Round(value*factor) / factor
}

// CalcSMA calculates simple moving average.
func CalcSMA(data []Value, period int) []Value {
	result := make([]Value, 0, len(data))
	if period <= 0 {
		for range data {
			result = append(result, Null())
		}
		return result
	}

	for i := range data {
		if i < period-1 {
			result = append(result, Null())
			continue
		}

		sum := 0.0
		count := 0
		for j := i - period + 1; j <= i; j++ {
			if data[j] != nil {
				sum += *data[j]
				count++
			}
		}

		if count == period {
			result = append(result, Float(round(sum/float64(period))))
		} else {
			result = append(result, Null())
		}
	}

	return result
}

// CalcEMA calculates exponential moving average.
func CalcEMA(data []Value, period int) []Value {
	result := make([]Value, 0, len(data))
	if period <= 0 {
		for range data {
			result = append(result, Null())
		}
		return result
	}

	alpha := 2 / float64(period+1)
	var ema float64
	initialized := false

	for i := range data {
		if i < period-1 {
			result = append(result, Null())
			continue
		}

		if !initialized {
			sum := 0.0
			count := 0
			for j := i - period + 1; j <= i; j++ {
				if data[j] != nil {
					sum += *data[j]
					count++
				}
			}
			if count == period {
				ema = sum / float64(period)
				initialized = true
				result = append(result, Float(round(ema)))
			} else {
				result = append(result, Null())
			}
			continue
		}

		if data[i] == nil {
			result = append(result, Float(round(ema)))
			continue
		}
		ema = alpha**data[i] + (1-alpha)*ema
		result = append(result, Float(round(ema)))
	}

	return result
}

// CalcWMA calculates weighted moving average.
func CalcWMA(data []Value, period int) []Value {
	result := make([]Value, 0, len(data))
	if period <= 0 {
		for range data {
			result = append(result, Null())
		}
		return result
	}

	weightSum := period * (period + 1) / 2
	for i := range data {
		if i < period-1 {
			result = append(result, Null())
			continue
		}

		sum := 0.0
		valid := true
		for j := 0; j < period; j++ {
			value := data[i-period+1+j]
			if value == nil {
				valid = false
				break
			}
			sum += *value * float64(j+1)
		}

		if valid {
			result = append(result, Float(round(sum/float64(weightSum))))
		} else {
			result = append(result, Null())
		}
	}

	return result
}

// CalcMA calculates moving averages for one or more periods.
func CalcMA(closes []Value, options MAOptions) []MAResult {
	periods := options.Periods
	if len(periods) == 0 {
		periods = []int{5, 10, 20, 30, 60, 120, 250}
	}

	calcFn := CalcSMA
	switch options.Type {
	case MATypeEMA:
		calcFn = CalcEMA
	case MATypeWMA:
		calcFn = CalcWMA
	}

	series := make(map[int][]Value, len(periods))
	for _, period := range periods {
		series[period] = calcFn(closes, period)
	}

	result := make([]MAResult, len(closes))
	for i := range closes {
		row := make(MAResult, len(periods))
		for _, period := range periods {
			row[fmt.Sprintf("ma%d", period)] = series[period][i]
		}
		result[i] = row
	}
	return result
}
