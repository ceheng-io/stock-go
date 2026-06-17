package stock

import "github.com/ceheng-io/stock-go/signals"

type SignalType = signals.SignalType

const (
	SignalMAGoldenCross   = signals.SignalMAGoldenCross
	SignalMADeathCross    = signals.SignalMADeathCross
	SignalMACDGoldenCross = signals.SignalMACDGoldenCross
	SignalMACDDeathCross  = signals.SignalMACDDeathCross
	SignalKDJGoldenCross  = signals.SignalKDJGoldenCross
	SignalKDJDeathCross   = signals.SignalKDJDeathCross
	SignalKDJOverbought   = signals.SignalKDJOverbought
	SignalKDJOversold     = signals.SignalKDJOversold
	SignalRSIOverbought   = signals.SignalRSIOverbought
	SignalRSIOversold     = signals.SignalRSIOversold
	SignalBOLLBreakUpper  = signals.SignalBOLLBreakUpper
	SignalBOLLBreakLower  = signals.SignalBOLLBreakLower
	SignalSARReversalUp   = signals.SignalSARReversalUp
	SignalSARReversalDown = signals.SignalSARReversalDown
)

type SignalKline = signals.Kline
type SignalMACD = signals.MACD
type SignalKDJ = signals.KDJ
type SignalBOLL = signals.BOLL
type SignalSAR = signals.SAR
type Signal = signals.Signal
type SignalOptions = signals.SignalOptions
type SignalMAOptions = signals.MAOptions
type SignalKDJOptions = signals.KDJOptions
type SignalRSIOptions = signals.RSIOptions

func CalcSignals(klines []SignalKline, options SignalOptions) ([]Signal, error) {
	return signals.CalcSignals(klines, options)
}

func SignalFloat(value float64) *float64 {
	return signals.Float(value)
}

func SignalTime(value int64) *int64 {
	return signals.FloatTime(value)
}

func SignalInt(value int) *int {
	return signals.Int(value)
}
