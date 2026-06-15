package signals

// SignalType identifies an indicator event.
type SignalType string

const (
	SignalMAGoldenCross   SignalType = "ma_golden_cross"
	SignalMADeathCross    SignalType = "ma_death_cross"
	SignalMACDGoldenCross SignalType = "macd_golden_cross"
	SignalMACDDeathCross  SignalType = "macd_death_cross"
	SignalKDJGoldenCross  SignalType = "kdj_golden_cross"
	SignalKDJDeathCross   SignalType = "kdj_death_cross"
	SignalKDJOverbought   SignalType = "kdj_overbought"
	SignalKDJOversold     SignalType = "kdj_oversold"
	SignalRSIOverbought   SignalType = "rsi_overbought"
	SignalRSIOversold     SignalType = "rsi_oversold"
	SignalBOLLBreakUpper  SignalType = "boll_break_upper"
	SignalBOLLBreakLower  SignalType = "boll_break_lower"
	SignalSARReversalUp   SignalType = "sar_reversal_up"
	SignalSARReversalDown SignalType = "sar_reversal_down"
)

// Kline is a minimal K-line plus indicator row used by CalcSignals.
type Kline struct {
	Timestamp *int64
	Close     *float64
	MA        map[string]*float64
	MACD      *MACD
	KDJ       *KDJ
	RSI       map[string]*float64
	BOLL      *BOLL
	SAR       *SAR
}

// MACD contains signal-layer MACD values.
type MACD struct {
	DIF *float64
	DEA *float64
}

// KDJ contains signal-layer KDJ values.
type KDJ struct {
	K *float64
	D *float64
}

// BOLL contains signal-layer Bollinger Band values.
type BOLL struct {
	Upper *float64
	Lower *float64
}

// SAR contains signal-layer Parabolic SAR values.
type SAR struct {
	Trend *int
}

// Signal is one indicator event.
type Signal struct {
	Type   SignalType
	At     int64
	Index  int
	Detail map[string]float64
}

// SignalOptions controls which signal families are detected.
type SignalOptions struct {
	MA   *MAOptions
	MACD bool
	KDJ  *KDJOptions
	RSI  *RSIOptions
	BOLL bool
	SAR  bool
}

// MAOptions configures MA cross signals.
type MAOptions struct {
	Fast int
	Slow int
}

// KDJOptions configures KDJ signals.
type KDJOptions struct {
	Overbought float64
	Oversold   float64
}

// RSIOptions configures RSI signals.
type RSIOptions struct {
	Period     int
	Overbought float64
	Oversold   float64
}

// Float returns a non-null float pointer.
func Float(value float64) *float64 {
	return &value
}

// FloatTime returns a non-null timestamp pointer.
func FloatTime(value int64) *int64 {
	return &value
}

// Int returns a non-null int pointer.
func Int(value int) *int {
	return &value
}
