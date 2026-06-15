package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/ceheng.io/stock-go/types"
)

const emFuturesGlobalSpotToken = "58b2fa8f54638b60b87d69b31969089c"

// GlobalFuturesSpotOptions configures global futures spot fetching.
type GlobalFuturesSpotOptions struct {
	PageSize int
}

// GlobalFuturesKlineOptions configures global futures historical K-line fetching.
type GlobalFuturesKlineOptions struct {
	Period     KlinePeriod
	StartDate  string
	EndDate    string
	MarketCode int
}

type globalFuturesSpotResponse struct {
	List  json.RawMessage `json:"list"`
	Total int             `json:"total"`
}

// GetGlobalFuturesSpot fetches all global futures spot quote rows.
func GetGlobalFuturesSpot(ctx context.Context, client KlineClient, endpoint string, options GlobalFuturesSpotOptions) ([]types.GlobalFuturesQuote, error) {
	pageSize := options.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	rows := make([]types.GlobalFuturesQuote, 0)
	total := 0
	for pageIndex := 0; ; pageIndex++ {
		params := url.Values{}
		params.Set("orderBy", "dm")
		params.Set("sort", "desc")
		params.Set("pageSize", strconv.Itoa(pageSize))
		params.Set("pageIndex", strconv.Itoa(pageIndex))
		params.Set("token", emFuturesGlobalSpotToken)
		params.Set("field", "dm,sc,name,p,zsjd,zde,zdf,f152,o,h,l,zjsj,vol,wp,np,ccl")
		params.Set("blockName", "callback")

		var payload globalFuturesSpotResponse
		if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
			return nil, err
		}
		list, err := decodeMapArray(payload.List)
		if err != nil {
			return nil, err
		}
		if !isJSONArrayPayload(payload.List) {
			break
		}
		if pageIndex == 0 {
			total = payload.Total
		}
		if len(list) == 0 {
			break
		}
		for _, item := range list {
			rows = append(rows, parseGlobalFuturesQuote(item))
		}
		if total <= 0 || len(rows) >= total {
			break
		}
	}
	return rows, nil
}

// GetGlobalFuturesKline fetches global futures historical K-line rows.
func GetGlobalFuturesKline(ctx context.Context, client KlineClient, symbol string, endpoint string, options GlobalFuturesKlineOptions) ([]types.FuturesKline, error) {
	options = normalizeGlobalFuturesKlineOptions(options)
	if err := validateFuturesKlineOptions(FuturesKlineOptions{Period: options.Period}); err != nil {
		return nil, err
	}
	marketCode := options.MarketCode
	if marketCode == 0 {
		variety, err := extractGlobalFuturesVariety(symbol)
		if err != nil {
			return nil, err
		}
		var ok bool
		marketCode, ok = globalFuturesMarketCode[variety]
		if !ok {
			return nil, invalidArgumentError(fmt.Sprintf("unknown global futures variety %q", variety))
		}
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

func normalizeGlobalFuturesKlineOptions(options GlobalFuturesKlineOptions) GlobalFuturesKlineOptions {
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

func parseGlobalFuturesQuote(item map[string]any) types.GlobalFuturesQuote {
	return types.GlobalFuturesQuote{
		Code:          stringValue(item["dm"]),
		Name:          stringValue(item["name"]),
		Price:         toNumberFromAny(item["p"]),
		Change:        toNumberFromAny(item["zde"]),
		ChangePercent: toNumberFromAny(item["zdf"]),
		Open:          toNumberFromAny(item["o"]),
		High:          toNumberFromAny(item["h"]),
		Low:           toNumberFromAny(item["l"]),
		PrevSettle:    toNumberFromAny(item["zjsj"]),
		Volume:        toNumberFromAny(item["vol"]),
		BuyVolume:     toNumberFromAny(item["wp"]),
		SellVolume:    toNumberFromAny(item["np"]),
		OpenInterest:  toNumberFromAny(item["ccl"]),
	}
}

func extractGlobalFuturesVariety(symbol string) (string, error) {
	match := regexp.MustCompile(`^([A-Z]+)`).FindStringSubmatch(strings.TrimSpace(symbol))
	if len(match) < 2 {
		return "", invalidArgumentError(fmt.Sprintf("invalid global futures symbol %q", symbol))
	}
	return match[1], nil
}

var globalFuturesMarketCode = map[string]int{
	"HG": 101, "GC": 101, "SI": 101, "QI": 101, "QO": 101, "MGC": 101,
	"CL": 102, "NG": 102, "RB": 102, "HO": 102, "PA": 102, "PL": 102,
	"ZW": 103, "ZM": 103, "ZS": 103, "ZC": 103, "ZL": 103, "ZR": 103,
	"YM": 103, "NQ": 103, "ES": 103,
	"SB": 108, "CT": 108,
	"LCPT": 109, "LZNT": 109, "LALT": 109,
}
