package eastmoney

import (
	"context"
	"fmt"

	"github.com/ceheng-io/stock-go/types"
)

type BlockTradeDateOptions = types.BlockTradeDateOptions

// BlockTradeClient is the minimal client interface required by block-trade providers.
type BlockTradeClient interface {
	GetJSON(context.Context, string, any) error
}

// GetBlockTradeMarketStat fetches block-trade market statistics rows.
func GetBlockTradeMarketStat(ctx context.Context, client BlockTradeClient, endpoint string) ([]types.BlockTradeMarketStatItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "PRT_BLOCKTRADE_MARKET_STA",
		columns:     "ALL",
		sortColumns: "TRADE_DATE",
		sortTypes:   "-1",
		pageSize:    "500",
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.BlockTradeMarketStatItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseBlockTradeMarketStat(item))
	}
	return rows, nil
}

// GetBlockTradeDetail fetches block-trade detail rows.
func GetBlockTradeDetail(ctx context.Context, client BlockTradeClient, endpoint string, options BlockTradeDateOptions) ([]types.BlockTradeDetailItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BLOCK_TRADE_DETAIL",
		columns:     "ALL",
		sortColumns: "TRADE_DATE,SECURITY_CODE",
		sortTypes:   "-1,1",
		pageSize:    "5000",
		filter:      blockTradeDateFilter(options),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.BlockTradeDetailItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseBlockTradeDetail(item))
	}
	return rows, nil
}

// GetBlockTradeDailyStat fetches block-trade daily stock statistics rows.
func GetBlockTradeDailyStat(ctx context.Context, client BlockTradeClient, endpoint string, options BlockTradeDateOptions) ([]types.BlockTradeDailyStatItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_BLOCK_TRADE_STA",
		columns:     "ALL",
		sortColumns: "TRADE_DATE,DEAL_AMT",
		sortTypes:   "-1,-1",
		pageSize:    "5000",
		filter:      blockTradeDateFilter(options),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.BlockTradeDailyStatItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseBlockTradeDailyStat(item))
	}
	return rows, nil
}

func blockTradeDateFilter(options BlockTradeDateOptions) string {
	filter := ""
	if options.StartDate != "" {
		filter += fmt.Sprintf(`(TRADE_DATE>='%s')`, toISODate(options.StartDate))
	}
	if options.EndDate != "" {
		filter += fmt.Sprintf(`(TRADE_DATE<='%s')`, toISODate(options.EndDate))
	}
	return filter
}

func parseBlockTradeMarketStat(item boardDynamicItem) types.BlockTradeMarketStatItem {
	return types.BlockTradeMarketStatItem{
		Date:            parseDatacenterDate(item["TRADE_DATE"]),
		SHClose:         numberFromKeys(item, "CLOSE_PRICE", "SH_CLOSE_PRICE"),
		SHChangePercent: numberFromKeys(item, "CHANGE_RATE", "SH_CHANGE_RATE"),
		TotalAmount:     numberFromKeys(item, "TURNOVER", "TOTAL_AMOUNT"),
		PremiumAmount:   numberFromKeys(item, "PREMIUM_TURNOVER", "PREMIUM_AMOUNT"),
		PremiumRatio:    toNumberFromAny(item["PREMIUM_RATIO"]),
		DiscountAmount:  numberFromKeys(item, "DISCOUNT_TURNOVER", "DISCOUNT_AMOUNT"),
		DiscountRatio:   toNumberFromAny(item["DISCOUNT_RATIO"]),
	}
}

func parseBlockTradeDetail(item boardDynamicItem) types.BlockTradeDetailItem {
	return types.BlockTradeDetailItem{
		Code:          stringValue(item["SECURITY_CODE"]),
		Name:          stringValue(item["SECURITY_NAME_ABBR"]),
		Date:          parseDatacenterDate(item["TRADE_DATE"]),
		Close:         toNumberFromAny(item["CLOSE_PRICE"]),
		ChangePercent: toNumberFromAny(item["CHANGE_RATE"]),
		DealPrice:     numberFromKeys(item, "DEAL_PRICE", "PRICE"),
		DealVolume:    numberFromKeys(item, "DEAL_VOLUME", "VOLUME"),
		DealAmount:    numberFromKeys(item, "DEAL_AMT", "TURNOVER"),
		PremiumRate:   numberFromKeys(item, "PREMIUM_RATIO", "PREMIUM_RATE"),
		BuyBranch:     stringFromKeys(item, "BUYER_DEPT", "BUYER_OPERATEDEPT_NAME"),
		SellBranch:    stringFromKeys(item, "SELLER_DEPT", "SELLER_OPERATEDEPT_NAME"),
	}
}

func parseBlockTradeDailyStat(item boardDynamicItem) types.BlockTradeDailyStatItem {
	return types.BlockTradeDailyStatItem{
		Code:            stringValue(item["SECURITY_CODE"]),
		Name:            stringValue(item["SECURITY_NAME_ABBR"]),
		Date:            parseDatacenterDate(item["TRADE_DATE"]),
		ChangePercent:   toNumberFromAny(item["CHANGE_RATE"]),
		Close:           toNumberFromAny(item["CLOSE_PRICE"]),
		DealCount:       numberFromKeys(item, "DEAL_NUM", "DEAL_COUNT"),
		DealTotalAmount: numberFromKeys(item, "DEAL_AMT", "TOTAL_AMOUNT"),
		DealTotalVolume: numberFromKeys(item, "DEAL_VOLUME", "TOTAL_VOLUME"),
		PremiumAmount:   numberFromKeys(item, "PREMIUM_AMT", "PREMIUM_AMOUNT"),
		DiscountAmount:  numberFromKeys(item, "DISCOUNT_AMT", "DISCOUNT_AMOUNT"),
	}
}

func numberFromKeys(item boardDynamicItem, keys ...string) *float64 {
	for _, key := range keys {
		value, ok := item[key]
		if !ok || value == nil {
			continue
		}
		return toNumberFromAny(value)
	}
	return nil
}

func stringFromKeys(item boardDynamicItem, keys ...string) string {
	for _, key := range keys {
		value, ok := item[key]
		if !ok || value == nil {
			continue
		}
		return stringValue(value)
	}
	return ""
}
