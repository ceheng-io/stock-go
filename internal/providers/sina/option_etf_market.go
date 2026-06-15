package sina

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/ceheng.io/stock-go/types"
)

const sinaJSONPCallback = "ceheng_jsonp"

type etfOptionMinuteItem struct {
	Time         string `json:"i"`
	Price        string `json:"p"`
	Volume       string `json:"v"`
	OpenInterest string `json:"t"`
	AvgPrice     string `json:"a"`
	Date         string `json:"d"`
}

type etfOptionMinuteResponse struct {
	Result struct {
		Data json.RawMessage `json:"data"`
	} `json:"result"`
}

type etfOption5DayMinuteResponse struct {
	Result struct {
		Data json.RawMessage `json:"data"`
	} `json:"result"`
}

// GetETFOptionMinute 获取新浪 ETF 期权当日分钟行情。
func GetETFOptionMinute(ctx context.Context, client JSONPClient, endpoint string, code string) ([]types.OptionMinute, error) {
	params := url.Values{}
	params.Set("symbol", etfOptionSymbol(code))

	var payload etfOptionMinuteResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return nil, err
	}
	return parseETFOptionMinutePayload(payload.Result.Data)
}

// GetETFOptionDailyKline 获取新浪 ETF 期权历史日 K 线。
func GetETFOptionDailyKline(ctx context.Context, client JSONPClient, endpoint string, code string) ([]types.OptionKline, error) {
	params := url.Values{}
	params.Set("symbol", etfOptionSymbol(code))

	return getSinaOptionKlinesJSONP(ctx, client, endpoint, params)
}

func getSinaOptionKlinesJSONP(ctx context.Context, client JSONPClient, endpoint string, params url.Values) ([]types.OptionKline, error) {
	var raw json.RawMessage
	if err := getSinaJSONP(ctx, client, sinaPathJSONPEndpoint(endpoint), params, &raw); err != nil {
		return nil, err
	}
	if !isJSONArray(raw) {
		return []types.OptionKline{}, nil
	}
	var payload []sinaOptionKlineItem
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	return parseOptionKlines(payload), nil
}

// GetETFOption5DayMinute 获取新浪 ETF 期权 5 日分钟行情。
func GetETFOption5DayMinute(ctx context.Context, client JSONPClient, endpoint string, code string) ([]types.OptionMinute, error) {
	params := url.Values{}
	params.Set("symbol", etfOptionSymbol(code))

	var payload etfOption5DayMinuteResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return nil, err
	}
	if !isJSONArray(payload.Result.Data) {
		return []types.OptionMinute{}, nil
	}
	var dayItems []json.RawMessage
	if err := json.Unmarshal(payload.Result.Data, &dayItems); err != nil {
		return nil, err
	}
	var rows []types.OptionMinute
	for _, rawDayItems := range dayItems {
		if !isJSONArray(rawDayItems) {
			continue
		}
		parsed, err := parseETFOptionMinutePayload(rawDayItems)
		if err != nil {
			return nil, err
		}
		rows = append(rows, parsed...)
	}
	return rows, nil
}

func parseETFOptionMinutePayload(raw json.RawMessage) ([]types.OptionMinute, error) {
	if !isJSONArray(raw) {
		return []types.OptionMinute{}, nil
	}
	var items []etfOptionMinuteItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return parseETFOptionMinuteList(items), nil
}

func parseETFOptionMinuteList(items []etfOptionMinuteItem) []types.OptionMinute {
	rows := make([]types.OptionMinute, 0, len(items))
	currentDate := ""
	for _, item := range items {
		if item.Date != "" {
			currentDate = item.Date
		}
		rows = append(rows, types.OptionMinute{
			Time:         item.Time,
			Date:         currentDate,
			Price:        sinaNumber(item.Price),
			Volume:       sinaNumber(item.Volume),
			OpenInterest: sinaNumber(item.OpenInterest),
			AvgPrice:     sinaNumber(item.AvgPrice),
		})
	}
	return rows
}

func etfOptionSymbol(code string) string {
	return "CON_OP_" + strings.TrimSpace(code)
}

func sinaPathJSONPEndpoint(endpoint string) string {
	return strings.Replace(endpoint, "{callback}", sinaJSONPCallback, 1)
}

func isJSONArray(raw json.RawMessage) bool {
	trimmed := strings.TrimSpace(string(raw))
	return strings.HasPrefix(trimmed, "[")
}

func sinaNumber(value string) *float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return nil
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &number
}
