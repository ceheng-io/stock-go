package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

const cnTimezone = "Asia/Shanghai"

// TimelineClient 是腾讯分时 provider 所需的最小请求客户端接口。
type TimelineClient interface {
	GetText(context.Context, string) (string, error)
}

type todayTimelineResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data map[string]struct {
		Data struct {
			Data []string `json:"data"`
			Date string   `json:"date"`
		} `json:"data"`
		QT map[string][]string `json:"qt"`
	} `json:"data"`
}

// GetTodayTimeline 获取腾讯当日分时走势数据。
func GetTodayTimeline(ctx context.Context, client TimelineClient, endpoint string, code string) (types.TodayTimelineResponse, error) {
	requestURL := endpoint + "?code=" + url.QueryEscape(code)
	text, err := client.GetText(ctx, requestURL)
	if err != nil {
		return types.TodayTimelineResponse{}, err
	}
	var payload todayTimelineResponse
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return types.TodayTimelineResponse{}, err
	}
	if payload.Code != 0 {
		msg := payload.Msg
		if msg == "" {
			msg = "API error"
		}
		return types.TodayTimelineResponse{}, upstreamError(fmt.Sprintf("tencent timeline upstream error: %s", msg))
	}
	stockData, ok := payload.Data[code]
	if !ok {
		zero := 0.0
		return types.TodayTimelineResponse{Code: code, TZ: cnTimezone, PreClose: &zero, Data: []types.TodayTimeline{}}, nil
	}
	date := stockData.Data.Date
	quoteFields := stockData.QT[code]
	preCloseValue := safeNumber(field(quoteFields, 4))
	rows := parseTodayTimelineRows(date, stockData.Data.Data)
	return types.TodayTimelineResponse{
		Code:      code,
		Date:      date,
		Timestamp: marketTimeMillis(date, ""),
		TZ:        cnTimezone,
		PreClose:  &preCloseValue,
		Data:      rows,
	}, nil
}

func upstreamError(message string) error {
	return core.NewCodedError("UPSTREAM_ERROR", message, nil)
}

func parseTodayTimelineRows(date string, rawRows []string) []types.TodayTimeline {
	isVolumeInLots := false
	for _, line := range rawRows {
		parts := strings.Fields(line)
		price := safeNumber(field(parts, 1))
		volume := safeNumber(field(parts, 2))
		amount := safeNumber(field(parts, 3))
		if volume > 0 && price > 0 {
			if amount/volume > price*50 {
				isVolumeInLots = true
			}
			break
		}
	}
	rows := make([]types.TodayTimeline, 0, len(rawRows))
	for _, line := range rawRows {
		parts := strings.Fields(line)
		rawTime := field(parts, 0)
		tickTime := formatHHMM(rawTime)
		rawVolume := safeNumber(field(parts, 2))
		amount := safeNumber(field(parts, 3))
		volume := rawVolume
		if isVolumeInLots {
			volume *= 100
		}
		avgPrice := 0.0
		if volume > 0 {
			avgPrice = math.Round(amount/volume*100) / 100
		}
		rows = append(rows, types.TodayTimeline{
			Time:      tickTime,
			Timestamp: marketTimeMillis(date, tickTime),
			TZ:        cnTimezone,
			Price:     safeNumber(field(parts, 1)),
			Volume:    volume,
			Amount:    amount,
			AvgPrice:  avgPrice,
		})
	}
	return rows
}

func formatHHMM(value string) string {
	if len(value) >= 4 {
		return value[:2] + ":" + value[2:4]
	}
	return value
}

func marketTimeMillis(date string, hm string) *int64 {
	if date == "" {
		return nil
	}
	layout := "20060102"
	value := date
	if strings.Contains(date, "-") {
		layout = "2006-01-02"
	}
	if hm != "" {
		value = date + " " + hm
		if strings.Contains(date, "-") {
			layout = "2006-01-02 15:04"
		} else {
			layout = "20060102 15:04"
		}
	}
	location, err := time.LoadLocation(cnTimezone)
	if err != nil {
		return nil
	}
	parsed, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		return nil
	}
	millis := parsed.UnixMilli()
	return &millis
}
