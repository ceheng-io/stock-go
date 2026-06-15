package indicators

// Value represents a nullable numeric input or indicator output.
type Value = *float64

// MAType controls the moving average algorithm.
type MAType string

const (
	MATypeSMA MAType = "sma"
	MATypeEMA MAType = "ema"
	MATypeWMA MAType = "wma"
)

// MAOptions configures moving average calculation.
type MAOptions struct {
	Periods []int
	Type    MAType
}

// MAResult stores moving average values by key, for example ma5 or ma20.
type MAResult map[string]Value

// OHLCV is the common indicator input bar.
type OHLCV struct {
	Open   Value
	High   Value
	Low    Value
	Close  Value
	Volume Value
}

// MACDOptions configures MACD calculation.
type MACDOptions struct {
	Short  int
	Long   int
	Signal int
}

// MACDResult contains DIF, DEA, and MACD bar values.
type MACDResult struct {
	DIF  Value
	DEA  Value
	MACD Value
}

// BOLLOptions configures Bollinger Bands calculation.
type BOLLOptions struct {
	Period int
	StdDev float64
}

// BOLLResult contains Bollinger Band values.
type BOLLResult struct {
	Mid       Value
	Upper     Value
	Lower     Value
	Bandwidth Value
}

// KDJOptions configures KDJ calculation.
type KDJOptions struct {
	Period  int
	KPeriod int
	DPeriod int
}

// KDJResult contains K, D, and J values.
type KDJResult struct {
	K Value
	D Value
	J Value
}

// RSIOptions configures RSI calculation.
type RSIOptions struct {
	Periods []int
}

// RSIResult stores RSI values by key, for example rsi6.
type RSIResult map[string]Value

// WROptions configures Williams %R calculation.
type WROptions struct {
	Periods []int
}

// WRResult stores Williams %R values by key, for example wr6.
type WRResult map[string]Value

// BIASOptions configures BIAS calculation.
type BIASOptions struct {
	Periods []int
}

// BIASResult stores BIAS values by key, for example bias6.
type BIASResult map[string]Value

// CCIOptions configures CCI calculation.
type CCIOptions struct {
	Period int
}

// CCIResult contains Commodity Channel Index values.
type CCIResult struct {
	CCI Value
}

// ATROptions configures ATR calculation.
type ATROptions struct {
	Period int
}

// ATRResult contains True Range and Average True Range values.
type ATRResult struct {
	TR  Value
	ATR Value
}

// OBVOptions configures OBV calculation.
type OBVOptions struct {
	MAPeriod int
}

// OBVResult contains On Balance Volume values.
type OBVResult struct {
	OBV   Value
	OBVMA Value
}

// ROCOptions configures ROC calculation.
type ROCOptions struct {
	Period       int
	SignalPeriod int
}

// ROCResult contains Rate of Change values.
type ROCResult struct {
	ROC    Value
	Signal Value
}

// DMIOptions configures DMI/ADX calculation.
type DMIOptions struct {
	Period    int
	ADXPeriod int
}

// DMIResult contains Directional Movement Index values.
type DMIResult struct {
	PDI  Value
	MDI  Value
	ADX  Value
	ADXR Value
}

// SAROptions configures Parabolic SAR calculation.
type SAROptions struct {
	AFStart     float64
	AFIncrement float64
	AFMax       float64
}

// SARResult contains Parabolic SAR values.
type SARResult struct {
	SAR   Value
	Trend *int
	EP    Value
	AF    Value
}

// KCOptions configures Keltner Channel calculation.
type KCOptions struct {
	EMAPeriod  int
	ATRPeriod  int
	Multiplier float64
}

// KCResult contains Keltner Channel values.
type KCResult struct {
	Mid   Value
	Upper Value
	Lower Value
	Width Value
}

// Float returns a non-null Value.
func Float(value float64) Value {
	return &value
}

// Null returns a null Value.
func Null() Value {
	return nil
}

// Values converts float64 inputs to non-null Values.
func Values(values ...float64) []Value {
	result := make([]Value, len(values))
	for i, value := range values {
		result[i] = Float(value)
	}
	return result
}
