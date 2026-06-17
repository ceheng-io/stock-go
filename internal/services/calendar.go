package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ceheng-io/stock-go/internal/core"
)

// MarketStatus 是市场实时状态。
type MarketStatus string

const (
	MarketStatusPreMarket  MarketStatus = "pre_market"
	MarketStatusOpen       MarketStatus = "open"
	MarketStatusLunchBreak MarketStatus = "lunch_break"
	MarketStatusAfterHours MarketStatus = "after_hours"
	MarketStatusClosed     MarketStatus = "closed"
)

// SupportedMarket 是市场状态判断支持的市场。
type SupportedMarket string

const (
	MarketA  SupportedMarket = "A"
	MarketHK SupportedMarket = "HK"
	MarketUS SupportedMarket = "US"
)

// CalendarQuoteService 是 CalendarService 依赖的行情日历接口。
type CalendarQuoteService interface {
	TradingCalendar(context.Context) ([]string, error)
}

// CalendarService 编排交易日历和市场状态能力。
type CalendarService struct {
	quotes CalendarQuoteService
}

// NewCalendarService 创建 CalendarService。
func NewCalendarService(quotes CalendarQuoteService) *CalendarService {
	return &CalendarService{quotes: quotes}
}

// IsTradingDay 判断给定日期是否为 A 股交易日。
func (s *CalendarService) IsTradingDay(ctx context.Context, date string) (bool, error) {
	target := normalizeCalendarDate(date, "Asia/Shanghai")
	calendar, err := s.sortedCalendar(ctx)
	if err != nil {
		return false, err
	}
	idx := sort.SearchStrings(calendar, target)
	return idx < len(calendar) && calendar[idx] == target, nil
}

// NextTradingDay 返回给定日期之后的下一个 A 股交易日。
func (s *CalendarService) NextTradingDay(ctx context.Context, date string) (string, error) {
	target := normalizeCalendarDate(date, "Asia/Shanghai")
	calendar, err := s.sortedCalendar(ctx)
	if err != nil {
		return "", err
	}
	idx := sort.SearchStrings(calendar, target)
	if idx < len(calendar) && calendar[idx] == target {
		idx++
	}
	if idx >= len(calendar) {
		return "", invalidArgumentError(fmt.Sprintf("next trading day after %s is out of calendar range", target))
	}
	return calendar[idx], nil
}

// PrevTradingDay 返回给定日期之前的上一个 A 股交易日。
func (s *CalendarService) PrevTradingDay(ctx context.Context, date string) (string, error) {
	target := normalizeCalendarDate(date, "Asia/Shanghai")
	calendar, err := s.sortedCalendar(ctx)
	if err != nil {
		return "", err
	}
	idx := sort.SearchStrings(calendar, target)
	prevIdx := idx - 1
	if prevIdx < 0 {
		return "", invalidArgumentError(fmt.Sprintf("previous trading day before %s is out of calendar range", target))
	}
	return calendar[prevIdx], nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

// MarketStatus 返回指定市场在给定时刻的时段状态。
func (s *CalendarService) MarketStatus(market SupportedMarket, now time.Time) MarketStatus {
	if market == "" {
		market = MarketA
	}
	session := marketSessions[market]
	if session.location == "" {
		return MarketStatusClosed
	}
	location, err := time.LoadLocation(session.location)
	if err != nil {
		return MarketStatusClosed
	}
	local := now.In(location)
	weekday := int(local.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	if !containsInt(session.tradingWeekdays, weekday) {
		return MarketStatusClosed
	}
	minutes := local.Hour()*60 + local.Minute()
	dayStart := session.open[0][0]
	dayEnd := session.open[len(session.open)-1][1]
	if minutes < dayStart {
		return MarketStatusPreMarket
	}
	if minutes >= dayEnd {
		return MarketStatusAfterHours
	}
	for _, period := range session.open {
		if minutes >= period[0] && minutes < period[1] {
			return MarketStatusOpen
		}
	}
	if session.lunchBreak != nil && minutes >= session.lunchBreak[0] && minutes < session.lunchBreak[1] {
		return MarketStatusLunchBreak
	}
	return MarketStatusClosed
}

func (s *CalendarService) sortedCalendar(ctx context.Context) ([]string, error) {
	calendar, err := s.quotes.TradingCalendar(ctx)
	if err != nil {
		return nil, err
	}
	rows := append([]string(nil), calendar...)
	sort.Strings(rows)
	return rows, nil
}

func normalizeCalendarDate(input string, locationName string) string {
	text := strings.TrimSpace(input)
	if text == "" {
		location, err := time.LoadLocation(locationName)
		if err != nil {
			return time.Now().Format("2006-01-02")
		}
		return time.Now().In(location).Format("2006-01-02")
	}
	if len(text) == 8 && strings.IndexByte(text, '-') < 0 {
		return text[:4] + "-" + text[4:6] + "-" + text[6:8]
	}
	return text
}

type marketSession struct {
	location        string
	open            [][2]int
	lunchBreak      *[2]int
	tradingWeekdays []int
}

func hm(hour int, minute int) int {
	return hour*60 + minute
}

var marketSessions = map[SupportedMarket]marketSession{
	MarketA: {
		location: "Asia/Shanghai",
		open: [][2]int{
			{hm(9, 30), hm(11, 30)},
			{hm(13, 0), hm(15, 0)},
		},
		lunchBreak:      &[2]int{hm(11, 30), hm(13, 0)},
		tradingWeekdays: []int{1, 2, 3, 4, 5},
	},
	MarketHK: {
		location: "Asia/Hong_Kong",
		open: [][2]int{
			{hm(9, 30), hm(12, 0)},
			{hm(13, 0), hm(16, 0)},
		},
		lunchBreak:      &[2]int{hm(12, 0), hm(13, 0)},
		tradingWeekdays: []int{1, 2, 3, 4, 5},
	},
	MarketUS: {
		location:        "America/New_York",
		open:            [][2]int{{hm(9, 30), hm(16, 0)}},
		tradingWeekdays: []int{1, 2, 3, 4, 5},
	},
}

func containsInt(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
