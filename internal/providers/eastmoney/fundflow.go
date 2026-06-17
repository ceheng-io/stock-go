package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ceheng-io/stock-go/symbols"
	"github.com/ceheng-io/stock-go/types"
)

const (
	emDataToken     = "b2884a393a59ad64002292a3e90d46a5"
	fflowFields2    = "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65"
	clistPageSize   = 100
	clistMaxPages   = 1000
	stockFundFlowFS = "m:0+t:6+f:!2,m:0+t:13+f:!2,m:0+t:80+f:!2,m:1+t:2+f:!2,m:1+t:23+f:!2,m:0+t:7+f:!2,m:1+t:3+f:!2"
)

// FundFlowPeriod is an Eastmoney fund-flow period.
type FundFlowPeriod string

const (
	FundFlowPeriodDaily   FundFlowPeriod = "daily"
	FundFlowPeriodWeekly  FundFlowPeriod = "weekly"
	FundFlowPeriodMonthly FundFlowPeriod = "monthly"
)

// FundFlowOptions configures fund-flow history fetching.
type FundFlowOptions struct {
	Period FundFlowPeriod
}

// FundFlowRankIndicator is an Eastmoney fund-flow rank window.
type FundFlowRankIndicator string

const (
	FundFlowRankToday    FundFlowRankIndicator = "today"
	FundFlowRankThreeDay FundFlowRankIndicator = "3day"
	FundFlowRankFiveDay  FundFlowRankIndicator = "5day"
	FundFlowRankTenDay   FundFlowRankIndicator = "10day"
)

// FundFlowSectorType is an Eastmoney sector type for fund-flow ranking.
type FundFlowSectorType string

const (
	FundFlowSectorIndustry FundFlowSectorType = "industry"
	FundFlowSectorConcept  FundFlowSectorType = "concept"
	FundFlowSectorRegion   FundFlowSectorType = "region"
)

// FundFlowRankOptions configures fund-flow ranking.
type FundFlowRankOptions struct {
	Indicator  FundFlowRankIndicator
	SectorType FundFlowSectorType
}

// FundFlowClient is the minimal client interface required by Eastmoney fund-flow providers.
type FundFlowClient interface {
	GetJSON(context.Context, string, any) error
}

type fundFlowResponse struct {
	Data struct {
		Klines json.RawMessage `json:"klines"`
	} `json:"data"`
}

type clistResponse struct {
	Data struct {
		Total int             `json:"total"`
		Diff  json.RawMessage `json:"diff"`
	} `json:"data"`
}

type rankIndicatorConfig struct {
	fid                string
	fields             string
	changePercentField string
	mainNet            string
	mainPct            string
	superLargeNet      string
	superLargePct      string
	largeNet           string
	largePct           string
	mediumNet          string
	mediumPct          string
	smallNet           string
	smallPct           string
}

var stockRankConfig = map[FundFlowRankIndicator]rankIndicatorConfig{
	FundFlowRankToday: {
		fid: "f62", fields: "f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f124", changePercentField: "f3",
		mainNet: "f62", mainPct: "f184", superLargeNet: "f66", superLargePct: "f69", largeNet: "f72", largePct: "f75", mediumNet: "f78", mediumPct: "f81", smallNet: "f84", smallPct: "f87",
	},
	FundFlowRankThreeDay: {
		fid: "f267", fields: "f12,f14,f2,f127,f267,f268,f269,f270,f271,f272,f273,f274,f275,f276,f124", changePercentField: "f127",
		mainNet: "f267", mainPct: "f268", superLargeNet: "f269", superLargePct: "f270", largeNet: "f271", largePct: "f272", mediumNet: "f273", mediumPct: "f274", smallNet: "f275", smallPct: "f276",
	},
	FundFlowRankFiveDay: {
		fid: "f164", fields: "f12,f14,f2,f109,f164,f165,f166,f167,f168,f169,f170,f171,f172,f173,f124", changePercentField: "f109",
		mainNet: "f164", mainPct: "f165", superLargeNet: "f166", superLargePct: "f167", largeNet: "f168", largePct: "f169", mediumNet: "f170", mediumPct: "f171", smallNet: "f172", smallPct: "f173",
	},
	FundFlowRankTenDay: {
		fid: "f174", fields: "f12,f14,f2,f160,f174,f175,f176,f177,f178,f179,f180,f181,f182,f183,f124", changePercentField: "f160",
		mainNet: "f174", mainPct: "f175", superLargeNet: "f176", superLargePct: "f177", largeNet: "f178", largePct: "f179", mediumNet: "f180", mediumPct: "f181", smallNet: "f182", smallPct: "f183",
	},
}

// GetIndividualFundFlow fetches stock fund-flow history rows.
func GetIndividualFundFlow(ctx context.Context, client FundFlowClient, symbol string, endpoint string, options FundFlowOptions) ([]types.StockFundFlow, error) {
	klt, err := fundFlowPeriodCode(options.Period)
	if err != nil {
		return nil, err
	}
	normalized, err := symbols.Normalize(symbol, &symbols.Hint{Market: symbols.MarketCN})
	if err != nil {
		return nil, err
	}
	secid, err := symbols.ToEastmoneySecIDE(normalized)
	if err != nil {
		return nil, err
	}
	return fetchStockFundFlow(ctx, client, endpoint, secid, klt)
}

// GetMarketFundFlow fetches market fund-flow rows for Shanghai and Shenzhen indexes.
func GetMarketFundFlow(ctx context.Context, client FundFlowClient, endpoint string) ([]types.MarketFundFlow, error) {
	params := fundFlowHistoryParams("1.000001", "101")
	params.Set("secid2", "0.399001")
	var payload fundFlowResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	klines, err := decodeStringArray(payload.Data.Klines)
	if err != nil {
		return nil, err
	}
	rows := make([]types.MarketFundFlow, 0, len(klines))
	for _, line := range klines {
		rows = append(rows, parseMarketFundFlow(line))
	}
	return rows, nil
}

// GetFundFlowRank fetches stock fund-flow ranking rows.
func GetFundFlowRank(ctx context.Context, client FundFlowClient, endpoint string, options FundFlowRankOptions) ([]types.FundFlowRankItem, error) {
	indicator := options.Indicator
	if indicator == "" {
		indicator = FundFlowRankToday
	}
	config, ok := stockRankConfig[indicator]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf("invalid fund flow rank indicator %q", indicator))
	}
	params := fundFlowRankParams(config, stockFundFlowFS, config.fields)
	items, err := fetchClistAllPages(ctx, client, endpoint, params, clistPageSize, clistMaxPages)
	if err != nil {
		return nil, err
	}
	rows := make([]types.FundFlowRankItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseFundFlowRankItem(item, config))
	}
	return rows, nil
}

// GetSectorFundFlowRank fetches sector fund-flow ranking rows.
func GetSectorFundFlowRank(ctx context.Context, client FundFlowClient, endpoint string, options FundFlowRankOptions) ([]types.SectorFundFlowItem, error) {
	indicator := options.Indicator
	if indicator == "" {
		indicator = FundFlowRankToday
	}
	config, ok := stockRankConfig[indicator]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf("invalid fund flow rank indicator %q", indicator))
	}
	sectorType := options.SectorType
	if sectorType == "" {
		sectorType = FundFlowSectorIndustry
	}
	fs, err := sectorFundFlowFS(sectorType)
	if err != nil {
		return nil, err
	}
	params := fundFlowRankParams(config, fs, config.fields+",f204,f205")
	items, err := fetchClistAllPages(ctx, client, endpoint, params, clistPageSize, clistMaxPages)
	if err != nil {
		return nil, err
	}
	rows := make([]types.SectorFundFlowItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseSectorFundFlowItem(item, config))
	}
	return rows, nil
}

// GetSectorFundFlowHistory fetches sector fund-flow history rows.
func GetSectorFundFlowHistory(ctx context.Context, client FundFlowClient, symbol string, endpoint string, options FundFlowOptions) ([]types.StockFundFlow, error) {
	klt, err := fundFlowPeriodCode(options.Period)
	if err != nil {
		return nil, err
	}
	secid := strings.TrimSpace(symbol)
	if !strings.Contains(secid, ".") {
		secid = "90." + strings.ToUpper(secid)
	}
	return fetchStockFundFlow(ctx, client, endpoint, secid, klt)
}

func fetchStockFundFlow(ctx context.Context, client FundFlowClient, endpoint string, secid string, klt string) ([]types.StockFundFlow, error) {
	params := fundFlowHistoryParams(secid, klt)
	var payload fundFlowResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	klines, err := decodeStringArray(payload.Data.Klines)
	if err != nil {
		return nil, err
	}
	rows := make([]types.StockFundFlow, 0, len(klines))
	for _, line := range klines {
		rows = append(rows, parseStockFundFlow(line))
	}
	return rows, nil
}

func fundFlowHistoryParams(secid string, klt string) url.Values {
	params := url.Values{}
	params.Set("lmt", "0")
	params.Set("klt", klt)
	params.Set("secid", secid)
	params.Set("fields1", "f1,f2,f3,f7")
	params.Set("fields2", fflowFields2)
	params.Set("ut", emDataToken)
	return params
}

func fundFlowPeriodCode(period FundFlowPeriod) (string, error) {
	switch period {
	case "", FundFlowPeriodDaily:
		return "101", nil
	case FundFlowPeriodWeekly:
		return "102", nil
	case FundFlowPeriodMonthly:
		return "103", nil
	default:
		return "", invalidArgumentError(fmt.Sprintf("invalid fund flow period %q", period))
	}
}

func fundFlowRankParams(config rankIndicatorConfig, fs string, fields string) url.Values {
	params := url.Values{}
	params.Set("fid", config.fid)
	params.Set("po", "1")
	params.Set("np", "1")
	params.Set("fltt", "2")
	params.Set("invt", "2")
	params.Set("ut", emDataToken)
	params.Set("fs", fs)
	params.Set("fields", fields)
	params.Set("pn", "1")
	params.Set("pz", "100")
	return params
}

func fetchClist(ctx context.Context, client FundFlowClient, endpoint string, params url.Values) ([]boardDynamicItem, error) {
	var payload clistResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	return decodeBoardDynamicArray(payload.Data.Diff)
}

func fetchClistAllPages(ctx context.Context, client FundFlowClient, endpoint string, baseParams url.Values, pageSize int, maxPages int) ([]boardDynamicItem, error) {
	if pageSize <= 0 {
		pageSize = clistPageSize
	}
	if maxPages <= 0 {
		maxPages = clistMaxPages
	}
	allItems := []boardDynamicItem{}
	total := 0
	for page := 1; page <= maxPages; page++ {
		params := cloneValues(baseParams)
		params.Set("pn", strconv.Itoa(page))
		params.Set("pz", strconv.Itoa(pageSize))

		var payload clistResponse
		if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
			return nil, err
		}
		if page == 1 {
			total = payload.Data.Total
		}
		diff, err := decodeBoardDynamicArray(payload.Data.Diff)
		if err != nil {
			return nil, err
		}
		if !isJSONArrayPayload(payload.Data.Diff) {
			break
		}
		allItems = append(allItems, diff...)
		if total <= 0 || len(allItems) >= total {
			break
		}
		if len(diff) < pageSize {
			break
		}
	}
	return allItems, nil
}

func cloneValues(values url.Values) url.Values {
	clone := make(url.Values, len(values))
	for key, items := range values {
		clone[key] = append([]string(nil), items...)
	}
	return clone
}

func sectorFundFlowFS(sectorType FundFlowSectorType) (string, error) {
	switch sectorType {
	case FundFlowSectorIndustry:
		return "m:90+t:2", nil
	case FundFlowSectorConcept:
		return "m:90+t:3", nil
	case FundFlowSectorRegion:
		return "m:90+t:1", nil
	default:
		return "", invalidArgumentError(fmt.Sprintf("invalid fund flow sector type %q", sectorType))
	}
}

func parseStockFundFlow(line string) types.StockFundFlow {
	fields := strings.Split(line, ",")
	return types.StockFundFlow{
		Date:                       boardField(fields, 0),
		MainNetInflow:              toNumber(boardField(fields, 1)),
		SmallNetInflow:             toNumber(boardField(fields, 2)),
		MediumNetInflow:            toNumber(boardField(fields, 3)),
		LargeNetInflow:             toNumber(boardField(fields, 4)),
		SuperLargeNetInflow:        toNumber(boardField(fields, 5)),
		MainNetInflowPercent:       toNumber(boardField(fields, 6)),
		SmallNetInflowPercent:      toNumber(boardField(fields, 7)),
		MediumNetInflowPercent:     toNumber(boardField(fields, 8)),
		LargeNetInflowPercent:      toNumber(boardField(fields, 9)),
		SuperLargeNetInflowPercent: toNumber(boardField(fields, 10)),
		Close:                      toNumber(boardField(fields, 11)),
		ChangePercent:              toNumber(boardField(fields, 12)),
	}
}

func parseMarketFundFlow(line string) types.MarketFundFlow {
	stock := parseStockFundFlow(line)
	fields := strings.Split(line, ",")
	return types.MarketFundFlow{
		Date:                       stock.Date,
		MainNetInflow:              stock.MainNetInflow,
		SmallNetInflow:             stock.SmallNetInflow,
		MediumNetInflow:            stock.MediumNetInflow,
		LargeNetInflow:             stock.LargeNetInflow,
		SuperLargeNetInflow:        stock.SuperLargeNetInflow,
		MainNetInflowPercent:       stock.MainNetInflowPercent,
		SmallNetInflowPercent:      stock.SmallNetInflowPercent,
		MediumNetInflowPercent:     stock.MediumNetInflowPercent,
		LargeNetInflowPercent:      stock.LargeNetInflowPercent,
		SuperLargeNetInflowPercent: stock.SuperLargeNetInflowPercent,
		SHClose:                    toNumber(boardField(fields, 11)),
		SHChangePercent:            toNumber(boardField(fields, 12)),
		SZClose:                    toNumber(boardField(fields, 13)),
		SZChangePercent:            toNumber(boardField(fields, 14)),
	}
}

func parseFundFlowRankItem(item boardDynamicItem, config rankIndicatorConfig) types.FundFlowRankItem {
	return types.FundFlowRankItem{
		Code:                       stringValue(item["f12"]),
		Name:                       stringValue(item["f14"]),
		Price:                      toNumberFromAny(item["f2"]),
		ChangePercent:              toNumberFromAny(item[config.changePercentField]),
		MainNetInflow:              toNumberFromAny(item[config.mainNet]),
		MainNetInflowPercent:       toNumberFromAny(item[config.mainPct]),
		SuperLargeNetInflow:        toNumberFromAny(item[config.superLargeNet]),
		SuperLargeNetInflowPercent: toNumberFromAny(item[config.superLargePct]),
		LargeNetInflow:             toNumberFromAny(item[config.largeNet]),
		LargeNetInflowPercent:      toNumberFromAny(item[config.largePct]),
		MediumNetInflow:            toNumberFromAny(item[config.mediumNet]),
		MediumNetInflowPercent:     toNumberFromAny(item[config.mediumPct]),
		SmallNetInflow:             toNumberFromAny(item[config.smallNet]),
		SmallNetInflowPercent:      toNumberFromAny(item[config.smallPct]),
	}
}

func parseSectorFundFlowItem(item boardDynamicItem, config rankIndicatorConfig) types.SectorFundFlowItem {
	return types.SectorFundFlowItem{
		Code:                 stringValue(item["f12"]),
		Name:                 stringValue(item["f14"]),
		ChangePercent:        toNumberFromAny(item[config.changePercentField]),
		MainNetInflow:        toNumberFromAny(item[config.mainNet]),
		MainNetInflowPercent: toNumberFromAny(item[config.mainPct]),
		SuperLargeNetInflow:  toNumberFromAny(item[config.superLargeNet]),
		LargeNetInflow:       toNumberFromAny(item[config.largeNet]),
		MediumNetInflow:      toNumberFromAny(item[config.mediumNet]),
		SmallNetInflow:       toNumberFromAny(item[config.smallNet]),
		TopStockCode:         nullableDatacenterString(item["f204"]),
		TopStockName:         nullableDatacenterString(item["f205"]),
	}
}
