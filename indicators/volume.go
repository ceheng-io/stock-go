package indicators

// CalcOBV calculates On Balance Volume.
func CalcOBV(data []OHLCV, options OBVOptions) []OBVResult {
	if len(data) == 0 {
		return nil
	}

	result := make([]OBVResult, 0, len(data))
	obvValue := 0.0
	if data[0].Volume != nil {
		obvValue = *data[0].Volume
	}
	result = append(result, OBVResult{OBV: Float(obvValue)})

	for i := 1; i < len(data); i++ {
		current := data[i]
		prev := data[i-1]
		if current.Close == nil || prev.Close == nil || current.Volume == nil {
			result = append(result, OBVResult{})
			continue
		}
		if *current.Close > *prev.Close {
			obvValue += *current.Volume
		} else if *current.Close < *prev.Close {
			obvValue -= *current.Volume
		}
		result = append(result, OBVResult{OBV: Float(obvValue)})
	}

	if options.MAPeriod > 0 {
		applyOBVMA(result, options.MAPeriod)
	}
	return result
}

func applyOBVMA(result []OBVResult, period int) {
	for i := period - 1; i < len(result); i++ {
		sum := 0.0
		count := 0
		for j := i - period + 1; j <= i; j++ {
			if result[j].OBV != nil {
				sum += *result[j].OBV
				count++
			}
		}
		if count == period {
			result[i].OBVMA = Float(round(sum / float64(period)))
		}
	}
}

// CalcROC calculates Rate of Change.
func CalcROC(data []OHLCV, options ROCOptions) []ROCResult {
	period := options.Period
	if period == 0 {
		period = 12
	}

	result := make([]ROCResult, 0, len(data))
	for i := range data {
		if i < period {
			result = append(result, ROCResult{})
			continue
		}
		current := data[i].Close
		previous := data[i-period].Close
		if current == nil || previous == nil || *previous == 0 {
			result = append(result, ROCResult{})
			continue
		}
		result = append(result, ROCResult{ROC: Float(round(((*current - *previous) / *previous) * 100))})
	}

	if options.SignalPeriod > 0 {
		applyROCSignal(result, period, options.SignalPeriod)
	}
	return result
}

func applyROCSignal(result []ROCResult, period int, signalPeriod int) {
	for i := period + signalPeriod - 1; i < len(result); i++ {
		sum := 0.0
		count := 0
		for j := i - signalPeriod + 1; j <= i; j++ {
			if result[j].ROC != nil {
				sum += *result[j].ROC
				count++
			}
		}
		if count == signalPeriod {
			result[i].Signal = Float(round(sum / float64(signalPeriod)))
		}
	}
}
