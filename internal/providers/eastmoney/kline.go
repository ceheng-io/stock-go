package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/symbols"
	"github.com/ceheng.io/stock-go/timeutil"
	"github.com/ceheng.io/stock-go/types"
)

const emPushToken = "7eea3edcaed734bea9cbfc24409ed989"

// KlinePeriod is an Eastmoney historical K-line period.
type KlinePeriod string

const (
	KlinePeriodDaily   KlinePeriod = "daily"
	KlinePeriodWeekly  KlinePeriod = "weekly"
	KlinePeriodMonthly KlinePeriod = "monthly"
)

// AdjustType is an Eastmoney adjustment mode.
type AdjustType string

const (
	AdjustNone AdjustType = "none"
	AdjustQFQ  AdjustType = "qfq"
	AdjustHFQ  AdjustType = "hfq"
)

// HistoryKlineOptions configures historical K-line fetching.
type HistoryKlineOptions struct {
	Period    KlinePeriod
	Adjust    AdjustType
	StartDate string
	EndDate   string
}

// MinutePeriod is an Eastmoney minute period.
type MinutePeriod string

const (
	MinutePeriodOne     MinutePeriod = "1"
	MinutePeriodFive    MinutePeriod = "5"
	MinutePeriodFifteen MinutePeriod = "15"
	MinutePeriodThirty  MinutePeriod = "30"
	MinutePeriodSixty   MinutePeriod = "60"
)

// MinuteKlineOptions configures minute K-line fetching.
type MinuteKlineOptions struct {
	Period    MinutePeriod
	Adjust    AdjustType
	StartDate string
	EndDate   string
	NDays     int
}

// KlineClient is the minimal client interface required by Eastmoney K-line providers.
type KlineClient interface {
	GetJSON(context.Context, string, any) error
}

type historyKlineResponse struct {
	Data struct {
		Klines []string `json:"klines"`
		Code   string   `json:"code"`
		Name   string   `json:"name"`
	} `json:"data"`
}

type minuteTimelineResponse struct {
	Data struct {
		Trends json.RawMessage `json:"trends"`
	} `json:"data"`
}

type minuteKlineResponse struct {
	Data struct {
		Klines json.RawMessage `json:"klines"`
	} `json:"data"`
}

// GetHistoryKline fetches CN historical K-line rows.
func GetHistoryKline(ctx context.Context, client KlineClient, symbol string, endpoint string, options HistoryKlineOptions) ([]types.HistoryKline, error) {
	options = normalizeHistoryKlineOptions(options)
	if err := validateHistoryKlineOptions(options); err != nil {
		return nil, err
	}
	normalized, err := symbols.Normalize(symbol, &symbols.Hint{Market: symbols.MarketCN})
	if err != nil {
		return nil, err
	}
	secid, err := symbols.ToEastmoneySecIDE(normalized)
	if err != nil {
		return nil, err
	}

	payload, err := fetchHistoryKline(ctx, client, endpoint, secid, options, false)
	if err != nil {
		return nil, err
	}
	rows := make([]types.HistoryKline, 0, len(payload.Data.Klines))
	for _, line := range payload.Data.Klines {
		rows = append(rows, parseHistoryKlineCSV(line, normalized.Code))
	}
	return rows, nil
}

// GetHKHistoryKline fetches HK historical K-line rows.
func GetHKHistoryKline(ctx context.Context, client KlineClient, symbol string, endpoint string, options HistoryKlineOptions) ([]types.HKHistoryKline, error) {
	options = normalizeHistoryKlineOptions(options)
	if err := validateHistoryKlineOptions(options); err != nil {
		return nil, err
	}
	normalized, err := symbols.Normalize(symbol, &symbols.Hint{Market: symbols.MarketHK})
	if err != nil {
		return nil, err
	}
	secid, err := symbols.ToEastmoneySecIDE(normalized)
	if err != nil {
		return nil, err
	}
	payload, err := fetchHistoryKline(ctx, client, endpoint, secid, options, true)
	if err != nil {
		return nil, err
	}
	code := payload.Data.Code
	if code == "" {
		code = normalized.Code
	}
	rows := make([]types.HKHistoryKline, 0, len(payload.Data.Klines))
	for _, line := range payload.Data.Klines {
		rows = append(rows, types.HKHistoryKline{
			ForeignHistoryKline: parseForeignHistoryKlineCSV(line, code, payload.Data.Name, "Asia/Hong_Kong"),
			Currency:            "HKD",
			LotSize:             nil,
		})
	}
	return rows, nil
}

// GetUSHistoryKline fetches US historical K-line rows.
func GetUSHistoryKline(ctx context.Context, client KlineClient, symbol string, endpoint string, options HistoryKlineOptions) ([]types.USHistoryKline, error) {
	options = normalizeHistoryKlineOptions(options)
	if err := validateHistoryKlineOptions(options); err != nil {
		return nil, err
	}
	secid := strings.TrimSpace(symbol)
	if secid == "" {
		return nil, invalidSymbolError(symbol)
	}
	payload, err := fetchHistoryKline(ctx, client, endpoint, secid, options, true)
	if err != nil {
		return nil, err
	}
	code := payload.Data.Code
	if code == "" {
		code = usTickerFromSecid(secid)
	}
	rows := make([]types.USHistoryKline, 0, len(payload.Data.Klines))
	for _, line := range payload.Data.Klines {
		rows = append(rows, types.USHistoryKline{
			ForeignHistoryKline: parseForeignHistoryKlineCSV(line, code, payload.Data.Name, "America/New_York"),
			Currency:            "USD",
		})
	}
	return rows, nil
}

func fetchHistoryKline(ctx context.Context, client KlineClient, endpoint string, secid string, options HistoryKlineOptions, includeLimit bool) (historyKlineResponse, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f116")
	params.Set("ut", emPushToken)
	params.Set("klt", periodCode(options.Period))
	params.Set("fqt", adjustCode(options.Adjust))
	params.Set("secid", secid)
	params.Set("beg", options.StartDate)
	params.Set("end", options.EndDate)
	if includeLimit {
		params.Set("lmt", "1000000")
	}

	var payload historyKlineResponse
	err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload)
	return payload, err
}

func normalizeHistoryKlineOptions(options HistoryKlineOptions) HistoryKlineOptions {
	if options.Period == "" {
		options.Period = KlinePeriodDaily
	}
	if options.Adjust == "" {
		options.Adjust = AdjustQFQ
	}
	if options.StartDate == "" {
		options.StartDate = "19700101"
	}
	if options.EndDate == "" {
		options.EndDate = "20500101"
	}
	return options
}

func validateHistoryKlineOptions(options HistoryKlineOptions) error {
	switch options.Period {
	case KlinePeriodDaily, KlinePeriodWeekly, KlinePeriodMonthly:
	default:
		return invalidArgumentError(fmt.Sprintf("invalid period %q", options.Period))
	}
	switch options.Adjust {
	case "", AdjustNone, AdjustQFQ, AdjustHFQ:
	default:
		return invalidArgumentError(fmt.Sprintf("invalid adjust %q", options.Adjust))
	}
	return nil
}

func periodCode(period KlinePeriod) string {
	switch period {
	case KlinePeriodWeekly:
		return "102"
	case KlinePeriodMonthly:
		return "103"
	default:
		return "101"
	}
}

func adjustCode(adjust AdjustType) string {
	switch adjust {
	case AdjustHFQ:
		return "2"
	case "", AdjustNone:
		return "0"
	default:
		return "1"
	}
}

func parseHistoryKlineCSV(line string, code string) types.HistoryKline {
	fields := strings.Split(line, ",")
	timeMeta := timeutil.BuildTimeMeta(field(fields, 0), timeutil.MarketTZ.CN)
	return types.HistoryKline{
		Date:          field(fields, 0),
		Timestamp:     timeMeta.Timestamp,
		TZ:            string(timeMeta.TZ),
		Code:          code,
		Open:          toNumber(field(fields, 1)),
		Close:         toNumber(field(fields, 2)),
		High:          toNumber(field(fields, 3)),
		Low:           toNumber(field(fields, 4)),
		Volume:        toNumber(field(fields, 5)),
		Amount:        toNumber(field(fields, 6)),
		Amplitude:     toNumber(field(fields, 7)),
		ChangePercent: toNumber(field(fields, 8)),
		Change:        toNumber(field(fields, 9)),
		TurnoverRate:  toNumber(field(fields, 10)),
	}
}

func parseForeignHistoryKlineCSV(line string, code string, name string, tzName string) types.ForeignHistoryKline {
	row := parseHistoryKlineCSV(line, code)
	return types.ForeignHistoryKline{
		Date:          row.Date,
		Timestamp:     historyDateTimestamp(row.Date, tzName),
		TZ:            tzName,
		Code:          code,
		Name:          name,
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

func historyDateTimestamp(date string, tzName string) *int64 {
	location, err := time.LoadLocation(tzName)
	if err != nil {
		return nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", date, location)
	if err != nil {
		return nil
	}
	timestamp := parsed.UnixMilli()
	return &timestamp
}

func usTickerFromSecid(secid string) string {
	parts := strings.SplitN(secid, ".", 2)
	if len(parts) == 2 && parts[1] != "" {
		return parts[1]
	}
	return secid
}

func field(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return fields[index]
}

func toNumber(value string) *float64 {
	if value == "" || value == "-" {
		return nil
	}
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &n
}

// GetMinuteKline fetches CN minute K-line or timeline rows.
func GetMinuteKline(ctx context.Context, client KlineClient, symbol string, klineEndpoint string, trendsEndpoint string, options MinuteKlineOptions) (types.MinuteKlineResult, error) {
	options = normalizeMinuteKlineOptions(options)
	if err := validateMinuteKlineOptions(options); err != nil {
		return types.MinuteKlineResult{}, err
	}
	normalized, err := symbols.Normalize(symbol, &symbols.Hint{Market: symbols.MarketCN})
	if err != nil {
		return types.MinuteKlineResult{}, err
	}
	secid, err := symbols.ToEastmoneySecIDE(normalized)
	if err != nil {
		return types.MinuteKlineResult{}, err
	}

	if options.Period == MinutePeriodOne {
		return fetchMinuteTimeline(ctx, client, normalized.Code, secid, trendsEndpoint, options)
	}
	return fetchMinuteKlineRows(ctx, client, normalized.Code, secid, klineEndpoint, options)
}

// GetHKMinuteKline fetches HK minute K-line or timeline rows.
func GetHKMinuteKline(ctx context.Context, client KlineClient, symbol string, klineEndpoint string, trendsEndpoint string, options MinuteKlineOptions) (types.HKMinuteKlineResult, error) {
	options = normalizeForeignMinuteKlineOptions(options)
	if err := validateMinuteKlineOptions(options); err != nil {
		return types.HKMinuteKlineResult{}, err
	}
	normalized, err := symbols.Normalize(symbol, &symbols.Hint{Market: symbols.MarketHK})
	if err != nil {
		return types.HKMinuteKlineResult{}, err
	}
	secid, err := symbols.ToEastmoneySecIDE(normalized)
	if err != nil {
		return types.HKMinuteKlineResult{}, err
	}

	if options.Period == MinutePeriodOne {
		timeline, err := fetchForeignMinuteTimeline(ctx, client, normalized.Code, secid, trendsEndpoint, options, "Asia/Hong_Kong")
		if err != nil {
			return types.HKMinuteKlineResult{}, err
		}
		rows := make([]types.HKMinuteTimeline, 0, len(timeline))
		for _, row := range timeline {
			rows = append(rows, types.HKMinuteTimeline{ForeignMinuteTimeline: row, Currency: "HKD"})
		}
		return types.HKMinuteKlineResult{Timeline: rows}, nil
	}

	klines, err := fetchForeignMinuteKlineRows(ctx, client, normalized.Code, secid, klineEndpoint, options, "Asia/Hong_Kong")
	if err != nil {
		return types.HKMinuteKlineResult{}, err
	}
	rows := make([]types.HKMinuteKline, 0, len(klines))
	for _, row := range klines {
		rows = append(rows, types.HKMinuteKline{ForeignMinuteKline: row, Currency: "HKD"})
	}
	return types.HKMinuteKlineResult{Klines: rows}, nil
}

// GetUSMinuteKline fetches US minute K-line or timeline rows.
func GetUSMinuteKline(ctx context.Context, client KlineClient, symbol string, klineEndpoint string, trendsEndpoint string, options MinuteKlineOptions) (types.USMinuteKlineResult, error) {
	options = normalizeForeignMinuteKlineOptions(options)
	if err := validateMinuteKlineOptions(options); err != nil {
		return types.USMinuteKlineResult{}, err
	}
	secid := strings.TrimSpace(symbol)
	if secid == "" {
		return types.USMinuteKlineResult{}, invalidSymbolError(symbol)
	}
	code := usTickerFromSecid(secid)

	if options.Period == MinutePeriodOne {
		timeline, err := fetchForeignMinuteTimeline(ctx, client, code, secid, trendsEndpoint, options, "America/New_York")
		if err != nil {
			return types.USMinuteKlineResult{}, err
		}
		rows := make([]types.USMinuteTimeline, 0, len(timeline))
		for _, row := range timeline {
			rows = append(rows, types.USMinuteTimeline{ForeignMinuteTimeline: row, Currency: "USD"})
		}
		return types.USMinuteKlineResult{Timeline: rows}, nil
	}

	klines, err := fetchForeignMinuteKlineRows(ctx, client, code, secid, klineEndpoint, options, "America/New_York")
	if err != nil {
		return types.USMinuteKlineResult{}, err
	}
	rows := make([]types.USMinuteKline, 0, len(klines))
	for _, row := range klines {
		rows = append(rows, types.USMinuteKline{ForeignMinuteKline: row, Currency: "USD"})
	}
	return types.USMinuteKlineResult{Klines: rows}, nil
}

func normalizeMinuteKlineOptions(options MinuteKlineOptions) MinuteKlineOptions {
	if options.Period == "" {
		options.Period = MinutePeriodOne
	}
	if options.Adjust == "" {
		options.Adjust = AdjustQFQ
	}
	if options.StartDate == "" {
		options.StartDate = "1979-09-01 09:32:00"
	}
	if options.EndDate == "" {
		options.EndDate = "2222-01-01 09:32:00"
	}
	if options.NDays <= 0 {
		options.NDays = 5
	}
	return options
}

func normalizeForeignMinuteKlineOptions(options MinuteKlineOptions) MinuteKlineOptions {
	if options.Period == "" {
		options.Period = MinutePeriodOne
	}
	if options.Adjust == "" {
		options.Adjust = AdjustQFQ
	}
	if options.StartDate == "" {
		options.StartDate = "1979-09-01 09:32:00"
	}
	if options.EndDate == "" {
		options.EndDate = "2222-01-01 09:32:00"
	}
	if options.NDays <= 0 {
		options.NDays = 1
	}
	return options
}

func validateMinuteKlineOptions(options MinuteKlineOptions) error {
	switch options.Period {
	case MinutePeriodOne, MinutePeriodFive, MinutePeriodFifteen, MinutePeriodThirty, MinutePeriodSixty:
	default:
		return invalidArgumentError(fmt.Sprintf("invalid minute period %q", options.Period))
	}
	switch options.Adjust {
	case AdjustNone, AdjustQFQ, AdjustHFQ:
	default:
		return invalidArgumentError(fmt.Sprintf("invalid adjust %q", options.Adjust))
	}
	return nil
}

func fetchMinuteTimeline(ctx context.Context, client KlineClient, code string, secid string, endpoint string, options MinuteKlineOptions) (types.MinuteKlineResult, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	params.Set("ut", emPushToken)
	params.Set("ndays", strconv.Itoa(options.NDays))
	params.Set("iscr", "0")
	params.Set("secid", secid)

	var payload minuteTimelineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return types.MinuteKlineResult{}, err
	}
	trends, err := decodeStringArray(payload.Data.Trends)
	if err != nil {
		return types.MinuteKlineResult{}, err
	}
	start, end := minuteBounds(options.StartDate, options.EndDate)
	rows := make([]types.MinuteTimeline, 0, len(trends))
	for _, line := range trends {
		row := parseMinuteTimelineCSV(line, code)
		if row.Time >= start && row.Time <= end {
			rows = append(rows, row)
		}
	}
	return types.MinuteKlineResult{Timeline: rows}, nil
}

func fetchMinuteKlineRows(ctx context.Context, client KlineClient, code string, secid string, endpoint string, options MinuteKlineOptions) (types.MinuteKlineResult, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	params.Set("ut", emPushToken)
	params.Set("klt", string(options.Period))
	params.Set("fqt", adjustCode(options.Adjust))
	params.Set("secid", secid)
	params.Set("beg", "0")
	params.Set("end", "20500000")

	var payload minuteKlineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return types.MinuteKlineResult{}, err
	}
	klines, err := decodeStringArray(payload.Data.Klines)
	if err != nil {
		return types.MinuteKlineResult{}, err
	}
	start, end := minuteBounds(options.StartDate, options.EndDate)
	rows := make([]types.MinuteKline, 0, len(klines))
	for _, line := range klines {
		row := parseMinuteKlineCSV(line, code)
		if row.Time >= start && row.Time <= end {
			rows = append(rows, row)
		}
	}
	return types.MinuteKlineResult{Klines: rows}, nil
}

func fetchForeignMinuteTimeline(ctx context.Context, client KlineClient, code string, secid string, endpoint string, options MinuteKlineOptions, targetTZ string) ([]types.ForeignMinuteTimeline, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	params.Set("ut", emPushToken)
	params.Set("ndays", strconv.Itoa(options.NDays))
	params.Set("iscr", "0")
	params.Set("secid", secid)

	var payload minuteTimelineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	trends, err := decodeStringArray(payload.Data.Trends)
	if err != nil {
		return nil, err
	}
	start, end := minuteBounds(options.StartDate, options.EndDate)
	rows := make([]types.ForeignMinuteTimeline, 0, len(trends))
	for _, line := range trends {
		row := parseForeignMinuteTimelineCSV(line, code, targetTZ)
		if row.Time >= start && row.Time <= end {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func fetchForeignMinuteKlineRows(ctx context.Context, client KlineClient, code string, secid string, endpoint string, options MinuteKlineOptions, targetTZ string) ([]types.ForeignMinuteKline, error) {
	params := url.Values{}
	params.Set("fields1", "f1,f2,f3,f4,f5,f6")
	params.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	params.Set("ut", emPushToken)
	params.Set("klt", string(options.Period))
	params.Set("fqt", adjustCode(options.Adjust))
	params.Set("secid", secid)
	params.Set("beg", "0")
	params.Set("end", "20500000")

	var payload minuteKlineResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	klines, err := decodeStringArray(payload.Data.Klines)
	if err != nil {
		return nil, err
	}
	start, end := minuteBounds(options.StartDate, options.EndDate)
	rows := make([]types.ForeignMinuteKline, 0, len(klines))
	for _, line := range klines {
		row := parseForeignMinuteKlineCSV(line, code, targetTZ)
		if row.Time >= start && row.Time <= end {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func minuteBounds(startDate, endDate string) (string, string) {
	start := strings.ReplaceAll(startDate, "T", " ")
	if len(start) > 16 {
		start = start[:16]
	}
	end := strings.ReplaceAll(endDate, "T", " ")
	if len(end) > 16 {
		end = end[:16]
	}
	if len(end) == 10 {
		end += " 23:59"
	}
	return start, end
}

func parseMinuteTimelineCSV(line string, code string) types.MinuteTimeline {
	fields := strings.Split(line, ",")
	timeMeta := timeutil.BuildTimeMeta(field(fields, 0), timeutil.MarketTZ.CN)
	return types.MinuteTimeline{
		Time:      field(fields, 0),
		Timestamp: timeMeta.Timestamp,
		TZ:        string(timeMeta.TZ),
		Code:      code,
		Open:      toNumber(field(fields, 1)),
		Close:     toNumber(field(fields, 2)),
		High:      toNumber(field(fields, 3)),
		Low:       toNumber(field(fields, 4)),
		Volume:    toNumber(field(fields, 5)),
		Amount:    toNumber(field(fields, 6)),
		AvgPrice:  toNumber(field(fields, 7)),
	}
}

func parseMinuteKlineCSV(line string, code string) types.MinuteKline {
	row := parseHistoryKlineCSV(line, code)
	return types.MinuteKline{
		Time:          row.Date,
		Timestamp:     row.Timestamp,
		TZ:            row.TZ,
		Code:          row.Code,
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

func parseForeignMinuteTimelineCSV(line string, code string, targetTZ string) types.ForeignMinuteTimeline {
	row := parseMinuteTimelineCSV(line, code)
	localTime, timestamp := marketLocalMinute(row.Time, targetTZ)
	return types.ForeignMinuteTimeline{
		Time:      localTime,
		Timestamp: timestamp,
		TZ:        targetTZ,
		Code:      code,
		Open:      row.Open,
		Close:     row.Close,
		High:      row.High,
		Low:       row.Low,
		Volume:    row.Volume,
		Amount:    row.Amount,
		AvgPrice:  row.AvgPrice,
	}
}

func parseForeignMinuteKlineCSV(line string, code string, targetTZ string) types.ForeignMinuteKline {
	row := parseMinuteKlineCSV(line, code)
	localTime, timestamp := marketLocalMinute(row.Time, targetTZ)
	return types.ForeignMinuteKline{
		Time:          localTime,
		Timestamp:     timestamp,
		TZ:            targetTZ,
		Code:          code,
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

func marketLocalMinute(rawTime string, targetTZ string) (string, *int64) {
	source, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return rawTime, nil
	}
	parsed, err := parseMinuteTimeInLocation(rawTime, source)
	if err != nil {
		return rawTime, nil
	}
	timestamp := parsed.UnixMilli()
	target, err := time.LoadLocation(targetTZ)
	if err != nil {
		return rawTime, &timestamp
	}
	return parsed.In(target).Format("2006-01-02 15:04"), &timestamp
}

func parseMinuteTimeInLocation(value string, location *time.Location) (time.Time, error) {
	normalized := strings.ReplaceAll(value, "T", " ")
	if len(normalized) > 16 {
		normalized = normalized[:16]
	}
	return time.ParseInLocation("2006-01-02 15:04", normalized, location)
}
