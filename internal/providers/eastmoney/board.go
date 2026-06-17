package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/ceheng-io/stock-go/types"
)

const emBoardToken = "bd1d9ddb04089700cf9c27f6f7426281"

const boardPageSize = 100

// BoardClient is the minimal client interface required by Eastmoney board providers.
type BoardClient interface {
	GetJSON(context.Context, string, any) error
}

type boardListResponse struct {
	Data struct {
		Total int             `json:"total"`
		Diff  json.RawMessage `json:"diff"`
	} `json:"data"`
}

type boardSpotResponse struct {
	Data map[string]float64 `json:"data"`
}

type boardKlineResponse struct {
	Data struct {
		Klines []string `json:"klines"`
	} `json:"data"`
}

type boardTimelineResponse struct {
	Data struct {
		Trends json.RawMessage `json:"trends"`
	} `json:"data"`
}

type boardMinuteKlineResponse struct {
	Data struct {
		Klines json.RawMessage `json:"klines"`
	} `json:"data"`
}

type boardDynamicItem map[string]any

// GetIndustryList fetches Eastmoney industry board list rows.
func GetIndustryList(ctx context.Context, client BoardClient, endpoint string) ([]types.Board, error) {
	return getBoardList(ctx, client, endpoint, "industry")
}

// GetConceptList fetches Eastmoney concept board list rows.
func GetConceptList(ctx context.Context, client BoardClient, endpoint string) ([]types.Board, error) {
	return getBoardList(ctx, client, endpoint, "concept")
}

func getBoardList(ctx context.Context, client BoardClient, endpoint string, boardType string) ([]types.Board, error) {
	params := url.Values{}
	params.Set("po", "1")
	params.Set("np", "1")
	params.Set("ut", emBoardToken)
	params.Set("fltt", "2")
	params.Set("invt", "2")
	params.Set("fid", boardSortField(boardType))
	params.Set("fs", boardFSFilter(boardType))
	params.Set("fields", boardListFields(boardType))

	items, err := fetchBoardPages(ctx, client, endpoint, params, boardPageSize)
	if err != nil {
		return nil, err
	}
	rows := make([]types.Board, 0, len(items))
	for index, item := range items {
		rows = append(rows, parseBoard(item, index+1))
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return numberValue(rows[i].ChangePercent) > numberValue(rows[j].ChangePercent)
	})
	for index := range rows {
		rows[index].Rank = index + 1
	}
	return rows, nil
}

func boardFSFilter(boardType string) string {
	if boardType == "concept" {
		return "m:90 t:3 f:!50"
	}
	return "m:90 t:2 f:!50"
}

func boardSortField(boardType string) string {
	if boardType == "concept" {
		return "f12"
	}
	return "f3"
}

func boardListFields(boardType string) string {
	if boardType == "concept" {
		return "f2,f3,f4,f8,f12,f14,f15,f16,f17,f18,f20,f21,f24,f25,f22,f33,f11,f62,f128,f124,f107,f104,f105,f136"
	}
	return "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152,f124,f107,f104,f105,f140,f141,f207,f208,f209,f222"
}

func parseBoard(item boardDynamicItem, rank int) types.Board {
	return types.Board{
		Rank:                      rank,
		Name:                      stringValue(item["f14"]),
		Code:                      stringValue(item["f12"]),
		Price:                     toNumberFromAny(item["f2"]),
		Change:                    toNumberFromAny(item["f4"]),
		ChangePercent:             toNumberFromAny(item["f3"]),
		TotalMarketCap:            toNumberFromAny(item["f20"]),
		TurnoverRate:              toNumberFromAny(item["f8"]),
		RiseCount:                 toNumberFromAny(item["f104"]),
		FallCount:                 toNumberFromAny(item["f105"]),
		LeadingStock:              nullableDatacenterString(item["f128"]),
		LeadingStockChangePercent: toNumberFromAny(item["f136"]),
	}
}

// GetBoardSpot fetches Eastmoney board spot metrics.
func GetBoardSpot(ctx context.Context, client BoardClient, boardCode string, endpoint string) ([]types.BoardSpot, error) {
	params := url.Values{}
	params.Set("fields", "f43,f44,f45,f46,f47,f48,f170,f171,f168,f169")
	params.Set("mpi", "1000")
	params.Set("invt", "2")
	params.Set("fltt", "1")
	params.Set("secid", "90."+normalizeBoardCode(boardCode))

	var payload boardSpotResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return []types.BoardSpot{}, nil
	}
	fieldMap := []struct {
		key   string
		name  string
		scale bool
	}{
		{key: "f43", name: "最新", scale: true},
		{key: "f44", name: "最高", scale: true},
		{key: "f45", name: "最低", scale: true},
		{key: "f46", name: "开盘", scale: true},
		{key: "f47", name: "成交量"},
		{key: "f48", name: "成交额"},
		{key: "f170", name: "涨跌幅", scale: true},
		{key: "f171", name: "振幅", scale: true},
		{key: "f168", name: "换手率", scale: true},
		{key: "f169", name: "涨跌额", scale: true},
	}
	rows := make([]types.BoardSpot, 0, len(fieldMap))
	for _, field := range fieldMap {
		value, ok := payload.Data[field.key]
		var number *float64
		if ok {
			if field.scale {
				value = value / 100
			}
			number = &value
		}
		rows = append(rows, types.BoardSpot{Item: field.name, Value: number})
	}
	return rows, nil
}

// GetBoardConstituents fetches Eastmoney board constituent stocks.
func GetBoardConstituents(ctx context.Context, client BoardClient, boardCode string, endpoint string) ([]types.BoardConstituent, error) {
	params := url.Values{}
	params.Set("po", "1")
	params.Set("np", "1")
	params.Set("ut", emBoardToken)
	params.Set("fltt", "2")
	params.Set("invt", "2")
	params.Set("fid", "f3")
	params.Set("fs", "b:"+normalizeBoardCode(boardCode)+" f:!50")
	params.Set("fields", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152,f45")

	items, err := fetchBoardPages(ctx, client, endpoint, params, boardPageSize)
	if err != nil {
		return nil, err
	}
	rows := make([]types.BoardConstituent, 0, len(items))
	for index, item := range items {
		rows = append(rows, types.BoardConstituent{
			Rank:          index + 1,
			Code:          stringValue(item["f12"]),
			Name:          stringValue(item["f14"]),
			Price:         toNumberFromAny(item["f2"]),
			ChangePercent: toNumberFromAny(item["f3"]),
			Change:        toNumberFromAny(item["f4"]),
			Volume:        toNumberFromAny(item["f5"]),
			Amount:        toNumberFromAny(item["f6"]),
			Amplitude:     toNumberFromAny(item["f7"]),
			High:          toNumberFromAny(item["f15"]),
			Low:           toNumberFromAny(item["f16"]),
			Open:          toNumberFromAny(item["f17"]),
			PrevClose:     toNumberFromAny(item["f18"]),
			TurnoverRate:  toNumberFromAny(item["f8"]),
			PE:            toNumberFromAny(item["f9"]),
			PB:            toNumberFromAny(item["f23"]),
		})
	}
	return rows, nil
}

func fetchBoardPages(ctx context.Context, client BoardClient, endpoint string, baseParams url.Values, pageSize int) ([]boardDynamicItem, error) {
	if pageSize <= 0 {
		pageSize = boardPageSize
	}
	allItems := []boardDynamicItem{}
	total := 0
	for page := 1; ; page++ {
		params := cloneValues(baseParams)
		params.Set("pn", strconv.Itoa(page))
		params.Set("pz", strconv.Itoa(pageSize))

		var payload boardListResponse
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
		if total <= 0 || len(allItems) >= total || len(diff) == 0 {
			break
		}
	}
	return allItems, nil
}

// GetBoardKline fetches Eastmoney board historical K-line rows.
func GetBoardKline(ctx context.Context, client BoardClient, boardCode string, endpoint string, options HistoryKlineOptions) ([]types.BoardKline, error) {
	options = normalizeBoardKlineOptions(options)
	if err := validateHistoryKlineOptions(options); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("secid", "90."+normalizeBoardCode(boardCode))
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	params.Set("klt", periodCode(options.Period))
	params.Set("fqt", adjustCode(options.Adjust))
	params.Set("beg", options.StartDate)
	params.Set("end", options.EndDate)
	params.Set("smplmt", "10000")
	params.Set("lmt", "1000000")

	var payload boardKlineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	rows := make([]types.BoardKline, 0, len(payload.Data.Klines))
	for _, line := range payload.Data.Klines {
		rows = append(rows, parseBoardKlineCSV(line))
	}
	return rows, nil
}

func normalizeBoardKlineOptions(options HistoryKlineOptions) HistoryKlineOptions {
	if options.Period == "" {
		options.Period = KlinePeriodDaily
	}
	if options.Adjust == "" {
		options.Adjust = AdjustNone
	}
	if options.StartDate == "" {
		options.StartDate = "19700101"
	}
	if options.EndDate == "" {
		options.EndDate = "20500101"
	}
	return options
}

// GetBoardMinuteKline fetches Eastmoney board minute K-line or timeline rows.
func GetBoardMinuteKline(ctx context.Context, client BoardClient, boardCode string, klineEndpoint string, trendsEndpoint string, options MinuteKlineOptions) (types.BoardMinuteKlineResult, error) {
	options = normalizeBoardMinuteOptions(options)
	if err := validateMinuteKlineOptions(options); err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	if options.Period == MinutePeriodOne {
		return getBoardMinuteTimeline(ctx, client, boardCode, trendsEndpoint)
	}
	return getBoardMinuteRows(ctx, client, boardCode, klineEndpoint, options)
}

func normalizeBoardMinuteOptions(options MinuteKlineOptions) MinuteKlineOptions {
	if options.Period == "" {
		options.Period = MinutePeriodFive
	}
	if options.Adjust == "" {
		options.Adjust = AdjustQFQ
	}
	return options
}

func getBoardMinuteTimeline(ctx context.Context, client BoardClient, boardCode string, endpoint string) (types.BoardMinuteKlineResult, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	params.Set("iscr", "0")
	params.Set("ndays", "1")
	params.Set("secid", "90."+normalizeBoardCode(boardCode))

	var payload boardTimelineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	trends, err := decodeStringArray(payload.Data.Trends)
	if err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	rows := make([]types.BoardMinuteTimeline, 0, len(trends))
	for _, line := range trends {
		rows = append(rows, parseBoardMinuteTimelineCSV(line))
	}
	return types.BoardMinuteKlineResult{Timeline: rows}, nil
}

func getBoardMinuteRows(ctx context.Context, client BoardClient, boardCode string, endpoint string, options MinuteKlineOptions) (types.BoardMinuteKlineResult, error) {
	params := url.Values{}
	params.Set("secid", "90."+normalizeBoardCode(boardCode))
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	params.Set("klt", string(options.Period))
	params.Set("fqt", "1")
	params.Set("beg", "0")
	params.Set("end", "20500101")
	params.Set("smplmt", "10000")
	params.Set("lmt", "1000000")

	var payload boardMinuteKlineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	klines, err := decodeStringArray(payload.Data.Klines)
	if err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	rows := make([]types.BoardMinuteKline, 0, len(klines))
	for _, line := range klines {
		rows = append(rows, parseBoardMinuteKlineCSV(line))
	}
	return types.BoardMinuteKlineResult{Klines: rows}, nil
}

func normalizeBoardCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case nil:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}

func toNumberFromAny(value any) *float64 {
	switch typed := value.(type) {
	case nil:
		return nil
	case float64:
		return &typed
	case float32:
		number := float64(typed)
		return &number
	case int:
		number := float64(typed)
		return &number
	case int64:
		number := float64(typed)
		return &number
	case string:
		if typed == "" || typed == "-" {
			return nil
		}
		number, err := strconv.ParseFloat(typed, 64)
		if err != nil {
			return nil
		}
		return &number
	default:
		return nil
	}
}

func numberValue(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func parseBoardKlineCSV(line string) types.BoardKline {
	fields := strings.Split(line, ",")
	return types.BoardKline{
		Date:          boardField(fields, 0),
		Open:          toNumber(boardField(fields, 1)),
		Close:         toNumber(boardField(fields, 2)),
		High:          toNumber(boardField(fields, 3)),
		Low:           toNumber(boardField(fields, 4)),
		Volume:        toNumber(boardField(fields, 5)),
		Amount:        toNumber(boardField(fields, 6)),
		Amplitude:     toNumber(boardField(fields, 7)),
		ChangePercent: toNumber(boardField(fields, 8)),
		Change:        toNumber(boardField(fields, 9)),
		TurnoverRate:  toNumber(boardField(fields, 10)),
	}
}

func parseBoardMinuteTimelineCSV(line string) types.BoardMinuteTimeline {
	fields := strings.Split(line, ",")
	return types.BoardMinuteTimeline{
		Time:   boardField(fields, 0),
		Open:   toNumber(boardField(fields, 1)),
		Close:  toNumber(boardField(fields, 2)),
		High:   toNumber(boardField(fields, 3)),
		Low:    toNumber(boardField(fields, 4)),
		Volume: toNumber(boardField(fields, 5)),
		Amount: toNumber(boardField(fields, 6)),
		Price:  toNumber(boardField(fields, 7)),
	}
}

func parseBoardMinuteKlineCSV(line string) types.BoardMinuteKline {
	row := parseBoardKlineCSV(line)
	return types.BoardMinuteKline{
		Time:          row.Date,
		Open:          row.Open,
		Close:         row.Close,
		High:          row.High,
		Low:           row.Low,
		Volume:        row.Volume,
		Amount:        row.Amount,
		Amplitude:     row.Amplitude,
		ChangePercent: row.ChangePercent,
		Change:        row.Change,
		TurnoverRate:  row.TurnoverRate,
	}
}

func boardField(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return fields[index]
}
