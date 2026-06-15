package signals

import (
	"fmt"

	"github.com/ceheng.io/stock-go/internal/core"
)

// CalcSignals detects indicator events from precomputed indicator rows.
func CalcSignals(klines []Kline, options SignalOptions) ([]Signal, error) {
	if err := validateSignalOptions(klines, options); err != nil {
		return nil, err
	}

	result := make([]Signal, 0)
	for i := 1; i < len(klines); i++ {
		prev := klines[i-1]
		cur := klines[i]
		if cur.Timestamp == nil {
			continue
		}
		at := *cur.Timestamp

		if options.MA != nil {
			result = append(result, maSignals(prev, cur, at, i, *options.MA)...)
		}
		if options.MACD {
			result = append(result, macdSignals(prev, cur, at, i)...)
		}
		if options.KDJ != nil {
			result = append(result, kdjSignals(prev, cur, at, i, *options.KDJ)...)
		}
		if options.RSI != nil {
			if signal := rsiSignal(cur, at, i, *options.RSI); signal != nil {
				result = append(result, *signal)
			}
		}
		if options.BOLL {
			if signal := bollSignal(cur, at, i); signal != nil {
				result = append(result, *signal)
			}
		}
		if options.SAR {
			if signal := sarSignal(prev, cur, at, i); signal != nil {
				result = append(result, *signal)
			}
		}
	}
	return result, nil
}

func validateSignalOptions(klines []Kline, options SignalOptions) error {
	if options.MA != nil && len(klines) > 0 {
		fastKey := fmt.Sprintf("ma%d", options.MA.Fast)
		slowKey := fmt.Sprintf("ma%d", options.MA.Slow)
		if !hasMAKeys(klines, fastKey, slowKey) {
			return invalidArgumentError(fmt.Sprintf("calcSignals: MA periods {fast:%d, slow:%d} not found on klines", options.MA.Fast, options.MA.Slow))
		}
	}
	if options.RSI != nil && len(klines) > 0 {
		period := rsiPeriod(*options.RSI)
		key := fmt.Sprintf("rsi%d", period)
		if !hasRSIKey(klines, key) {
			return invalidArgumentError(fmt.Sprintf("calcSignals: RSI period %d not found on klines", period))
		}
	}
	return nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

func hasMAKeys(klines []Kline, fastKey string, slowKey string) bool {
	for _, item := range klines {
		if item.MA == nil {
			continue
		}
		_, hasFast := item.MA[fastKey]
		_, hasSlow := item.MA[slowKey]
		if hasFast && hasSlow {
			return true
		}
	}
	return false
}

func hasRSIKey(klines []Kline, key string) bool {
	for _, item := range klines {
		if item.RSI == nil {
			continue
		}
		if _, ok := item.RSI[key]; ok {
			return true
		}
	}
	return false
}

func maSignals(prev Kline, cur Kline, at int64, index int, options MAOptions) []Signal {
	fastKey := fmt.Sprintf("ma%d", options.Fast)
	slowKey := fmt.Sprintf("ma%d", options.Slow)
	pf, ps := prev.MA[fastKey], prev.MA[slowKey]
	cf, cs := cur.MA[fastKey], cur.MA[slowKey]
	if pf == nil || ps == nil || cf == nil || cs == nil {
		return nil
	}
	detail := map[string]float64{"fast": float64(options.Fast), "slow": float64(options.Slow)}
	if *pf <= *ps && *cf > *cs {
		return []Signal{{Type: SignalMAGoldenCross, At: at, Index: index, Detail: detail}}
	}
	if *pf >= *ps && *cf < *cs {
		return []Signal{{Type: SignalMADeathCross, At: at, Index: index, Detail: detail}}
	}
	return nil
}

func macdSignals(prev Kline, cur Kline, at int64, index int) []Signal {
	if prev.MACD == nil || cur.MACD == nil || prev.MACD.DIF == nil || prev.MACD.DEA == nil || cur.MACD.DIF == nil || cur.MACD.DEA == nil {
		return nil
	}
	if *prev.MACD.DIF <= *prev.MACD.DEA && *cur.MACD.DIF > *cur.MACD.DEA {
		return []Signal{{Type: SignalMACDGoldenCross, At: at, Index: index}}
	}
	if *prev.MACD.DIF >= *prev.MACD.DEA && *cur.MACD.DIF < *cur.MACD.DEA {
		return []Signal{{Type: SignalMACDDeathCross, At: at, Index: index}}
	}
	return nil
}

func kdjSignals(prev Kline, cur Kline, at int64, index int, options KDJOptions) []Signal {
	overbought := optionOrDefault(options.Overbought, 80)
	oversold := optionOrDefault(options.Oversold, 20)
	result := make([]Signal, 0, 2)
	if prev.KDJ != nil && cur.KDJ != nil && prev.KDJ.K != nil && prev.KDJ.D != nil && cur.KDJ.K != nil && cur.KDJ.D != nil {
		if *prev.KDJ.K <= *prev.KDJ.D && *cur.KDJ.K > *cur.KDJ.D {
			result = append(result, Signal{Type: SignalKDJGoldenCross, At: at, Index: index})
		} else if *prev.KDJ.K >= *prev.KDJ.D && *cur.KDJ.K < *cur.KDJ.D {
			result = append(result, Signal{Type: SignalKDJDeathCross, At: at, Index: index})
		}
	}
	if cur.KDJ != nil && cur.KDJ.K != nil {
		if *cur.KDJ.K > overbought {
			result = append(result, Signal{Type: SignalKDJOverbought, At: at, Index: index, Detail: map[string]float64{"k": *cur.KDJ.K}})
		} else if *cur.KDJ.K < oversold {
			result = append(result, Signal{Type: SignalKDJOversold, At: at, Index: index, Detail: map[string]float64{"k": *cur.KDJ.K}})
		}
	}
	return result
}

func rsiSignal(cur Kline, at int64, index int, options RSIOptions) *Signal {
	period := rsiPeriod(options)
	value := cur.RSI[fmt.Sprintf("rsi%d", period)]
	if value == nil {
		return nil
	}
	overbought := optionOrDefault(options.Overbought, 70)
	oversold := optionOrDefault(options.Oversold, 30)
	if *value > overbought {
		return &Signal{Type: SignalRSIOverbought, At: at, Index: index, Detail: map[string]float64{"rsi": *value}}
	}
	if *value < oversold {
		return &Signal{Type: SignalRSIOversold, At: at, Index: index, Detail: map[string]float64{"rsi": *value}}
	}
	return nil
}

func bollSignal(cur Kline, at int64, index int) *Signal {
	if cur.BOLL == nil || cur.Close == nil {
		return nil
	}
	if cur.BOLL.Upper != nil && *cur.Close > *cur.BOLL.Upper {
		return &Signal{Type: SignalBOLLBreakUpper, At: at, Index: index}
	}
	if cur.BOLL.Lower != nil && *cur.Close < *cur.BOLL.Lower {
		return &Signal{Type: SignalBOLLBreakLower, At: at, Index: index}
	}
	return nil
}

func sarSignal(prev Kline, cur Kline, at int64, index int) *Signal {
	if prev.SAR == nil || cur.SAR == nil || prev.SAR.Trend == nil || cur.SAR.Trend == nil || *prev.SAR.Trend == *cur.SAR.Trend {
		return nil
	}
	if *cur.SAR.Trend == 1 {
		return &Signal{Type: SignalSARReversalUp, At: at, Index: index}
	}
	return &Signal{Type: SignalSARReversalDown, At: at, Index: index}
}

func rsiPeriod(options RSIOptions) int {
	if options.Period == 0 {
		return 6
	}
	return options.Period
}

func optionOrDefault(value float64, fallback float64) float64 {
	if value == 0 {
		return fallback
	}
	return value
}
