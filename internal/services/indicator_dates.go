package services

import (
	"math"
	"sort"
	"strings"
	"time"
)

func actualStartDateByCalendar(startDate string, requiredBars int, calendar []string) string {
	if len(calendar) == 0 {
		return ""
	}
	normalized := indicatorHyphenDate(startDate)
	days := make([]string, 0, len(calendar))
	for _, date := range calendar {
		value := indicatorHyphenDate(date)
		if value != "" {
			days = append(days, value)
		}
	}
	if len(days) == 0 {
		return ""
	}
	sort.Strings(days)
	startIndex := sort.SearchStrings(days, normalized)
	if startIndex >= len(days) {
		startIndex = len(days) - 1
	}
	targetIndex := startIndex - requiredBars
	if targetIndex < 0 {
		targetIndex = 0
	}
	return indicatorCompactDate(days[targetIndex])
}

func actualStartDateByNaturalDays(startDate string, requiredBars int, ratio float64) string {
	compact := indicatorCompactDate(startDate)
	if len(compact) != 8 {
		return compact
	}
	date, err := time.ParseInLocation("20060102", compact, time.Local)
	if err != nil {
		return compact
	}
	naturalDays := int(math.Ceil(float64(requiredBars) * ratio))
	return date.AddDate(0, 0, -naturalDays).Format("20060102")
}

func indicatorRatio(market SupportedMarket) float64 {
	switch market {
	case MarketHK:
		return 1.46
	case MarketUS:
		return 1.45
	default:
		return 1.5
	}
}

func indicatorCompactDate(date string) string {
	return strings.ReplaceAll(strings.TrimSpace(date), "-", "")
}

func indicatorHyphenDate(date string) string {
	compact := indicatorCompactDate(date)
	if len(compact) != 8 {
		return strings.TrimSpace(date)
	}
	return compact[:4] + "-" + compact[4:6] + "-" + compact[6:8]
}

func indicatorDateKey(date string) string {
	return indicatorCompactDate(date)
}
