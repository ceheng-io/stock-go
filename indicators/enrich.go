package indicators

// Options configures which indicators should be attached to K-line rows.
type Options struct {
	MA   *MAOptions
	MACD *MACDOptions
	BOLL *BOLLOptions
	KDJ  *KDJOptions
	RSI  *RSIOptions
	WR   *WROptions
	BIAS *BIASOptions
	CCI  *CCIOptions
	ATR  *ATROptions
	OBV  *OBVOptions
	ROC  *ROCOptions
	DMI  *DMIOptions
	SAR  *SAROptions
	KC   *KCOptions
}

// KlineInput is the normalized input used by indicator enrichment.
type KlineInput struct {
	Date   string
	Open   Value
	High   Value
	Low    Value
	Close  Value
	Volume Value
}

// KlineWithIndicators contains one K-line row plus enabled indicator outputs.
type KlineWithIndicators struct {
	Date   string
	Open   Value
	High   Value
	Low    Value
	Close  Value
	Volume Value

	MA   MAResult
	MACD *MACDResult
	BOLL *BOLLResult
	KDJ  *KDJResult
	RSI  RSIResult
	WR   WRResult
	BIAS BIASResult
	CCI  *CCIResult
	ATR  *ATRResult
	OBV  *OBVResult
	ROC  *ROCResult
	DMI  *DMIResult
	SAR  *SARResult
	KC   *KCResult
}

// AddIndicators attaches enabled technical indicators to normalized K-line rows.
func AddIndicators(rows []KlineInput, options Options) []KlineWithIndicators {
	if len(rows) == 0 {
		return nil
	}

	context := BuildContext(rows)

	result := make([]KlineWithIndicators, len(rows))
	for i, row := range rows {
		result[i] = KlineWithIndicators{
			Date:   row.Date,
			Open:   row.Open,
			High:   row.High,
			Low:    row.Low,
			Close:  row.Close,
			Volume: row.Volume,
		}
	}

	if options.MA != nil {
		values := CalcMA(context.Closes, *options.MA)
		for i := range result {
			result[i].MA = values[i]
		}
	}
	if options.MACD != nil {
		values := CalcMACD(context.Closes, *options.MACD)
		for i := range result {
			value := values[i]
			result[i].MACD = &value
		}
	}
	if options.BOLL != nil {
		values := CalcBOLL(context.Closes, *options.BOLL)
		for i := range result {
			value := values[i]
			result[i].BOLL = &value
		}
	}
	if options.KDJ != nil {
		values := CalcKDJ(context.OHLCV, *options.KDJ)
		for i := range result {
			value := values[i]
			result[i].KDJ = &value
		}
	}
	if options.RSI != nil {
		values := CalcRSI(context.Closes, *options.RSI)
		for i := range result {
			result[i].RSI = values[i]
		}
	}
	if options.WR != nil {
		values := CalcWR(context.OHLCV, *options.WR)
		for i := range result {
			result[i].WR = values[i]
		}
	}
	if options.BIAS != nil {
		values := CalcBIAS(context.Closes, *options.BIAS)
		for i := range result {
			result[i].BIAS = values[i]
		}
	}
	if options.CCI != nil {
		values := CalcCCI(context.OHLCV, *options.CCI)
		for i := range result {
			value := values[i]
			result[i].CCI = &value
		}
	}
	if options.ATR != nil {
		values := CalcATR(context.OHLCV, *options.ATR)
		for i := range result {
			value := values[i]
			result[i].ATR = &value
		}
	}
	if options.OBV != nil {
		values := CalcOBV(context.OHLCV, *options.OBV)
		for i := range result {
			value := values[i]
			result[i].OBV = &value
		}
	}
	if options.ROC != nil {
		values := CalcROC(context.OHLCV, *options.ROC)
		for i := range result {
			value := values[i]
			result[i].ROC = &value
		}
	}
	if options.DMI != nil {
		values := CalcDMI(context.OHLCV, *options.DMI)
		for i := range result {
			value := values[i]
			result[i].DMI = &value
		}
	}
	if options.SAR != nil {
		values := CalcSAR(context.OHLCV, *options.SAR)
		for i := range result {
			value := values[i]
			result[i].SAR = &value
		}
	}
	if options.KC != nil {
		values := CalcKC(context.OHLCV, *options.KC)
		for i := range result {
			value := values[i]
			result[i].KC = &value
		}
	}

	return result
}
