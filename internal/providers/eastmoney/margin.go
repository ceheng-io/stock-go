package eastmoney

import (
	"context"
	"fmt"

	"github.com/ceheng-io/stock-go/types"
)

// MarginClient is the minimal client interface required by margin providers.
type MarginClient interface {
	GetJSON(context.Context, string, any) error
}

// GetMarginAccountInfo fetches daily margin account statistics rows.
func GetMarginAccountInfo(ctx context.Context, client MarginClient, endpoint string) ([]types.MarginAccountItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPTA_WEB_MARGIN_DAILYTRADE",
		columns:     "ALL",
		sortColumns: "STATISTICS_DATE",
		sortTypes:   "-1",
		pageSize:    "500",
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.MarginAccountItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseMarginAccount(item))
	}
	return rows, nil
}

// GetMarginTargetList fetches stock margin target detail rows.
func GetMarginTargetList(ctx context.Context, client MarginClient, endpoint string, date string) ([]types.MarginTargetItem, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_MARGIN_TRADE_DETAIL",
		columns:     "ALL",
		sortColumns: "FIN_BALANCE",
		sortTypes:   "-1",
		pageSize:    "5000",
		filter:      marginDateFilter(date),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.MarginTargetItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseMarginTarget(item))
	}
	return rows, nil
}

func marginDateFilter(date string) string {
	if date == "" {
		return ""
	}
	return fmt.Sprintf(`(TRADE_DATE='%s')`, toISODate(date))
}

func parseMarginAccount(item boardDynamicItem) types.MarginAccountItem {
	return types.MarginAccountItem{
		Date:                   parseDatacenterDate(firstValue(item, "STATISTICS_DATE", "TRADE_DATE")),
		FinBalance:             toNumberFromAny(item["FIN_BALANCE"]),
		LoanBalance:            toNumberFromAny(item["LOAN_BALANCE"]),
		FinBuyAmount:           toNumberFromAny(item["FIN_BUY_AMT"]),
		LoanSellAmount:         toNumberFromAny(item["LOAN_SELL_AMT"]),
		InvestorCount:          numberFromKeys(item, "OPERATE_INVESTOR_NUM", "INVESTOR_NUM"),
		LiabilityInvestorCount: toNumberFromAny(item["MARGIN_INVESTOR_NUM"]),
		TotalGuarantee:         toNumberFromAny(item["TOTAL_GUARANTEE"]),
		AvgGuaranteeRatio:      toNumberFromAny(item["AVG_GUARANTEE_RATIO"]),
	}
}

func parseMarginTarget(item boardDynamicItem) types.MarginTargetItem {
	return types.MarginTargetItem{
		Code:            stringValue(item["SECURITY_CODE"]),
		Name:            stringValue(item["SECURITY_NAME_ABBR"]),
		Date:            parseDatacenterDate(item["TRADE_DATE"]),
		FinBalance:      toNumberFromAny(item["FIN_BALANCE"]),
		FinBuyAmount:    toNumberFromAny(item["FIN_BUY_AMT"]),
		FinRepayAmount:  toNumberFromAny(item["FIN_REPAY_AMT"]),
		LoanBalance:     toNumberFromAny(item["LOAN_BALANCE"]),
		LoanSellVolume:  toNumberFromAny(item["LOAN_SELL_VOLUME"]),
		LoanRepayVolume: toNumberFromAny(item["LOAN_REPAY_VOLUME"]),
	}
}

func firstValue(item boardDynamicItem, keys ...string) any {
	for _, key := range keys {
		if value, ok := item[key]; ok && value != nil {
			return value
		}
	}
	return nil
}
