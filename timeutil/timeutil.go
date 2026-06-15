package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MarketTz 是市场使用的 IANA 时区名。
type MarketTz string

// MarketTZ 保存常用市场时区。
var MarketTZ = struct {
	CN MarketTz
	HK MarketTz
	US MarketTz
}{
	CN: "Asia/Shanghai",
	HK: "Asia/Hong_Kong",
	US: "America/New_York",
}

// TimeMeta 是原始时间字符串对应的 UTC 毫秒时间戳和市场时区。
type TimeMeta struct {
	Timestamp *int64
	TZ        MarketTz
}

type wallClock struct {
	year   int
	month  time.Month
	day    int
	hour   int
	minute int
	second int
}

// ParseMarketTime 把市场本地时间字符串解析为 UTC unix 毫秒。
func ParseMarketTime(local string, tz MarketTz) (int64, bool) {
	wall, ok := parseWallClock(local)
	if !ok {
		return 0, false
	}
	return wallTimeToUnixMilli(wall, tz)
}

// BuildTimeMeta 构造时间元信息；解析失败时 Timestamp 为 nil。
func BuildTimeMeta(local string, tz MarketTz) TimeMeta {
	timestamp, ok := ParseMarketTime(local, tz)
	if !ok {
		return TimeMeta{TZ: tz}
	}
	return TimeMeta{Timestamp: &timestamp, TZ: tz}
}

// BuildTimeMetaFromDateAndTime 将基础日期和 HH:mm/HH:mm:ss 组合后构造时间元信息。
func BuildTimeMetaFromDateAndTime(baseDate string, hhmm string, tz MarketTz) TimeMeta {
	wall, ok := combineDateAndTime(baseDate, hhmm)
	if !ok {
		return TimeMeta{TZ: tz}
	}
	timestamp, ok := wallTimeToUnixMilli(wall, tz)
	if !ok {
		return TimeMeta{TZ: tz}
	}
	return TimeMeta{Timestamp: &timestamp, TZ: tz}
}

// FormatInTz 把 UTC 毫秒时间戳格式化为指定市场时区的 YYYY-MM-DD HH:mm。
func FormatInTz(epoch *int64, tz MarketTz) string {
	if epoch == nil {
		return ""
	}
	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return ""
	}
	return time.UnixMilli(*epoch).In(location).Format("2006-01-02 15:04")
}

func parseWallClock(input string) (wallClock, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return wallClock{}, false
	}
	for _, layout := range []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02",
		"20060102150405",
		"20060102",
	} {
		parsed, err := time.Parse(layout, trimmed)
		if err != nil {
			continue
		}
		return wallClock{
			year:   parsed.Year(),
			month:  parsed.Month(),
			day:    parsed.Day(),
			hour:   parsed.Hour(),
			minute: parsed.Minute(),
			second: parsed.Second(),
		}, true
	}
	return wallClock{}, false
}

func combineDateAndTime(baseDate string, hhmm string) (wallClock, bool) {
	date := strings.TrimSpace(baseDate)
	timeText := strings.TrimSpace(hhmm)
	if date == "" || timeText == "" {
		return wallClock{}, false
	}

	var year int
	var month int
	var day int
	var err error
	switch {
	case len(date) == 10 && date[4] == '-' && date[7] == '-':
		year, err = atoi(date[0:4])
		if err != nil {
			return wallClock{}, false
		}
		month, err = atoi(date[5:7])
		if err != nil {
			return wallClock{}, false
		}
		day, err = atoi(date[8:10])
		if err != nil {
			return wallClock{}, false
		}
	case len(date) == 8:
		year, err = atoi(date[0:4])
		if err != nil {
			return wallClock{}, false
		}
		month, err = atoi(date[4:6])
		if err != nil {
			return wallClock{}, false
		}
		day, err = atoi(date[6:8])
		if err != nil {
			return wallClock{}, false
		}
	default:
		return wallClock{}, false
	}

	parts := strings.Split(timeText, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return wallClock{}, false
	}
	hour, err := atoi(parts[0])
	if err != nil {
		return wallClock{}, false
	}
	minute, err := atoi(parts[1])
	if err != nil {
		return wallClock{}, false
	}
	second := 0
	if len(parts) == 3 {
		second, err = atoi(parts[2])
		if err != nil {
			return wallClock{}, false
		}
	}
	return validateWallClock(wallClock{
		year:   year,
		month:  time.Month(month),
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	})
}

func wallTimeToUnixMilli(wall wallClock, tz MarketTz) (int64, bool) {
	wall, ok := validateWallClock(wall)
	if !ok {
		return 0, false
	}
	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return 0, false
	}
	return time.Date(wall.year, wall.month, wall.day, wall.hour, wall.minute, wall.second, 0, location).UnixMilli(), true
}

func validateWallClock(wall wallClock) (wallClock, bool) {
	if wall.year <= 0 || wall.month < time.January || wall.month > time.December {
		return wallClock{}, false
	}
	if wall.hour < 0 || wall.hour > 23 || wall.minute < 0 || wall.minute > 59 || wall.second < 0 || wall.second > 59 {
		return wallClock{}, false
	}
	utc := time.Date(wall.year, wall.month, wall.day, wall.hour, wall.minute, wall.second, 0, time.UTC)
	if utc.Year() != wall.year || utc.Month() != wall.month || utc.Day() != wall.day ||
		utc.Hour() != wall.hour || utc.Minute() != wall.minute || utc.Second() != wall.second {
		return wallClock{}, false
	}
	return wall, true
}

func atoi(value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("empty number")
	}
	return strconv.Atoi(value)
}
