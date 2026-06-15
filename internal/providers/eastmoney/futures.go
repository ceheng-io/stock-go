package eastmoney

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

// FuturesKlineOptions configures domestic futures historical K-line fetching.
type FuturesKlineOptions struct {
	Period    KlinePeriod
	StartDate string
	EndDate   string
}

// GetFuturesHistoryKline fetches domestic futures historical K-line rows.
func GetFuturesHistoryKline(ctx context.Context, client KlineClient, symbol string, endpoint string, options FuturesKlineOptions) ([]types.FuturesKline, error) {
	options = normalizeFuturesKlineOptions(options)
	if err := validateFuturesKlineOptions(options); err != nil {
		return nil, err
	}
	variety, err := extractFuturesVariety(symbol)
	if err != nil {
		return nil, err
	}
	marketCode, err := futuresMarketCode(variety)
	if err != nil {
		return nil, err
	}
	secid := fmt.Sprintf("%d.%s", marketCode, strings.TrimSpace(symbol))

	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64")
	params.Set("ut", emPushToken)
	params.Set("klt", periodCode(options.Period))
	params.Set("fqt", "0")
	params.Set("secid", secid)
	params.Set("beg", options.StartDate)
	params.Set("end", options.EndDate)

	var payload historyKlineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	code := payload.Data.Code
	if code == "" {
		code = strings.TrimSpace(symbol)
	}
	rows := make([]types.FuturesKline, 0, len(payload.Data.Klines))
	for _, line := range payload.Data.Klines {
		rows = append(rows, parseFuturesKlineCSV(line, code, payload.Data.Name))
	}
	return rows, nil
}

func normalizeFuturesKlineOptions(options FuturesKlineOptions) FuturesKlineOptions {
	if options.Period == "" {
		options.Period = KlinePeriodDaily
	}
	if options.StartDate == "" {
		options.StartDate = "19700101"
	}
	if options.EndDate == "" {
		options.EndDate = "20500101"
	}
	return options
}

func validateFuturesKlineOptions(options FuturesKlineOptions) error {
	switch options.Period {
	case KlinePeriodDaily, KlinePeriodWeekly, KlinePeriodMonthly:
		return nil
	default:
		return invalidArgumentError(fmt.Sprintf("invalid period %q", options.Period))
	}
}

func parseFuturesKlineCSV(line string, code string, name string) types.FuturesKline {
	fields := strings.Split(line, ",")
	return types.FuturesKline{
		Date:          field(fields, 0),
		Code:          code,
		Name:          name,
		Open:          toNumber(field(fields, 1)),
		Close:         toNumber(field(fields, 2)),
		High:          toNumber(field(fields, 3)),
		Low:           toNumber(field(fields, 4)),
		Volume:        toNumber(field(fields, 5)),
		Amount:        toNumber(field(fields, 6)),
		Amplitude:     toNumber(field(fields, 7)),
		ChangePercent: toNumber(field(fields, 8)),
		Change:        toNumber(field(fields, 9)),
		TurnoverRate:  toNumber(field(fields, 10)),
		OpenInterest:  toNumber(field(fields, 12)),
	}
}

func extractFuturesVariety(symbol string) (string, error) {
	match := regexp.MustCompile(`^([a-zA-Z]+)`).FindStringSubmatch(strings.TrimSpace(symbol))
	if len(match) < 2 {
		return "", invalidArgumentError(fmt.Sprintf("invalid futures symbol %q", symbol))
	}
	return match[1], nil
}

func futuresMarketCode(variety string) (int, error) {
	exchange, ok := futuresVarietyExchange[variety]
	if !ok {
		exchange, ok = futuresVarietyExchange[strings.ToLower(variety)]
	}
	if !ok {
		exchange, ok = futuresVarietyExchange[strings.ToUpper(variety)]
	}
	if !ok && len(variety) > 1 && strings.HasSuffix(variety, "M") {
		exchange, ok = futuresVarietyExchange[variety[:len(variety)-1]]
		if !ok {
			exchange, ok = futuresVarietyExchange[strings.ToLower(variety[:len(variety)-1])]
		}
	}
	if !ok {
		return 0, invalidArgumentError(fmt.Sprintf("unknown futures variety %q", variety))
	}
	code, ok := futuresExchangeMarketCode[exchange]
	if !ok {
		return 0, invalidArgumentError(fmt.Sprintf("no market code found for exchange %q", exchange))
	}
	return code, nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

func invalidSymbolError(symbol string) error {
	return core.NewCodedError("INVALID_SYMBOL", fmt.Sprintf("Invalid symbol: %s", symbol), nil)
}

var futuresExchangeMarketCode = map[string]int{
	"SHFE":  113,
	"DCE":   114,
	"CZCE":  115,
	"INE":   142,
	"CFFEX": 220,
	"GFEX":  225,
}

var futuresVarietyExchange = map[string]string{
	"cu": "SHFE", "al": "SHFE", "zn": "SHFE", "pb": "SHFE", "au": "SHFE", "ag": "SHFE",
	"rb": "SHFE", "wr": "SHFE", "fu": "SHFE", "ru": "SHFE", "bu": "SHFE", "hc": "SHFE",
	"ni": "SHFE", "sn": "SHFE", "sp": "SHFE", "ss": "SHFE", "ao": "SHFE", "br": "SHFE",
	"c": "DCE", "a": "DCE", "b": "DCE", "m": "DCE", "y": "DCE", "p": "DCE",
	"l": "DCE", "v": "DCE", "j": "DCE", "jm": "DCE", "i": "DCE", "jd": "DCE",
	"pp": "DCE", "cs": "DCE", "eg": "DCE", "eb": "DCE", "pg": "DCE", "lh": "DCE",
	"WH": "CZCE", "CF": "CZCE", "SR": "CZCE", "TA": "CZCE", "OI": "CZCE", "MA": "CZCE",
	"FG": "CZCE", "RM": "CZCE", "SF": "CZCE", "SM": "CZCE", "ZC": "CZCE", "AP": "CZCE",
	"CJ": "CZCE", "UR": "CZCE", "SA": "CZCE", "PF": "CZCE", "PK": "CZCE", "PX": "CZCE",
	"SH": "CZCE",
	"sc": "INE", "nr": "INE", "lu": "INE", "bc": "INE", "ec": "INE",
	"IF": "CFFEX", "IC": "CFFEX", "IH": "CFFEX", "IM": "CFFEX",
	"TS": "CFFEX", "TF": "CFFEX", "T": "CFFEX", "TL": "CFFEX",
	"si": "GFEX", "lc": "GFEX", "ps": "GFEX", "pt": "GFEX", "pd": "GFEX",
}
