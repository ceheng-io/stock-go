package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type datacenterClient interface {
	GetJSON(context.Context, string, any) error
}

type datacenterResponse struct {
	Result *datacenterResultPayload `json:"result"`
}

type datacenterResultPayload struct {
	Pages int             `json:"pages"`
	Count int             `json:"count"`
	Data  json.RawMessage `json:"data"`
}

// DatacenterQuery describes Eastmoney datacenter-web query parameters.
type DatacenterQuery struct {
	ReportName    string
	Columns       string
	Filter        string
	SortColumns   string
	SortTypes     string
	PageSize      int
	StartPage     int
	FetchAllPages *bool
	MaxPages      int
	QuoteColumns  string
	QuoteType     string
	ExtraParams   map[string]string
}

// DatacenterResult contains merged datacenter-web page data and first-page metadata.
type DatacenterResult[T any] struct {
	Data  []T
	Total int
	Pages int
}

type datacenterOptions struct {
	reportName    string
	columns       string
	quoteColumns  string
	quoteType     string
	sortColumns   string
	sortTypes     string
	pageSize      string
	filter        string
	fetchAllPages *bool
}

func fetchDatacenter(ctx context.Context, client datacenterClient, endpoint string, options datacenterOptions) ([]boardDynamicItem, error) {
	pageSize := 0
	if options.pageSize != "" {
		if parsed, err := strconv.Atoi(options.pageSize); err == nil {
			pageSize = parsed
		}
	}
	result, err := FetchDatacenter(ctx, client, endpoint, DatacenterQuery{
		ReportName:    options.reportName,
		Columns:       options.columns,
		QuoteColumns:  options.quoteColumns,
		QuoteType:     options.quoteType,
		SortColumns:   options.sortColumns,
		SortTypes:     options.sortTypes,
		PageSize:      pageSize,
		Filter:        options.filter,
		FetchAllPages: options.fetchAllPages,
	}, func(item map[string]any, _ int) boardDynamicItem {
		return boardDynamicItem(item)
	})
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// FetchDatacenter fetches and merges Eastmoney datacenter-web pages.
func FetchDatacenter[T any](
	ctx context.Context,
	client datacenterClient,
	endpoint string,
	query DatacenterQuery,
	mapper func(item map[string]any, index int) T,
) (DatacenterResult[T], error) {
	query = normalizeDatacenterQuery(query)
	allData := []T{}
	page := query.StartPage
	totalPages := 1
	totalCount := 0
	pagesFetched := 0

	for {
		payload, err := fetchDatacenterPage(ctx, client, endpoint, query, page)
		if err != nil {
			return DatacenterResult[T]{}, err
		}
		if payload.Result == nil {
			break
		}
		data, err := decodeBoardDynamicArray(payload.Result.Data)
		if err != nil {
			return DatacenterResult[T]{}, err
		}
		if len(data) == 0 && !isJSONArrayPayload(payload.Result.Data) {
			break
		}
		if page == query.StartPage {
			totalPages = payload.Result.Pages
			if totalPages <= 0 {
				totalPages = 1
			}
			totalCount = payload.Result.Count
			if totalCount == 0 {
				totalCount = len(data)
			}
		}

		for _, item := range data {
			allData = append(allData, mapper(map[string]any(item), len(allData)))
		}
		pagesFetched++
		if !datacenterFetchAllPages(query) || page >= totalPages || pagesFetched >= query.MaxPages {
			break
		}
		page++
	}

	return DatacenterResult[T]{
		Data:  allData,
		Total: totalCount,
		Pages: totalPages,
	}, nil
}

// FetchDatacenterList fetches datacenter-web pages and returns only mapped data.
func FetchDatacenterList[T any](
	ctx context.Context,
	client datacenterClient,
	endpoint string,
	query DatacenterQuery,
	mapper func(item map[string]any, index int) T,
) ([]T, error) {
	result, err := FetchDatacenter(ctx, client, endpoint, query, mapper)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func fetchDatacenterPage(ctx context.Context, client datacenterClient, endpoint string, options DatacenterQuery, page int) (datacenterResponse, error) {
	params := url.Values{}
	params.Set("reportName", options.ReportName)
	params.Set("columns", options.Columns)
	params.Set("pageSize", strconv.Itoa(options.PageSize))
	params.Set("pageNumber", strconv.Itoa(page))
	if options.Filter != "" {
		params.Set("filter", options.Filter)
	}
	if options.SortColumns != "" {
		params.Set("sortColumns", options.SortColumns)
	}
	if options.SortTypes != "" {
		params.Set("sortTypes", options.SortTypes)
	}
	if options.QuoteColumns != "" {
		params.Set("quoteColumns", options.QuoteColumns)
	}
	if options.QuoteType != "" {
		params.Set("quoteType", options.QuoteType)
	}
	for key, value := range options.ExtraParams {
		params.Set(key, value)
	}
	params.Set("source", "WEB")
	params.Set("client", "WEB")

	var payload datacenterResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return datacenterResponse{}, err
	}
	return payload, nil
}

func normalizeDatacenterQuery(query DatacenterQuery) DatacenterQuery {
	if query.Columns == "" {
		query.Columns = "ALL"
	}
	if query.PageSize <= 0 {
		query.PageSize = 500
	}
	if query.StartPage <= 0 {
		query.StartPage = 1
	}
	if query.MaxPages <= 0 {
		query.MaxPages = 1000
	}
	return query
}

func datacenterFetchAllPages(query DatacenterQuery) bool {
	if query.FetchAllPages == nil {
		return true
	}
	return *query.FetchAllPages
}

func datacenterBool(value bool) *bool {
	return &value
}

func parseDatacenterDate(value any) string {
	return ParseDCDate(value)
}

// ParseDCDate extracts common datacenter date values as YYYY-MM-DD.
func ParseDCDate(value any) string {
	text := stringValue(value)
	if text == "" {
		return ""
	}
	if len(text) >= 10 && text[4] == '-' {
		return text[:10]
	}
	if len(text) >= 8 {
		if parsed, err := time.Parse("20060102", text[:8]); err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return strings.TrimSpace(text)
}
