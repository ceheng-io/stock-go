package tencent

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/ceheng.io/stock-go/types"
)

// SearchClient is the minimal client interface required by Tencent search.
type SearchClient interface {
	GetText(context.Context, string) (string, error)
	TencentSearchURL(string) string
}

// CalendarClient is the minimal client interface required by the trading calendar provider.
type CalendarClient interface {
	GetText(context.Context, string) (string, error)
	CalendarURL() string
}

// NormalizeSearchType maps Tencent raw asset types to normalized categories.
func NormalizeSearchType(rawType string) types.SearchResultType {
	upper := strings.ToUpper(rawType)
	if strings.HasPrefix(upper, "QDII") ||
		strings.HasPrefix(upper, "ETF") ||
		strings.HasPrefix(upper, "LOF") ||
		strings.HasPrefix(upper, "KJ") ||
		strings.HasPrefix(upper, "JJ") ||
		strings.Contains(upper, "FUND") {
		return types.SearchFund
	}
	if strings.HasPrefix(upper, "GP") || strings.Contains(upper, "STOCK") {
		return types.SearchStock
	}
	if upper == "ZS" || strings.Contains(upper, "INDEX") {
		return types.SearchIndex
	}
	if strings.HasPrefix(upper, "ZQ") || strings.Contains(upper, "BOND") {
		return types.SearchBond
	}
	if strings.HasPrefix(upper, "QH") || strings.Contains(upper, "FUTURE") {
		return types.SearchFutures
	}
	if strings.HasPrefix(upper, "QZ") || strings.Contains(upper, "OPTION") {
		return types.SearchOption
	}
	return types.SearchOther
}

// Search queries Tencent Smartbox.
func Search(ctx context.Context, client SearchClient, keyword string) ([]types.SearchResult, error) {
	if strings.TrimSpace(keyword) == "" {
		return []types.SearchResult{}, nil
	}
	text, err := client.GetText(ctx, client.TencentSearchURL(keyword))
	if err != nil {
		return nil, err
	}
	raw := extractHint(text)
	return parseSearchResult(raw), nil
}

var hintPattern = regexp.MustCompile(`v_hint="([^"]*)"`)

func extractHint(text string) string {
	match := hintPattern.FindStringSubmatch(text)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func parseSearchResult(raw string) []types.SearchResult {
	if raw == "" || raw == "N" {
		return []types.SearchResult{}
	}
	records := strings.Split(raw, "^")
	results := make([]types.SearchResult, 0, len(records))
	for _, record := range records {
		if record == "" {
			continue
		}
		fields := strings.Split(record, "~")
		market := field(fields, 0)
		pureCode := field(fields, 1)
		rawType := field(fields, 4)
		results = append(results, types.SearchResult{
			Code:     market + pureCode,
			Name:     decodeUnicode(field(fields, 2)),
			Market:   market,
			Type:     rawType,
			Category: NormalizeSearchType(rawType),
		})
	}
	return results
}

var unicodeEscapePattern = regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)

func decodeUnicode(value string) string {
	return unicodeEscapePattern.ReplaceAllStringFunc(value, func(match string) string {
		hex := strings.TrimPrefix(match, `\u`)
		code, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return match
		}
		return string(rune(code))
	})
}

// GetTradingCalendar fetches and parses A-share trading dates.
func GetTradingCalendar(ctx context.Context, client CalendarClient) ([]string, error) {
	text, err := client.GetText(ctx, client.CalendarURL())
	if err != nil {
		return nil, err
	}
	return ParseTradingCalendar(text), nil
}

// ParseTradingCalendar parses comma-separated trading date text.
func ParseTradingCalendar(text string) []string {
	if strings.TrimSpace(text) == "" {
		return []string{}
	}
	parts := strings.Split(strings.TrimSpace(text), ",")
	dates := make([]string, 0, len(parts))
	for _, part := range parts {
		date := strings.TrimSpace(part)
		if date != "" {
			dates = append(dates, date)
		}
	}
	return dates
}
