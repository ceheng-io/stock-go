package eastmoney

import (
	"context"
	"fmt"

	"github.com/ceheng-io/stock-go/types"
)

// DividendClient is the minimal client interface required by dividend providers.
type DividendClient interface {
	GetJSON(context.Context, string, any) error
}

// GetDividendDetail fetches stock dividend detail rows.
func GetDividendDetail(ctx context.Context, client DividendClient, endpoint string, symbol string) ([]types.DividendDetail, error) {
	items, err := fetchDatacenter(ctx, client, endpoint, datacenterOptions{
		reportName:  "RPT_SHAREBONUS_DET",
		columns:     "ALL",
		sortColumns: "REPORT_DATE",
		sortTypes:   "-1",
		pageSize:    "500",
		filter:      fmt.Sprintf(`(SECURITY_CODE="%s")`, stripAsharePrefix(symbol)),
	})
	if err != nil {
		return nil, err
	}
	rows := make([]types.DividendDetail, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseDividendDetail(item))
	}
	return rows, nil
}

func parseDividendDetail(item boardDynamicItem) types.DividendDetail {
	return types.DividendDetail{
		Code:                stringValue(item["SECURITY_CODE"]),
		Name:                stringValue(item["SECURITY_NAME_ABBR"]),
		ReportDate:          nullableDatacenterDate(item["REPORT_DATE"]),
		PlanNoticeDate:      nullableDatacenterDate(item["PLAN_NOTICE_DATE"]),
		DisclosureDate:      nullableDatacenterDate(firstValue(item, "PUBLISH_DATE", "PLAN_NOTICE_DATE")),
		AssignTransferRatio: toNumberFromAny(item["BONUS_IT_RATIO"]),
		BonusRatio:          toNumberFromAny(item["BONUS_RATIO"]),
		TransferRatio:       toNumberFromAny(item["IT_RATIO"]),
		DividendPretax:      toNumberFromAny(item["PRETAX_BONUS_RMB"]),
		DividendDesc:        nullableDatacenterText(item["IMPL_PLAN_PROFILE"]),
		DividendYield:       toNumberFromAny(item["DIVIDENT_RATIO"]),
		EPS:                 toNumberFromAny(item["BASIC_EPS"]),
		BPS:                 toNumberFromAny(item["BVPS"]),
		CapitalReserve:      toNumberFromAny(item["PER_CAPITAL_RESERVE"]),
		UnassignedProfit:    toNumberFromAny(item["PER_UNASSIGN_PROFIT"]),
		NetProfitYoY:        toNumberFromAny(item["PNP_YOY_RATIO"]),
		TotalShares:         toNumberFromAny(item["TOTAL_SHARES"]),
		EquityRecordDate:    nullableDatacenterDate(item["EQUITY_RECORD_DATE"]),
		ExDividendDate:      nullableDatacenterDate(item["EX_DIVIDEND_DATE"]),
		PayDate:             nullableDatacenterDate(item["PAY_DATE"]),
		AssignProgress:      nullableDatacenterText(item["ASSIGN_PROGRESS"]),
		NoticeDate:          nullableDatacenterDate(item["NOTICE_DATE"]),
	}
}

func nullableDatacenterDate(value any) *string {
	return nullableDatacenterString(parseDatacenterDate(value))
}

func nullableDatacenterString(value any) *string {
	text := stringValue(value)
	if text == "" {
		return nil
	}
	return &text
}

func nullableDatacenterText(value any) *string {
	if value == nil {
		return nil
	}
	text := stringValue(value)
	return &text
}
