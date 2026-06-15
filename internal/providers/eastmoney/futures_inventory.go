package eastmoney

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ceheng.io/stock-go/types"
)

// FuturesInventoryOptions 配置国内期货库存查询。
type FuturesInventoryOptions struct {
	StartDate string
	PageSize  int
}

// ComexInventoryOptions 配置 COMEX 库存查询。
type ComexInventoryOptions struct {
	PageSize int
}

// GetFuturesInventorySymbols 获取国内期货库存品种列表。
func GetFuturesInventorySymbols(ctx context.Context, client datacenterClient, endpoint string) ([]types.FuturesInventorySymbol, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:    "RPT_FUTU_POSITIONCODE",
		columns:       "TRADE_MARKET_CODE,TRADE_CODE,TRADE_TYPE",
		filter:        `(IS_MAINCODE="1")`,
		pageSize:      "500",
		fetchAllPages: datacenterBool(false),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.FuturesInventorySymbol, 0, len(items))
	for _, item := range items {
		rows = append(rows, types.FuturesInventorySymbol{
			Code:       stringValue(item["TRADE_CODE"]),
			Name:       stringValue(item["TRADE_TYPE"]),
			MarketCode: stringValue(item["TRADE_MARKET_CODE"]),
		})
	}
	return rows, nil
}

// GetFuturesInventory 获取国内期货库存数据。
func GetFuturesInventory(ctx context.Context, client datacenterClient, endpoint string, symbol string, options FuturesInventoryOptions) ([]types.FuturesInventory, error) {
	code := strings.ToUpper(strings.TrimSpace(symbol))
	options = normalizeFuturesInventoryOptions(options)
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_FUTU_STOCKDATA",
		columns:     "SECURITY_CODE,TRADE_DATE,ON_WARRANT_NUM,ADDCHANGE",
		filter:      fmt.Sprintf(`(SECURITY_CODE="%s")(TRADE_DATE>='%s')`, code, options.StartDate),
		pageSize:    strconv.Itoa(options.PageSize),
		sortColumns: "TRADE_DATE",
		sortTypes:   "-1",
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.FuturesInventory, 0, len(items))
	for _, item := range items {
		codeValue, ok := item["SECURITY_CODE"]
		rowCode := stringValue(codeValue)
		if !ok || codeValue == nil {
			rowCode = code
		}
		rows = append(rows, types.FuturesInventory{
			Code:      rowCode,
			Date:      parseDatacenterDate(item["TRADE_DATE"]),
			Inventory: toNumberFromAny(item["ON_WARRANT_NUM"]),
			Change:    toNumberFromAny(item["ADDCHANGE"]),
		})
	}
	return rows, nil
}

// GetComexInventory 获取 COMEX 黄金或白银库存数据。
func GetComexInventory(ctx context.Context, client datacenterClient, endpoint string, symbol string, options ComexInventoryOptions) ([]types.ComexInventory, error) {
	normalized := strings.ToLower(strings.TrimSpace(symbol))
	indicatorID, ok := comexInventorySymbolMap[normalized]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf(`invalid COMEX symbol %q: must be "gold" or "silver"`, symbol))
	}
	pageSize := options.PageSize
	if pageSize <= 0 {
		pageSize = 500
	}
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_FUTUOPT_GOLDSIL",
		filter:      fmt.Sprintf(`(INDICATOR_ID1="%s")(@STORAGE_TON<>"NULL")`, indicatorID),
		pageSize:    strconv.Itoa(pageSize),
		sortColumns: "REPORT_DATE",
		sortTypes:   "-1",
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.ComexInventory, 0, len(items))
	for _, item := range items {
		rows = append(rows, types.ComexInventory{
			Date:         parseDatacenterDate(item["REPORT_DATE"]),
			Name:         comexInventoryNameMap[normalized],
			StorageTon:   toNumberFromAny(item["STORAGE_TON"]),
			StorageOunce: toNumberFromAny(item["STORAGE_OUNCE"]),
		})
	}
	return rows, nil
}

func normalizeFuturesInventoryOptions(options FuturesInventoryOptions) FuturesInventoryOptions {
	if options.StartDate == "" {
		options.StartDate = "2020-10-28"
	}
	if options.PageSize <= 0 {
		options.PageSize = 500
	}
	return options
}

var comexInventorySymbolMap = map[string]string{
	"gold":   "EMI00069026",
	"silver": "EMI00069027",
}

var comexInventoryNameMap = map[string]string{
	"gold":   "黄金",
	"silver": "白银",
}
