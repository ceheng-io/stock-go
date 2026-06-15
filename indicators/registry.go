package indicators

// IndicatorKey identifies an enabled indicator in stable registry order.
type IndicatorKey string

const (
	IndicatorMA   IndicatorKey = "ma"
	IndicatorMACD IndicatorKey = "macd"
	IndicatorBOLL IndicatorKey = "boll"
	IndicatorKDJ  IndicatorKey = "kdj"
	IndicatorRSI  IndicatorKey = "rsi"
	IndicatorWR   IndicatorKey = "wr"
	IndicatorBIAS IndicatorKey = "bias"
	IndicatorCCI  IndicatorKey = "cci"
	IndicatorATR  IndicatorKey = "atr"
	IndicatorOBV  IndicatorKey = "obv"
	IndicatorROC  IndicatorKey = "roc"
	IndicatorDMI  IndicatorKey = "dmi"
	IndicatorSAR  IndicatorKey = "sar"
	IndicatorKC   IndicatorKey = "kc"
)

var indicatorKeyOrder = []IndicatorKey{
	IndicatorMA,
	IndicatorMACD,
	IndicatorBOLL,
	IndicatorKDJ,
	IndicatorRSI,
	IndicatorWR,
	IndicatorBIAS,
	IndicatorCCI,
	IndicatorATR,
	IndicatorOBV,
	IndicatorROC,
	IndicatorDMI,
	IndicatorSAR,
	IndicatorKC,
}

// Context stores reusable indicator inputs derived from K-line rows.
type Context struct {
	Closes []Value
	OHLCV  []OHLCV
}

// BuildContext extracts close prices and OHLCV bars from normalized K-line rows.
func BuildContext(rows []KlineInput) Context {
	context := Context{
		Closes: make([]Value, len(rows)),
		OHLCV:  make([]OHLCV, len(rows)),
	}
	for i, row := range rows {
		context.Closes[i] = row.Close
		context.OHLCV[i] = OHLCV{
			Open:   row.Open,
			High:   row.High,
			Low:    row.Low,
			Close:  row.Close,
			Volume: row.Volume,
		}
	}
	return context
}

// GetEnabledKeys returns enabled indicator keys in registry order.
func GetEnabledKeys(options Options) []IndicatorKey {
	keys := make([]IndicatorKey, 0, len(indicatorKeyOrder))
	for _, key := range indicatorKeyOrder {
		if isEnabled(options, key) {
			keys = append(keys, key)
		}
	}
	return keys
}

func isEnabled(options Options, key IndicatorKey) bool {
	switch key {
	case IndicatorMA:
		return options.MA != nil
	case IndicatorMACD:
		return options.MACD != nil
	case IndicatorBOLL:
		return options.BOLL != nil
	case IndicatorKDJ:
		return options.KDJ != nil
	case IndicatorRSI:
		return options.RSI != nil
	case IndicatorWR:
		return options.WR != nil
	case IndicatorBIAS:
		return options.BIAS != nil
	case IndicatorCCI:
		return options.CCI != nil
	case IndicatorATR:
		return options.ATR != nil
	case IndicatorOBV:
		return options.OBV != nil
	case IndicatorROC:
		return options.ROC != nil
	case IndicatorDMI:
		return options.DMI != nil
	case IndicatorSAR:
		return options.SAR != nil
	case IndicatorKC:
		return options.KC != nil
	default:
		return false
	}
}
