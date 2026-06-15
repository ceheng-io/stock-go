package indicators

import "math"

// Lookback describes how many historical bars are needed before the requested range.
type Lookback struct {
	MaxLookback          int
	HasEMABasedIndicator bool
	RequiredBars         int
}

// EstimateLookback returns the amount of historical context needed by enabled indicators.
func EstimateLookback(options Options) Lookback {
	lookback := Lookback{}

	if options.MA != nil {
		periods := options.MA.Periods
		if len(periods) == 0 {
			periods = []int{5, 10, 20, 30, 60, 120, 250}
		}
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxPeriod(periods, 20))
		lookback.HasEMABasedIndicator = lookback.HasEMABasedIndicator || options.MA.Type == MATypeEMA
	}
	if options.MACD != nil {
		long := defaultInt(options.MACD.Long, 26)
		signal := defaultInt(options.MACD.Signal, 9)
		lookback.MaxLookback = maxInt(lookback.MaxLookback, long*3+signal)
		lookback.HasEMABasedIndicator = true
	}
	if options.BOLL != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, defaultInt(options.BOLL.Period, 20))
	}
	if options.KDJ != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, defaultInt(options.KDJ.Period, 9))
	}
	if options.RSI != nil {
		periods := options.RSI.Periods
		if len(periods) == 0 {
			periods = []int{6, 12, 24}
		}
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxPeriod(periods, 14)+1)
	}
	if options.WR != nil {
		periods := options.WR.Periods
		if len(periods) == 0 {
			periods = []int{6, 10}
		}
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxPeriod(periods, 10))
	}
	if options.BIAS != nil {
		periods := options.BIAS.Periods
		if len(periods) == 0 {
			periods = []int{6, 12, 24}
		}
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxPeriod(periods, 12))
	}
	if options.CCI != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, defaultInt(options.CCI.Period, 14))
	}
	if options.ATR != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, defaultInt(options.ATR.Period, 14))
	}
	if options.OBV != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxInt(2, options.OBV.MAPeriod))
	}
	if options.ROC != nil {
		period := defaultInt(options.ROC.Period, 12)
		lookback.MaxLookback = maxInt(lookback.MaxLookback, period+options.ROC.SignalPeriod)
	}
	if options.DMI != nil {
		period := defaultInt(options.DMI.Period, 14)
		adxPeriod := defaultInt(options.DMI.ADXPeriod, period)
		lookback.MaxLookback = maxInt(lookback.MaxLookback, period*2+adxPeriod)
	}
	if options.SAR != nil {
		lookback.MaxLookback = maxInt(lookback.MaxLookback, 5)
	}
	if options.KC != nil {
		emaPeriod := defaultInt(options.KC.EMAPeriod, 20)
		atrPeriod := defaultInt(options.KC.ATRPeriod, 10)
		lookback.MaxLookback = maxInt(lookback.MaxLookback, maxInt(emaPeriod*3, atrPeriod))
		lookback.HasEMABasedIndicator = true
	}

	buffer := 1.2
	if lookback.HasEMABasedIndicator {
		buffer = 1.5
	}
	lookback.RequiredBars = int(math.Ceil(float64(lookback.MaxLookback) * buffer))
	return lookback
}

func defaultInt(value int, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}

func maxPeriod(values []int, fallback int) int {
	if len(values) == 0 {
		return fallback
	}
	result := values[0]
	for _, value := range values[1:] {
		if value > result {
			result = value
		}
	}
	return result
}
