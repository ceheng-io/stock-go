package indicators

import "math"

// CalcDMI calculates DMI/ADX.
func CalcDMI(data []OHLCV, options DMIOptions) []DMIResult {
	period := options.Period
	if period == 0 {
		period = 14
	}
	adxPeriod := options.ADXPeriod
	if adxPeriod == 0 {
		adxPeriod = period
	}

	result := make([]DMIResult, len(data))
	if len(data) < 2 {
		return result
	}

	plusDM := make([]float64, len(data))
	minusDM := make([]float64, len(data))
	tr := make([]float64, len(data))
	for i := 1; i < len(data); i++ {
		current := data[i]
		prev := data[i-1]
		if current.High == nil || current.Low == nil || current.Close == nil ||
			prev.High == nil || prev.Low == nil || prev.Close == nil {
			continue
		}
		hl := *current.High - *current.Low
		hc := math.Abs(*current.High - *prev.Close)
		lc := math.Abs(*current.Low - *prev.Close)
		tr[i] = math.Max(hl, math.Max(hc, lc))

		upMove := *current.High - *prev.High
		downMove := *prev.Low - *current.Low
		if upMove > downMove && upMove > 0 {
			plusDM[i] = upMove
		}
		if downMove > upMove && downMove > 0 {
			minusDM[i] = downMove
		}
	}

	dx := make([]float64, len(data))
	smoothTR := 0.0
	smoothPlusDM := 0.0
	smoothMinusDM := 0.0
	for i := 1; i < len(data); i++ {
		if i < period {
			smoothTR += tr[i]
			smoothPlusDM += plusDM[i]
			smoothMinusDM += minusDM[i]
			continue
		}

		if i == period {
			smoothTR += tr[i]
			smoothPlusDM += plusDM[i]
			smoothMinusDM += minusDM[i]
		} else {
			smoothTR = smoothTR - smoothTR/float64(period) + tr[i]
			smoothPlusDM = smoothPlusDM - smoothPlusDM/float64(period) + plusDM[i]
			smoothMinusDM = smoothMinusDM - smoothMinusDM/float64(period) + minusDM[i]
		}

		pdi := 0.0
		mdi := 0.0
		if smoothTR > 0 {
			pdi = smoothPlusDM / smoothTR * 100
			mdi = smoothMinusDM / smoothTR * 100
		}
		result[i].PDI = Float(round(pdi))
		result[i].MDI = Float(round(mdi))
		diSum := pdi + mdi
		if diSum > 0 {
			dx[i] = math.Abs(pdi-mdi) / diSum * 100
		}
	}

	adxSum := 0.0
	prevADX := 0.0
	for i := period; i < len(data); i++ {
		if i < period*2-1 {
			adxSum += dx[i]
			continue
		}
		if i == period*2-1 {
			adxSum += dx[i]
			prevADX = adxSum / float64(adxPeriod)
		} else {
			prevADX = (prevADX*float64(adxPeriod-1) + dx[i]) / float64(adxPeriod)
		}
		result[i].ADX = Float(round(prevADX))
	}

	for i := period*2 - 1 + adxPeriod; i < len(data); i++ {
		if result[i].ADX != nil && result[i-adxPeriod].ADX != nil {
			result[i].ADXR = Float(round((*result[i].ADX + *result[i-adxPeriod].ADX) / 2))
		}
	}
	return result
}

// CalcSAR calculates Parabolic SAR.
func CalcSAR(data []OHLCV, options SAROptions) []SARResult {
	result := make([]SARResult, len(data))
	if len(data) < 2 {
		return result
	}

	afStart := optionOrDefault(options.AFStart, 0.02)
	afIncrement := optionOrDefault(options.AFIncrement, 0.02)
	afMax := optionOrDefault(options.AFMax, 0.2)

	trend := 1
	af := afStart
	ep := valueOrZero(data[0].High)
	sar := valueOrZero(data[0].Low)
	if data[0].Close != nil && data[1].Close != nil && *data[1].Close < *data[0].Close {
		trend = -1
		ep = valueOrZero(data[0].Low)
		sar = valueOrZero(data[0].High)
	}

	for i := 1; i < len(data); i++ {
		current := data[i]
		prev := data[i-1]
		if current.High == nil || current.Low == nil || prev.High == nil || prev.Low == nil {
			continue
		}

		nextSAR := sar + af*(ep-sar)
		twoBack := data[max(0, i-2)]
		if trend == 1 {
			nextSAR = math.Min(nextSAR, *prev.Low)
			if twoBack.Low != nil {
				nextSAR = math.Min(nextSAR, *twoBack.Low)
			}
			if *current.Low < nextSAR {
				trend = -1
				nextSAR = ep
				ep = *current.Low
				af = afStart
			} else if *current.High > ep {
				ep = *current.High
				af = math.Min(af+afIncrement, afMax)
			}
		} else {
			nextSAR = math.Max(nextSAR, *prev.High)
			if twoBack.High != nil {
				nextSAR = math.Max(nextSAR, *twoBack.High)
			}
			if *current.High > nextSAR {
				trend = 1
				nextSAR = ep
				ep = *current.High
				af = afStart
			} else if *current.Low < ep {
				ep = *current.Low
				af = math.Min(af+afIncrement, afMax)
			}
		}

		sar = nextSAR
		trendValue := trend
		result[i] = SARResult{
			SAR:   Float(round(sar)),
			Trend: &trendValue,
			EP:    Float(round(ep)),
			AF:    Float(round(af)),
		}
	}
	return result
}

// CalcKC calculates Keltner Channel.
func CalcKC(data []OHLCV, options KCOptions) []KCResult {
	emaPeriod := options.EMAPeriod
	if emaPeriod == 0 {
		emaPeriod = 20
	}
	atrPeriod := options.ATRPeriod
	if atrPeriod == 0 {
		atrPeriod = 10
	}
	multiplier := optionOrDefault(options.Multiplier, 2)

	closes := make([]Value, len(data))
	for i, item := range data {
		closes[i] = item.Close
	}
	ema := CalcEMA(closes, emaPeriod)
	atr := CalcATR(data, ATROptions{Period: atrPeriod})

	result := make([]KCResult, len(data))
	for i := range data {
		if ema[i] == nil || atr[i].ATR == nil {
			continue
		}
		upper := *ema[i] + multiplier**atr[i].ATR
		lower := *ema[i] - multiplier**atr[i].ATR
		row := KCResult{
			Mid:   ema[i],
			Upper: Float(round(upper)),
			Lower: Float(round(lower)),
		}
		if *ema[i] > 0 {
			row.Width = Float(round(((upper - lower) / *ema[i]) * 100))
		}
		result[i] = row
	}
	return result
}

func optionOrDefault(value float64, fallback float64) float64 {
	if value == 0 {
		return fallback
	}
	return value
}

func valueOrZero(value Value) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func max(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
