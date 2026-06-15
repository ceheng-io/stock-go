package eastmoney

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

// FundDividendRank identifies a fund dividend list sort column.
type FundDividendRank = types.FundDividendRank

const (
	FundDividendRankCode             FundDividendRank = types.FundDividendRankCode
	FundDividendRankName             FundDividendRank = types.FundDividendRankName
	FundDividendRankEquityRecordDate FundDividendRank = types.FundDividendRankEquityRecordDate
	FundDividendRankExDividendDate   FundDividendRank = types.FundDividendRankExDividendDate
	FundDividendRankDividendPerShare FundDividendRank = types.FundDividendRankDividendPerShare
	FundDividendRankPayDate          FundDividendRank = types.FundDividendRankPayDate
)

// FundSortDirection identifies a fund list sort direction.
type FundSortDirection = types.FundSortDirection

const (
	FundSortAsc  FundSortDirection = types.FundSortAsc
	FundSortDesc FundSortDirection = types.FundSortDesc
)

// FundDividendListOptions configures public fund dividend list requests.
type FundDividendListOptions = types.FundDividendListOptions

type normalizedFundDividendListOptions struct {
	Year     string
	Page     int
	AllPages bool
	FundType string
	Rank     FundDividendRank
	Sort     FundSortDirection
	Code     string
}

// GetFundDividendList fetches public fund dividend distribution rows.
func GetFundDividendList(ctx context.Context, client FundClient, endpoint string, options FundDividendListOptions) (types.FundDividendListResult, error) {
	normalized := normalizeFundDividendListOptions(options)
	if !normalized.AllPages {
		return fetchFundDividendPage(ctx, client, endpoint, normalized)
	}
	first, err := fetchFundDividendPage(ctx, client, endpoint, normalized)
	if err != nil {
		return types.FundDividendListResult{}, err
	}
	items := first.Items
	for page := 2; page <= first.TotalPages; page++ {
		pageOptions := normalized
		pageOptions.Page = page
		next, err := fetchFundDividendPage(ctx, client, endpoint, pageOptions)
		if err != nil {
			return types.FundDividendListResult{}, err
		}
		items = append(items, next.Items...)
	}
	return types.FundDividendListResult{
		Items:       filterFundDividendsByCode(items, normalized.Code),
		TotalPages:  first.TotalPages,
		PageSize:    first.PageSize,
		CurrentPage: -1,
	}, nil
}

func fetchFundDividendPage(ctx context.Context, client FundClient, endpoint string, options normalizedFundDividendListOptions) (types.FundDividendListResult, error) {
	text, err := client.GetText(ctx, fundDividendListURL(endpoint, options))
	if err != nil {
		return types.FundDividendListResult{}, err
	}
	pageInfo := []float64{0, 0, float64(options.Page)}
	_ = core.ExtractJSVar(text, "pageinfo", &pageInfo)
	var rows [][]string
	_ = core.ExtractJSVar(text, "jjfh_data", &rows)
	items := make([]types.FundDividend, 0, len(rows))
	for _, row := range rows {
		items = append(items, parseFundDividend(row))
	}
	return types.FundDividendListResult{
		Items:       filterFundDividendsByCode(items, options.Code),
		TotalPages:  fundPageInfoValue(pageInfo, 0),
		PageSize:    fundPageInfoValue(pageInfo, 1),
		CurrentPage: fundPageInfoValue(pageInfo, 2),
	}, nil
}

func normalizeFundDividendListOptions(options FundDividendListOptions) normalizedFundDividendListOptions {
	normalized := normalizedFundDividendListOptions{
		Year:     options.Year,
		FundType: options.FundType,
		Rank:     FundDividendRank(options.Rank),
		Sort:     FundSortDirection(options.Sort),
		Code:     strings.TrimSpace(options.Code),
		AllPages: options.AllPages,
	}
	if strings.TrimSpace(options.Year) == "" {
		normalized.Year = strconv.Itoa(time.Now().Year())
	}
	normalized.Page = 1
	switch page := options.Page.(type) {
	case string:
		if strings.EqualFold(strings.TrimSpace(page), "all") {
			normalized.AllPages = true
		} else if parsed, err := strconv.Atoi(strings.TrimSpace(page)); err == nil && parsed > 0 {
			normalized.Page = parsed
		}
	case int:
		if page > 0 {
			normalized.Page = page
		}
	case int64:
		if page > 0 {
			normalized.Page = int(page)
		}
	case float64:
		if page > 0 {
			normalized.Page = int(page)
		}
	}
	if normalized.Rank == "" {
		normalized.Rank = FundDividendRankExDividendDate
	}
	if normalized.Sort == "" {
		normalized.Sort = FundSortDesc
	}
	return normalized
}

func fundDividendListURL(endpoint string, options normalizedFundDividendListOptions) string {
	params := url.Values{}
	params.Set("dt", "8")
	params.Set("page", strconv.Itoa(options.Page))
	params.Set("rank", string(options.Rank))
	params.Set("sort", string(options.Sort))
	params.Set("gs", "")
	params.Set("ftype", options.FundType)
	params.Set("year", options.Year)
	return endpoint + "?" + params.Encode()
}

func parseFundDividend(row []string) types.FundDividend {
	return types.FundDividend{
		Code:             fundRowValue(row, 0),
		Name:             fundRowValue(row, 1),
		EquityRecordDate: nullableDate(parseDatacenterDate(fundRowValue(row, 2))),
		ExDividendDate:   nullableDate(parseDatacenterDate(fundRowValue(row, 3))),
		DividendPerShare: nullableNumber(fundRowValue(row, 4)),
		PayDate:          nullableDate(parseDatacenterDate(fundRowValue(row, 5))),
	}
}

func nullableDate(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func fundRowValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func fundPageInfoValue(pageInfo []float64, index int) int {
	if index < 0 || index >= len(pageInfo) {
		return 0
	}
	return int(pageInfo[index])
}

func filterFundDividendsByCode(items []types.FundDividend, code string) []types.FundDividend {
	if code == "" {
		return items
	}
	filtered := make([]types.FundDividend, 0, len(items))
	for _, item := range items {
		if item.Code == code {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
