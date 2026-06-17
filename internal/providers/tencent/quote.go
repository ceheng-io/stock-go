package tencent

import (
	"context"
	"strconv"

	"github.com/ceheng-io/stock-go/internal/core"
	"github.com/ceheng-io/stock-go/timeutil"
	"github.com/ceheng-io/stock-go/types"
)

// QuoteClient is the minimal client interface required by Tencent quote providers.
type QuoteClient interface {
	GetTencentQuote(context.Context, string) ([]core.TencentQuoteItem, error)
}

// GetSimpleQuotes fetches and parses CN simple quotes.
func GetSimpleQuotes(ctx context.Context, client QuoteClient, codes []string) ([]types.SimpleQuote, error) {
	if len(codes) == 0 {
		return []types.SimpleQuote{}, nil
	}

	prefixedCodes := make([]string, len(codes))
	wanted := make(map[string]struct{}, len(codes))
	for i, code := range codes {
		prefixed := "s_" + code
		prefixedCodes[i] = prefixed
		wanted[prefixed] = struct{}{}
	}

	items, err := client.GetTencentQuote(ctx, joinComma(prefixedCodes))
	if err != nil {
		return nil, err
	}

	quotes := make([]types.SimpleQuote, 0, len(items))
	for _, item := range items {
		if _, ok := wanted[item.Key]; !ok {
			continue
		}
		if len(item.Fields) <= 5 || item.Fields[0] == "" {
			continue
		}
		quotes = append(quotes, parseSimpleQuote(item.Fields))
	}
	return quotes, nil
}

// GetFullQuotes fetches and parses CN full quotes.
func GetFullQuotes(ctx context.Context, client QuoteClient, codes []string) ([]types.FullQuote, error) {
	if len(codes) == 0 {
		return []types.FullQuote{}, nil
	}

	wanted := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		wanted[code] = struct{}{}
	}

	items, err := client.GetTencentQuote(ctx, joinComma(codes))
	if err != nil {
		return nil, err
	}

	quotes := make([]types.FullQuote, 0, len(items))
	for _, item := range items {
		if _, ok := wanted[item.Key]; !ok {
			continue
		}
		if len(item.Fields) <= 5 || item.Fields[0] == "" {
			continue
		}
		quotes = append(quotes, parseFullQuote(item.Fields))
	}
	return quotes, nil
}

func parseFullQuote(f []string) types.FullQuote {
	timeMeta := timeutil.BuildTimeMeta(field(f, 30), timeutil.MarketTZ.CN)
	bid := make([]types.PriceLevel, 0, 5)
	for i := 0; i < 5; i++ {
		bid = append(bid, types.PriceLevel{
			Price:  safeNumber(field(f, 9+i*2)),
			Volume: safeNumber(field(f, 10+i*2)),
		})
	}
	ask := make([]types.PriceLevel, 0, 5)
	for i := 0; i < 5; i++ {
		ask = append(ask, types.PriceLevel{
			Price:  safeNumber(field(f, 19+i*2)),
			Volume: safeNumber(field(f, 20+i*2)),
		})
	}

	return types.FullQuote{
		MarketID:             field(f, 0),
		Name:                 field(f, 1),
		Code:                 field(f, 2),
		Price:                safeNumber(field(f, 3)),
		PrevClose:            safeNumber(field(f, 4)),
		Open:                 safeNumber(field(f, 5)),
		Volume:               safeNumber(field(f, 6)),
		OuterVolume:          safeNumber(field(f, 7)),
		InnerVolume:          safeNumber(field(f, 8)),
		Bid:                  bid,
		Ask:                  ask,
		Time:                 field(f, 30),
		Timestamp:            timeMeta.Timestamp,
		TZ:                   string(timeMeta.TZ),
		Change:               safeNumber(field(f, 31)),
		ChangePercent:        safeNumber(field(f, 32)),
		High:                 safeNumber(field(f, 33)),
		Low:                  safeNumber(field(f, 34)),
		Volume2:              safeNumber(field(f, 36)),
		Amount:               safeNumber(field(f, 37)),
		TurnoverRate:         safeNumberPtr(field(f, 38)),
		PE:                   safeNumberPtr(field(f, 39)),
		Amplitude:            safeNumberPtr(field(f, 43)),
		CirculatingMarketCap: safeNumberPtr(field(f, 44)),
		TotalMarketCap:       safeNumberPtr(field(f, 45)),
		PB:                   safeNumberPtr(field(f, 46)),
		LimitUp:              safeNumberPtr(field(f, 47)),
		LimitDown:            safeNumberPtr(field(f, 48)),
		VolumeRatio:          safeNumberPtr(field(f, 49)),
		AvgPrice:             safeNumberPtr(field(f, 51)),
		PEStatic:             safeNumberPtr(field(f, 52)),
		PEDynamic:            safeNumberPtr(field(f, 53)),
		High52W:              safeNumberPtr(field(f, 67)),
		Low52W:               safeNumberPtr(field(f, 68)),
		CirculatingShares:    safeNumberPtr(field(f, 72)),
		TotalShares:          safeNumberPtr(field(f, 73)),
		Market:               types.MarketCN,
		AssetType:            "stock",
		Source:               "tencent",
	}
}

// GetHKQuotes fetches and parses Hong Kong quotes.
func GetHKQuotes(ctx context.Context, client QuoteClient, codes []string) ([]types.HKQuote, error) {
	return getPrefixedQuotes(ctx, client, codes, "hk", 6, parseHKQuote)
}

// GetUSQuotes fetches and parses US quotes.
func GetUSQuotes(ctx context.Context, client QuoteClient, codes []string) ([]types.USQuote, error) {
	return getPrefixedQuotes(ctx, client, codes, "us", 6, parseUSQuote)
}

// GetFundQuotes fetches and parses public fund quotes.
func GetFundQuotes(ctx context.Context, client QuoteClient, codes []string) ([]types.FundQuote, error) {
	return getPrefixedQuotes(ctx, client, codes, "jj", 9, parseFundQuote)
}

// GetFundFlow fetches Tencent fund-flow rows.
func GetFundFlow(ctx context.Context, client QuoteClient, codes []string) ([]types.FundFlow, error) {
	return getPrefixedQuotes(ctx, client, codes, "ff_", 14, parseFundFlow)
}

// GetPanelLargeOrder fetches Tencent panel large-order ratio rows.
func GetPanelLargeOrder(ctx context.Context, client QuoteClient, codes []string) ([]types.PanelLargeOrder, error) {
	return getPrefixedQuotes(ctx, client, codes, "s_pk", 4, parsePanelLargeOrder)
}

func getPrefixedQuotes[T any](
	ctx context.Context,
	client QuoteClient,
	codes []string,
	prefix string,
	minFields int,
	parse func([]string) T,
) ([]T, error) {
	if len(codes) == 0 {
		return []T{}, nil
	}

	prefixedCodes := make([]string, len(codes))
	wanted := make(map[string]struct{}, len(codes))
	for i, code := range codes {
		prefixed := prefix + code
		prefixedCodes[i] = prefixed
		wanted[prefixed] = struct{}{}
	}

	items, err := client.GetTencentQuote(ctx, joinComma(prefixedCodes))
	if err != nil {
		return nil, err
	}

	quotes := make([]T, 0, len(items))
	for _, item := range items {
		if _, ok := wanted[item.Key]; !ok {
			continue
		}
		if len(item.Fields) < minFields || item.Fields[0] == "" {
			continue
		}
		quotes = append(quotes, parse(item.Fields))
	}
	return quotes, nil
}

func parseHKQuote(f []string) types.HKQuote {
	timeMeta := timeutil.BuildTimeMeta(field(f, 30), timeutil.MarketTZ.HK)
	return types.HKQuote{
		MarketID:             field(f, 0),
		Name:                 field(f, 1),
		Code:                 field(f, 2),
		Price:                safeNumber(field(f, 3)),
		PrevClose:            safeNumber(field(f, 4)),
		Open:                 safeNumber(field(f, 5)),
		Volume:               safeNumber(field(f, 6)),
		Time:                 field(f, 30),
		Timestamp:            timeMeta.Timestamp,
		TZ:                   string(timeMeta.TZ),
		Change:               safeNumber(field(f, 31)),
		ChangePercent:        safeNumber(field(f, 32)),
		High:                 safeNumber(field(f, 33)),
		Low:                  safeNumber(field(f, 34)),
		Amount:               safeNumber(field(f, 37)),
		LotSize:              safeNumberPtr(field(f, 40)),
		CirculatingMarketCap: safeNumberPtr(field(f, 44)),
		TotalMarketCap:       safeNumberPtr(field(f, 45)),
		Currency:             field(f, len(f)-3),
		Market:               types.MarketHK,
		AssetType:            "stock",
		Source:               "tencent",
	}
}

func parseUSQuote(f []string) types.USQuote {
	timeMeta := timeutil.BuildTimeMeta(field(f, 30), timeutil.MarketTZ.US)
	return types.USQuote{
		MarketID:       field(f, 0),
		Name:           field(f, 1),
		Code:           field(f, 2),
		Price:          safeNumber(field(f, 3)),
		PrevClose:      safeNumber(field(f, 4)),
		Open:           safeNumber(field(f, 5)),
		Volume:         safeNumber(field(f, 6)),
		Time:           field(f, 30),
		Timestamp:      timeMeta.Timestamp,
		TZ:             string(timeMeta.TZ),
		Change:         safeNumber(field(f, 31)),
		ChangePercent:  safeNumber(field(f, 32)),
		High:           safeNumber(field(f, 33)),
		Low:            safeNumber(field(f, 34)),
		Amount:         safeNumber(field(f, 37)),
		TurnoverRate:   safeNumberPtr(field(f, 38)),
		PE:             safeNumberPtr(field(f, 39)),
		Amplitude:      safeNumberPtr(field(f, 43)),
		TotalMarketCap: safeNumberPtr(field(f, 45)),
		PB:             safeNumberPtr(field(f, 47)),
		High52W:        safeNumberPtr(field(f, 48)),
		Low52W:         safeNumberPtr(field(f, 49)),
		Market:         types.MarketUS,
		AssetType:      "stock",
		Source:         "tencent",
	}
}

func parseFundQuote(f []string) types.FundQuote {
	timeMeta := timeutil.BuildTimeMeta(field(f, 8), timeutil.MarketTZ.CN)
	return types.FundQuote{
		Code:      field(f, 0),
		Name:      field(f, 1),
		NAV:       safeNumber(field(f, 5)),
		AccNAV:    safeNumber(field(f, 6)),
		Change:    safeNumber(field(f, 7)),
		NavDate:   field(f, 8),
		Timestamp: timeMeta.Timestamp,
		TZ:        string(timeMeta.TZ),
		Market:    types.MarketCN,
		AssetType: "fund",
		Source:    "tencent",
	}
}

func parseFundFlow(f []string) types.FundFlow {
	timeMeta := timeutil.BuildTimeMeta(field(f, 13), timeutil.MarketTZ.CN)
	return types.FundFlow{
		Code:           field(f, 0),
		MainInflow:     safeNumber(field(f, 1)),
		MainOutflow:    safeNumber(field(f, 2)),
		MainNet:        safeNumber(field(f, 3)),
		MainNetRatio:   safeNumber(field(f, 4)),
		RetailInflow:   safeNumber(field(f, 5)),
		RetailOutflow:  safeNumber(field(f, 6)),
		RetailNet:      safeNumber(field(f, 7)),
		RetailNetRatio: safeNumber(field(f, 8)),
		TotalFlow:      safeNumber(field(f, 9)),
		Name:           field(f, 12),
		Date:           field(f, 13),
		Timestamp:      timeMeta.Timestamp,
		TZ:             string(timeMeta.TZ),
	}
}

func parsePanelLargeOrder(f []string) types.PanelLargeOrder {
	return types.PanelLargeOrder{
		BuyLargeRatio:  safeNumber(field(f, 0)),
		BuySmallRatio:  safeNumber(field(f, 1)),
		SellLargeRatio: safeNumber(field(f, 2)),
		SellSmallRatio: safeNumber(field(f, 3)),
	}
}

func parseSimpleQuote(f []string) types.SimpleQuote {
	return types.SimpleQuote{
		MarketID:      field(f, 0),
		Name:          field(f, 1),
		Code:          field(f, 2),
		Price:         safeNumber(field(f, 3)),
		Change:        safeNumber(field(f, 4)),
		ChangePercent: safeNumber(field(f, 5)),
		Volume:        safeNumber(field(f, 6)),
		Amount:        safeNumber(field(f, 7)),
		MarketCap:     safeNumberPtr(field(f, 9)),
		MarketType:    field(f, 10),
		Market:        types.MarketCN,
		AssetType:     "stock",
		Source:        "tencent",
	}
}

func field(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return fields[index]
}

func safeNumber(value string) float64 {
	if value == "" {
		return 0
	}
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return n
}

func safeNumberPtr(value string) *float64 {
	if value == "" {
		return nil
	}
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &n
}

func joinComma(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for _, value := range values[1:] {
		result += "," + value
	}
	return result
}
