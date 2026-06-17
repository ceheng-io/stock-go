package eastmoney

import (
	"context"
	"fmt"
	"strings"

	"github.com/ceheng-io/stock-go/types"
)

type DragonTigerPeriod = types.DragonTigerPeriod

const (
	DragonTigerPeriodOneMonth   = types.DragonTigerPeriodOneMonth
	DragonTigerPeriodThreeMonth = types.DragonTigerPeriodThreeMonth
	DragonTigerPeriodSixMonth   = types.DragonTigerPeriodSixMonth
	DragonTigerPeriodOneYear    = types.DragonTigerPeriodOneYear
)

type DragonTigerDateOptions = types.DragonTigerDateOptions

// DragonTigerClient is the minimal client interface required by dragon-tiger providers.
type DragonTigerClient interface {
	GetJSON(context.Context, string, any) error
}

var dragonTigerPeriodMap = map[DragonTigerPeriod]string{
	DragonTigerPeriodOneMonth:   "01",
	DragonTigerPeriodThreeMonth: "02",
	DragonTigerPeriodSixMonth:   "03",
	DragonTigerPeriodOneYear:    "04",
}

// GetDragonTigerDetail fetches dragon-tiger billboard detail rows.
func GetDragonTigerDetail(ctx context.Context, client DragonTigerClient, endpoint string, options DragonTigerDateOptions) ([]types.DragonTigerDetailItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_DAILYBILLBOARD_DETAILSNEW",
		columns:     "SECURITY_CODE,SECUCODE,SECURITY_NAME_ABBR,TRADE_DATE,EXPLAIN,CLOSE_PRICE,CHANGE_RATE,BILLBOARD_NET_AMT,BILLBOARD_BUY_AMT,BILLBOARD_SELL_AMT,BILLBOARD_DEAL_AMT,ACCUM_AMOUNT,DEAL_NET_RATIO,DEAL_AMOUNT_RATIO,TURNOVERRATE,FREE_MARKET_CAP,EXPLANATION,D1_CLOSE_ADJCHRATE,D2_CLOSE_ADJCHRATE,D5_CLOSE_ADJCHRATE,D10_CLOSE_ADJCHRATE",
		sortColumns: "SECURITY_CODE,TRADE_DATE",
		sortTypes:   "1,-1",
		pageSize:    "5000",
		filter:      dragonTigerDateFilter(options),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DragonTigerDetailItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseDragonTigerDetail(item))
	}
	return rows, nil
}

// GetDragonTigerStockStats fetches stock billboard statistics rows.
func GetDragonTigerStockStats(ctx context.Context, client DragonTigerClient, endpoint string, period DragonTigerPeriod) ([]types.DragonTigerStockStatItem, error) {
	cycle, err := dragonTigerCycle(period)
	if err != nil {
		return nil, err
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BILLBOARD_TRADEALL",
		columns:     "ALL",
		sortColumns: "BILLBOARD_TIMES,LATEST_TDATE,SECURITY_CODE",
		sortTypes:   "-1,-1,1",
		pageSize:    "5000",
		filter:      fmt.Sprintf(`(STATISTICS_CYCLE="%s")`, cycle),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DragonTigerStockStatItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseDragonTigerStockStat(item))
	}
	return rows, nil
}

// GetDragonTigerInstitution fetches institution trading rows.
func GetDragonTigerInstitution(ctx context.Context, client DragonTigerClient, endpoint string, options DragonTigerDateOptions) ([]types.DragonTigerInstitutionItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_ORGANIZATION_TRADE_DETAILS",
		columns:     "ALL",
		sortColumns: "TRADE_DATE,SECURITY_CODE",
		sortTypes:   "-1,1",
		pageSize:    "5000",
		filter:      dragonTigerDateFilter(options),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DragonTigerInstitutionItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseDragonTigerInstitution(item))
	}
	return rows, nil
}

// GetDragonTigerBranchRank fetches brokerage branch ranking rows.
func GetDragonTigerBranchRank(ctx context.Context, client DragonTigerClient, endpoint string, period DragonTigerPeriod) ([]types.DragonTigerBranchItem, error) {
	cycle, err := dragonTigerCycle(period)
	if err != nil {
		return nil, err
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BILLBOARD_TRADEDETAILS",
		columns:     "ALL",
		sortColumns: "TOTAL_BUYER_SALESTIMES",
		sortTypes:   "-1",
		pageSize:    "5000",
		filter:      fmt.Sprintf(`(STATISTICS_CYCLE="%s")`, cycle),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DragonTigerBranchItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseDragonTigerBranch(item))
	}
	return rows, nil
}

// GetDragonTigerStockSeatDetail fetches buy and sell seat detail rows for a stock on a date.
func GetDragonTigerStockSeatDetail(ctx context.Context, client DragonTigerClient, endpoint string, symbol string, date string) ([]types.DragonTigerSeatItem, error) {
	code := stripAsharePrefix(symbol)
	queryDate := toISODate(date)
	filter := fmt.Sprintf(`(SECURITY_CODE="%s")(TRADE_DATE='%s')`, code, queryDate)
	buyItems, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BILLBOARD_DAILYDETAILSBUY",
		columns:     "ALL",
		sortColumns: "BUY_AMT_REAL",
		sortTypes:   "-1",
		pageSize:    "100",
		filter:      filter,
	})
	if err != nil {
		return nil, err
	}
	sellItems, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BILLBOARD_DAILYDETAILSSELL",
		columns:     "ALL",
		sortColumns: "SELL_AMT_REAL",
		sortTypes:   "-1",
		pageSize:    "100",
		filter:      filter,
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DragonTigerSeatItem, 0, len(buyItems)+len(sellItems))
	for index, item := range buyItems {
		rows = append(rows, parseDragonTigerSeat(item, "buy", index+1))
	}
	for index, item := range sellItems {
		rows = append(rows, parseDragonTigerSeat(item, "sell", index+1))
	}
	return rows, nil
}

func dragonTigerCycle(period DragonTigerPeriod) (string, error) {
	if period == "" {
		period = DragonTigerPeriodOneMonth
	}
	cycle, ok := dragonTigerPeriodMap[period]
	if !ok {
		return "", invalidArgumentError(fmt.Sprintf("invalid dragon tiger period %q", period))
	}
	return cycle, nil
}

func dragonTigerDateFilter(options DragonTigerDateOptions) string {
	return fmt.Sprintf(`(TRADE_DATE<='%s')(TRADE_DATE>='%s')`, toISODate(options.EndDate), toISODate(options.StartDate))
}

func toISODate(date string) string {
	if len(date) == 8 && date[4] != '-' {
		return date[:4] + "-" + date[4:6] + "-" + date[6:8]
	}
	return date
}

func stripAsharePrefix(symbol string) string {
	code := strings.TrimSpace(symbol)
	lower := strings.ToLower(code)
	for _, prefix := range []string{"sh", "sz", "bj"} {
		if strings.HasPrefix(lower, prefix) {
			return code[len(prefix):]
		}
	}
	return code
}

func parseDragonTigerDetail(item boardDynamicItem) types.DragonTigerDetailItem {
	reason := stringValue(item["EXPLANATION"])
	if reason == "" {
		reason = stringValue(item["EXPLAIN"])
	}
	return types.DragonTigerDetailItem{
		Code:             stringValue(item["SECURITY_CODE"]),
		Name:             stringValue(item["SECURITY_NAME_ABBR"]),
		Date:             parseDatacenterDate(item["TRADE_DATE"]),
		Close:            toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent:    toNumberFromAny(item["CHANGE_RATE"]),
		NetBuyAmount:     toNumberFromAny(item["BILLBOARD_NET_AMT"]),
		BuyAmount:        toNumberFromAny(item["BILLBOARD_BUY_AMT"]),
		SellAmount:       toNumberFromAny(item["BILLBOARD_SELL_AMT"]),
		DealAmount:       toNumberFromAny(item["BILLBOARD_DEAL_AMT"]),
		TotalAmount:      toNumberFromAny(item["ACCUM_AMOUNT"]),
		NetBuyRatio:      toNumberFromAny(item["DEAL_NET_RATIO"]),
		DealAmountRatio:  toNumberFromAny(item["DEAL_AMOUNT_RATIO"]),
		TurnoverRate:     toNumberFromAny(item["TURNOVERRATE"]),
		FloatMarketValue: toNumberFromAny(item["FREE_MARKET_CAP"]),
		Reason:           reason,
		AfterChange1D:    toNumberFromAny(item["D1_CLOSE_ADJCHRATE"]),
		AfterChange2D:    toNumberFromAny(item["D2_CLOSE_ADJCHRATE"]),
		AfterChange5D:    toNumberFromAny(item["D5_CLOSE_ADJCHRATE"]),
		AfterChange10D:   toNumberFromAny(item["D10_CLOSE_ADJCHRATE"]),
	}
}

func parseDragonTigerStockStat(item boardDynamicItem) types.DragonTigerStockStatItem {
	return types.DragonTigerStockStatItem{
		Code:            stringValue(item["SECURITY_CODE"]),
		Name:            stringValue(item["SECURITY_NAME_ABBR"]),
		LatestDate:      parseDatacenterDate(item["LATEST_TDATE"]),
		Close:           toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent:   toNumberFromAny(item["CHANGE_RATE"]),
		Count:           toNumberFromAny(item["BILLBOARD_TIMES"]),
		TotalBuyAmount:  toNumberFromAny(item["BILLBOARD_BUY_AMT"]),
		TotalSellAmount: toNumberFromAny(item["BILLBOARD_SELL_AMT"]),
		TotalNetAmount:  toNumberFromAny(item["BILLBOARD_NET_AMT"]),
		TotalDealAmount: toNumberFromAny(item["BILLBOARD_DEAL_AMT"]),
		BuyOrgCount:     toNumberFromAny(item["ORG_BUY_TIMES"]),
		SellOrgCount:    toNumberFromAny(item["ORG_SELL_TIMES"]),
	}
}

func parseDragonTigerInstitution(item boardDynamicItem) types.DragonTigerInstitutionItem {
	return types.DragonTigerInstitutionItem{
		Code:          stringValue(item["SECURITY_CODE"]),
		Name:          stringValue(item["SECURITY_NAME_ABBR"]),
		Date:          parseDatacenterDate(item["TRADE_DATE"]),
		Close:         toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent: toNumberFromAny(item["CHANGE_RATE"]),
		BuyOrgCount:   toNumberFromAny(item["BUY_TIMES"]),
		SellOrgCount:  toNumberFromAny(item["SELL_TIMES"]),
		OrgBuyAmount:  toNumberFromAny(item["BUY_AMT"]),
		OrgSellAmount: toNumberFromAny(item["SELL_AMT"]),
		OrgNetAmount:  toNumberFromAny(item["NET_AMT"]),
	}
}

func parseDragonTigerBranch(item boardDynamicItem) types.DragonTigerBranchItem {
	return types.DragonTigerBranchItem{
		Code:            stringValue(item["OPERATEDEPT_CODE"]),
		Name:            stringValue(item["OPERATEDEPT_NAME"]),
		TotalBuyAmount:  dragonTigerNullishFallbackNumber(item, "TOTAL_BUYAMT", "BUY_AMT"),
		TotalSellAmount: dragonTigerNullishFallbackNumber(item, "TOTAL_SELLAMT", "SELL_AMT"),
		BuyCount:        dragonTigerNullishFallbackNumber(item, "TOTAL_BUYER_SALESTIMES", "BUY_TIMES"),
		SellCount:       dragonTigerNullishFallbackNumber(item, "TOTAL_SELLER_SALESTIMES", "SELL_TIMES"),
		TotalCount:      toNumberFromAny(item["TOTAL_TIMES"]),
	}
}

func parseDragonTigerSeat(item boardDynamicItem, side string, defaultRank int) types.DragonTigerSeatItem {
	rank := toNumberFromAny(item["RANK"])
	if rank == nil {
		value := float64(defaultRank)
		rank = &value
	}
	return types.DragonTigerSeatItem{
		Rank:            rank,
		BranchName:      stringValue(item["OPERATEDEPT_NAME"]),
		BuyAmount:       dragonTigerNullishFallbackNumber(item, "BUY_AMT_REAL", "BUY_AMT"),
		BuyAmountRatio:  dragonTigerNullishFallbackNumber(item, "BUY_RATIO_TOTAL", "BUY_AMT_RATIO"),
		SellAmount:      dragonTigerNullishFallbackNumber(item, "SELL_AMT_REAL", "SELL_AMT"),
		SellAmountRatio: dragonTigerNullishFallbackNumber(item, "SELL_RATIO_TOTAL", "SELL_AMT_RATIO"),
		NetAmount:       toNumberFromAny(item["NET_AMT"]),
		Side:            side,
	}
}

func dragonTigerNullishFallbackNumber(item boardDynamicItem, primaryKey, fallbackKey string) *float64 {
	value, ok := item[primaryKey]
	if !ok || value == nil {
		return toNumberFromAny(item[fallbackKey])
	}
	return toNumberFromAny(value)
}
