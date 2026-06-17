package eastmoney

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ceheng-io/stock-go/types"
)

type NorthboundDirection = types.NorthboundDirection

const (
	NorthboundNorth = types.NorthboundNorth
	NorthboundSouth = types.NorthboundSouth
)

type NorthboundMarket = types.NorthboundMarket

const (
	NorthboundMarketAll      = types.NorthboundMarketAll
	NorthboundMarketShanghai = types.NorthboundMarketShanghai
	NorthboundMarketShenzhen = types.NorthboundMarketShenzhen
)

type NorthboundRankPeriod = types.NorthboundRankPeriod

const (
	NorthboundRankToday    = types.NorthboundRankToday
	NorthboundRankThreeDay = types.NorthboundRankThreeDay
	NorthboundRankFiveDay  = types.NorthboundRankFiveDay
	NorthboundRankTenDay   = types.NorthboundRankTenDay
	NorthboundRankMonth    = types.NorthboundRankMonth
	NorthboundRankQuarter  = types.NorthboundRankQuarter
	NorthboundRankYear     = types.NorthboundRankYear
)

// NorthboundHoldingRankOptions configures northbound holding ranking.
type NorthboundHoldingRankOptions struct {
	Market NorthboundMarket
	Period NorthboundRankPeriod
	Date   string
}

// NorthboundHistoryOptions configures northbound history queries.
type NorthboundHistoryOptions struct {
	StartDate string
	EndDate   string
}

// NorthboundClient is the minimal client interface required by northbound providers.
type NorthboundClient interface {
	GetJSON(context.Context, string, any) error
}

type northboundMinuteResponse struct {
	Data struct {
		S2N     []string `json:"s2n"`
		S2NDate string   `json:"s2nDate"`
		N2S     []string `json:"n2s"`
		N2SDate string   `json:"n2sDate"`
	} `json:"data"`
}

var northboundRankPeriodMap = map[NorthboundRankPeriod]string{
	NorthboundRankToday:    "1",
	NorthboundRankThreeDay: "3",
	NorthboundRankFiveDay:  "5",
	NorthboundRankTenDay:   "10",
	NorthboundRankMonth:    "M",
	NorthboundRankQuarter:  "Q",
	NorthboundRankYear:     "Y",
}

// GetNorthboundMinute fetches northbound or southbound intraday flow rows.
func GetNorthboundMinute(ctx context.Context, client NorthboundClient, endpoint string, direction NorthboundDirection) ([]types.NorthboundMinuteItem, error) {
	if direction == "" {
		direction = NorthboundNorth
	}
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4")
	params.Set("fields2", "f51,f54,f52,f58,f53,f62,f56,f57,f60,f61")
	params.Set("ut", emDataToken)

	var payload northboundMinuteResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	list := payload.Data.S2N
	date := formatNorthboundDate(payload.Data.S2NDate)
	if direction == NorthboundSouth {
		list = payload.Data.N2S
		date = formatNorthboundDate(payload.Data.N2SDate)
	}
	rows := make([]types.NorthboundMinuteItem, 0, len(list))
	for _, line := range list {
		rows = append(rows, parseNorthboundMinuteRow(line, date))
	}
	return rows, nil
}

// GetNorthboundFlowSummary fetches Shanghai/Shenzhen/HK connect flow summary rows.
func GetNorthboundFlowSummary(ctx context.Context, client NorthboundClient, endpoint string) ([]types.NorthboundFlowSummary, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:    "RPT_MUTUAL_QUOTA",
		columns:       "TRADE_DATE,MUTUAL_TYPE,BOARD_TYPE,MUTUAL_TYPE_NAME,FUNDS_DIRECTION,INDEX_CODE,INDEX_NAME,BOARD_CODE",
		quoteColumns:  "status~07~BOARD_CODE,dayNetAmtIn~07~BOARD_CODE,dayAmtRemain~07~BOARD_CODE,dayAmtThreshold~07~BOARD_CODE,f104~07~BOARD_CODE,f105~07~BOARD_CODE,f106~07~BOARD_CODE,f3~03~INDEX_CODE~INDEX_f3,netBuyAmt~07~BOARD_CODE",
		quoteType:     "0",
		sortColumns:   "MUTUAL_TYPE",
		sortTypes:     "1",
		pageSize:      "2000",
		fetchAllPages: datacenterBool(false),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.NorthboundFlowSummary, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseNorthboundFlowSummary(item))
	}
	return rows, nil
}

// GetNorthboundHoldingRank fetches northbound holding ranking rows.
func GetNorthboundHoldingRank(ctx context.Context, client NorthboundClient, endpoint string, options NorthboundHoldingRankOptions) ([]types.NorthboundHoldingRankItem, error) {
	period := options.Period
	if period == "" {
		period = NorthboundRankFiveDay
	}
	interval, ok := northboundRankPeriodMap[period]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf("invalid northbound rank period %q", period))
	}
	market := options.Market
	if market == "" {
		market = NorthboundMarketAll
	}
	filters := []string{fmt.Sprintf(`(INTERVAL_TYPE="%s")`, interval)}
	if options.Date != "" {
		filters = append(filters, fmt.Sprintf(`(TRADE_DATE='%s')`, options.Date))
	}
	if market != NorthboundMarketAll {
		code, err := northboundMarketCode(market)
		if err != nil {
			return nil, err
		}
		filters = append(filters, fmt.Sprintf(`(MUTUAL_TYPE="%s")`, code))
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_MUTUAL_STOCK_NORTHSTA",
		columns:     "ALL",
		sortColumns: "ADD_MARKET_CAP",
		sortTypes:   "-1",
		pageSize:    "500",
		filter:      strings.Join(filters, ""),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.NorthboundHoldingRankItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseNorthboundHoldingRank(item))
	}
	return rows, nil
}

// GetNorthboundHistory fetches northbound or southbound daily flow history rows.
func GetNorthboundHistory(ctx context.Context, client NorthboundClient, endpoint string, direction NorthboundDirection, options NorthboundHistoryOptions) ([]types.NorthboundHistoryItem, error) {
	filters := []string{`(BOARD_TYPE="1")`}
	if direction == NorthboundSouth {
		filters = []string{`(BOARD_TYPE="0")`}
	}
	if options.StartDate != "" {
		filters = append(filters, fmt.Sprintf(`(TRADE_DATE>='%s')`, options.StartDate))
	}
	if options.EndDate != "" {
		filters = append(filters, fmt.Sprintf(`(TRADE_DATE<='%s')`, options.EndDate))
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_MUTUAL_DEAL_HISTORY",
		columns:     "ALL",
		sortColumns: "TRADE_DATE",
		sortTypes:   "-1",
		pageSize:    "500",
		filter:      strings.Join(filters, ""),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.NorthboundHistoryItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseNorthboundHistory(item))
	}
	return rows, nil
}

// GetNorthboundIndividual fetches a stock's northbound holding history rows.
func GetNorthboundIndividual(ctx context.Context, client NorthboundClient, endpoint string, symbol string, options NorthboundHistoryOptions) ([]types.NorthboundIndividualItem, error) {
	code := strings.TrimSpace(symbol)
	lower := strings.ToLower(code)
	for _, prefix := range []string{"sh", "sz", "bj"} {
		if strings.HasPrefix(lower, prefix) {
			code = code[len(prefix):]
			break
		}
	}
	filters := []string{fmt.Sprintf(`(SECURITY_CODE="%s")`, code)}
	if options.StartDate != "" {
		filters = append(filters, fmt.Sprintf(`(TRADE_DATE>='%s')`, options.StartDate))
	}
	if options.EndDate != "" {
		filters = append(filters, fmt.Sprintf(`(TRADE_DATE<='%s')`, options.EndDate))
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_MUTUAL_HOLDSTOCKNORTH_STA",
		columns:     "ALL",
		sortColumns: "TRADE_DATE",
		sortTypes:   "-1",
		pageSize:    "500",
		filter:      strings.Join(filters, ""),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.NorthboundIndividualItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseNorthboundIndividual(item))
	}
	return rows, nil
}

func northboundMarketCode(market NorthboundMarket) (string, error) {
	switch market {
	case NorthboundMarketShanghai:
		return "001", nil
	case NorthboundMarketShenzhen:
		return "003", nil
	default:
		return "", invalidArgumentError(fmt.Sprintf("invalid northbound market %q", market))
	}
}

func parseNorthboundMinuteRow(line string, date string) types.NorthboundMinuteItem {
	fields := strings.Split(line, ",")
	return types.NorthboundMinuteItem{
		Date:              date,
		Time:              boardField(fields, 0),
		ShanghaiNetInflow: toNumber(boardField(fields, 1)),
		ShenzhenNetInflow: toNumber(boardField(fields, 3)),
		TotalNetInflow:    toNumber(boardField(fields, 5)),
	}
}

func formatNorthboundDate(value string) string {
	if len(value) == 8 && value[4] != '-' {
		return value[:4] + "-" + value[4:6] + "-" + value[6:8]
	}
	return value
}

func parseNorthboundFlowSummary(item boardDynamicItem) types.NorthboundFlowSummary {
	return types.NorthboundFlowSummary{
		Date:               parseDatacenterDate(item["TRADE_DATE"]),
		Type:               stringValue(item["MUTUAL_TYPE"]),
		BoardName:          stringValue(item["MUTUAL_TYPE_NAME"]),
		Direction:          stringValue(item["FUNDS_DIRECTION"]),
		Status:             stringValue(item["status"]),
		NetBuyAmount:       toNumberFromAny(item["netBuyAmt"]),
		NetInflow:          toNumberFromAny(item["dayNetAmtIn"]),
		RemainAmount:       toNumberFromAny(item["dayAmtRemain"]),
		UpCount:            toNumberFromAny(item["f104"]),
		FlatCount:          toNumberFromAny(item["f106"]),
		DownCount:          toNumberFromAny(item["f105"]),
		IndexCode:          stringValue(item["INDEX_CODE"]),
		IndexName:          stringValue(item["INDEX_NAME"]),
		IndexChangePercent: toNumberFromAny(item["INDEX_f3"]),
	}
}

func parseNorthboundHoldingRank(item boardDynamicItem) types.NorthboundHoldingRankItem {
	nameValue, ok := item["SECURITY_NAME"]
	name := stringValue(nameValue)
	if !ok || nameValue == nil {
		name = stringValue(item["SECURITY_NAME_ABBR"])
	}
	return types.NorthboundHoldingRankItem{
		Date:                  parseDatacenterDate(item["TRADE_DATE"]),
		Code:                  stringValue(item["SECURITY_CODE"]),
		Name:                  name,
		Close:                 toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent:         toNumberFromAny(item["CHANGE_RATE"]),
		HoldShares:            toNumberFromAny(item["HOLD_SHARES"]),
		HoldMarketValue:       toNumberFromAny(item["HOLD_MARKET_CAP"]),
		HoldRatioFloat:        toNumberFromAny(item["HOLD_RATIO"]),
		HoldRatioTotal:        toNumberFromAny(item["A_SHARES_RATIO"]),
		AddShares:             toNumberFromAny(item["ADD_SHARES"]),
		AddMarketValue:        toNumberFromAny(item["ADD_MARKET_CAP"]),
		AddMarketValuePercent: toNumberFromAny(item["ADD_MARKET_CAP_PROPORTION"]),
		Sector:                stringValue(item["BOARD_NAME"]),
	}
}

func parseNorthboundHistory(item boardDynamicItem) types.NorthboundHistoryItem {
	return types.NorthboundHistoryItem{
		Date:                  parseDatacenterDate(item["TRADE_DATE"]),
		NetBuyAmount:          toNumberFromAny(item["NET_DEAL_AMT"]),
		BuyAmount:             toNumberFromAny(item["BUY_AMT"]),
		SellAmount:            toNumberFromAny(item["SELL_AMT"]),
		AccNetBuyAmount:       toNumberFromAny(item["ACCUM_DEAL_AMT"]),
		NetInflow:             toNumberFromAny(item["FUND_INFLOW"]),
		RemainAmount:          toNumberFromAny(item["QUOTA_BALANCE"]),
		TopStockCode:          nullableDatacenterString(item["LEAD_STOCKS_CODE"]),
		TopStockName:          nullableDatacenterString(item["LEAD_STOCKS_NAME"]),
		TopStockChangePercent: toNumberFromAny(item["LS_CHANGE_RATE"]),
	}
}

func parseNorthboundIndividual(item boardDynamicItem) types.NorthboundIndividualItem {
	return types.NorthboundIndividualItem{
		Date:            parseDatacenterDate(item["TRADE_DATE"]),
		HoldShares:      toNumberFromAny(item["HOLD_SHARES"]),
		HoldMarketValue: toNumberFromAny(item["HOLD_MARKET_CAP"]),
		HoldRatioFloat:  toNumberFromAny(item["HOLD_RATIO"]),
		HoldRatioTotal:  toNumberFromAny(item["A_SHARES_RATIO"]),
		Close:           toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent:   toNumberFromAny(item["CHANGE_RATE"]),
	}
}
