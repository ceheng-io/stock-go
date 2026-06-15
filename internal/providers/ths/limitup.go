package ths

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

const (
	defaultLimitUpFields = "199112,10,9001,330323,330324,330325,9002,330329,133971,133970,1968584,3475914,9003,9004"
	defaultLimitUpFilter = "HS,GEM2STAR"
)

// LimitUpClient is the minimal request client required by Tonghuashun providers.
type LimitUpClient interface {
	GetJSON(context.Context, string, any) error
}

type LimitUpPoolOptions = types.THSLimitUpPoolOptions
type LimitUpOrderField = types.THSLimitUpOrderField
type LimitUpOrderType = types.THSLimitUpOrderType

const (
	LimitUpOrderFirstLimitUpTime LimitUpOrderField = types.THSLimitUpOrderFirstLimitUpTime
	LimitUpOrderLastLimitUpTime  LimitUpOrderField = types.THSLimitUpOrderLastLimitUpTime
	LimitUpOrderOpenNum          LimitUpOrderField = types.THSLimitUpOrderOpenNum
	LimitUpOrderDesc             LimitUpOrderType  = types.THSLimitUpOrderDesc
	LimitUpOrderAsc              LimitUpOrderType  = types.THSLimitUpOrderAsc
)

type limitUpPoolResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	Data       struct {
		Page           thsLimitUpPage      `json:"page"`
		Info           []thsLimitUpRawItem `json:"info"`
		LimitUpCount   thsLimitStatGroup   `json:"limit_up_count"`
		LimitDownCount thsLimitStatGroup   `json:"limit_down_count"`
		Date           string              `json:"date"`
		Message        string              `json:"msg"`
		TradeStatus    thsTradeStatus      `json:"trade_status"`
	} `json:"data"`
}

type thsLimitUpPage struct {
	Limit int `json:"limit"`
	Total int `json:"total"`
	Count int `json:"count"`
	Page  int `json:"page"`
}

type thsLimitStatGroup struct {
	Today     thsLimitStat `json:"today"`
	Yesterday thsLimitStat `json:"yesterday"`
}

type thsLimitStat struct {
	Num        int      `json:"num"`
	HistoryNum int      `json:"history_num"`
	Rate       *float64 `json:"rate"`
	OpenNum    int      `json:"open_num"`
}

type thsTradeStatus struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type thsLimitUpRawItem struct {
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	Latest             *float64  `json:"latest"`
	ChangeRate         *float64  `json:"change_rate"`
	FirstLimitUpTime   string    `json:"first_limit_up_time"`
	LastLimitUpTime    string    `json:"last_limit_up_time"`
	OpenNum            *int      `json:"open_num"`
	LimitUpType        string    `json:"limit_up_type"`
	OrderVolume        *float64  `json:"order_volume"`
	OrderAmount        *float64  `json:"order_amount"`
	TurnoverRate       *float64  `json:"turnover_rate"`
	CurrencyValue      *float64  `json:"currency_value"`
	ReasonType         string    `json:"reason_type"`
	HighDays           string    `json:"high_days"`
	HighDaysValue      *int      `json:"high_days_value"`
	ChangeTag          string    `json:"change_tag"`
	MarketType         string    `json:"market_type"`
	MarketID           *int      `json:"market_id"`
	IsNew              *int      `json:"is_new"`
	IsAgainLimit       *int      `json:"is_again_limit"`
	LimitUpSuccessRate *float64  `json:"limit_up_suc_rate"`
	TimePreview        []float64 `json:"time_preview"`
}

// GetLimitUpPool fetches Tonghuashun limit-up pool rows.
func GetLimitUpPool(ctx context.Context, client LimitUpClient, endpoint string, options LimitUpPoolOptions) (types.THSLimitUpPoolResult, error) {
	normalized := normalizeLimitUpOptions(options)
	params := url.Values{}
	params.Set("page", strconv.Itoa(normalized.Page))
	params.Set("limit", strconv.Itoa(normalized.Limit))
	params.Set("field", defaultLimitUpFields)
	params.Set("filter", normalized.Filter)
	params.Set("order_field", string(normalized.OrderField))
	params.Set("order_type", string(normalized.OrderType))
	params.Set("date", normalizeLimitUpDate(normalized.Date))
	params.Set("_", strconv.FormatInt(time.Now().UnixMilli(), 10))

	var payload limitUpPoolResponse
	if err := client.GetJSON(ctx, strings.TrimRight(endpoint, "?")+"?"+params.Encode(), &payload); err != nil {
		return types.THSLimitUpPoolResult{}, err
	}
	if payload.StatusCode != 0 {
		return types.THSLimitUpPoolResult{}, core.NewCodedError("UPSTREAM_ERROR", fmt.Sprintf("ths limit up pool upstream status %d: %s", payload.StatusCode, payload.StatusMsg), nil)
	}
	items := make([]types.THSLimitUpItem, 0, len(payload.Data.Info))
	for _, item := range payload.Data.Info {
		items = append(items, parseLimitUpItem(item))
	}
	return types.THSLimitUpPoolResult{
		Page:           parseLimitUpPage(payload.Data.Page),
		Items:          items,
		LimitUpCount:   parseLimitStatGroup(payload.Data.LimitUpCount),
		LimitDownCount: parseLimitStatGroup(payload.Data.LimitDownCount),
		Date:           payload.Data.Date,
		Message:        payload.Data.Message,
		TradeStatus:    parseTradeStatus(payload.Data.TradeStatus),
	}, nil
}

func normalizeLimitUpOptions(options LimitUpPoolOptions) LimitUpPoolOptions {
	if options.Page <= 0 {
		options.Page = 1
	}
	if options.Limit <= 0 {
		options.Limit = 50
	}
	if options.Filter == "" {
		options.Filter = defaultLimitUpFilter
	}
	if options.OrderField == "" {
		options.OrderField = LimitUpOrderLastLimitUpTime
	}
	if options.OrderType == "" {
		options.OrderType = LimitUpOrderDesc
	}
	return options
}

func normalizeLimitUpDate(date string) string {
	date = strings.TrimSpace(date)
	if date == "" {
		return ""
	}
	return strings.ReplaceAll(date, "-", "")
}

func parseLimitUpItem(item thsLimitUpRawItem) types.THSLimitUpItem {
	firstTime, firstText := parseTHSTimestamp(item.FirstLimitUpTime)
	lastTime, lastText := parseTHSTimestamp(item.LastLimitUpTime)
	return types.THSLimitUpItem{
		Code:                 item.Code,
		Name:                 item.Name,
		Latest:               cloneFloat64Ptr(item.Latest),
		ChangeRate:           cloneFloat64Ptr(item.ChangeRate),
		FirstLimitUpTime:     firstTime,
		FirstLimitUpTimeText: firstText,
		LastLimitUpTime:      lastTime,
		LastLimitUpTimeText:  lastText,
		OpenNum:              cloneIntPtr(item.OpenNum),
		LimitUpType:          item.LimitUpType,
		OrderVolume:          cloneFloat64Ptr(item.OrderVolume),
		OrderAmount:          cloneFloat64Ptr(item.OrderAmount),
		TurnoverRate:         cloneFloat64Ptr(item.TurnoverRate),
		CurrencyValue:        cloneFloat64Ptr(item.CurrencyValue),
		ReasonType:           item.ReasonType,
		HighDays:             item.HighDays,
		HighDaysValue:        cloneIntPtr(item.HighDaysValue),
		ChangeTag:            item.ChangeTag,
		MarketType:           item.MarketType,
		MarketID:             cloneIntPtr(item.MarketID),
		IsNew:                cloneIntPtr(item.IsNew),
		IsAgainLimit:         cloneIntPtr(item.IsAgainLimit),
		LimitUpSuccessRate:   cloneFloat64Ptr(item.LimitUpSuccessRate),
		TimePreview:          append([]float64{}, item.TimePreview...),
	}
}

func parseTHSTimestamp(value string) (*int64, string) {
	if value == "" {
		return nil, ""
	}
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, ""
	}
	return &seconds, time.Unix(seconds, 0).UTC().Add(8 * time.Hour).Format("15:04:05")
}

func parseLimitUpPage(page thsLimitUpPage) types.THSLimitUpPage {
	return types.THSLimitUpPage{
		Limit: page.Limit,
		Total: page.Total,
		Count: page.Count,
		Page:  page.Page,
	}
}

func parseLimitStatGroup(group thsLimitStatGroup) types.THSLimitStatGroup {
	return types.THSLimitStatGroup{
		Today:     parseLimitStat(group.Today),
		Yesterday: parseLimitStat(group.Yesterday),
	}
}

func parseLimitStat(stat thsLimitStat) types.THSLimitStat {
	return types.THSLimitStat{
		Num:        stat.Num,
		HistoryNum: stat.HistoryNum,
		Rate:       cloneFloat64Ptr(stat.Rate),
		OpenNum:    stat.OpenNum,
	}
}

func parseTradeStatus(status thsTradeStatus) types.THSTradeStatus {
	return types.THSTradeStatus{
		ID:        status.ID,
		Name:      status.Name,
		StartTime: status.StartTime,
		EndTime:   status.EndTime,
	}
}

func cloneFloat64Ptr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
