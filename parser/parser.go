package parser

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// TencentQuoteItem 是解析后的腾讯行情赋值项。
type TencentQuoteItem struct {
	Key    string
	Fields []string
}

// DecodeGBK 将 GBK 字节解码为 UTF-8 字符串。
func DecodeGBK(data []byte) (string, error) {
	return simplifiedchinese.GBK.NewDecoder().String(string(data))
}

// ParseResponse 解析腾讯财经 v_xxx="~" 响应文本。
func ParseResponse(text string) []TencentQuoteItem {
	lines := strings.Split(text, ";")
	items := make([]TencentQuoteItem, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		key = strings.TrimPrefix(key, "v_")
		raw := strings.TrimSpace(line[eq+1:])
		if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
			raw = strings.TrimPrefix(raw, `"`)
			raw = strings.TrimSuffix(raw, `"`)
		}
		items = append(items, TencentQuoteItem{
			Key:    key,
			Fields: strings.Split(raw, "~"),
		})
	}
	return items
}

// SafeNumber 将字符串安全转换为数字，空值或非法值返回 0。
func SafeNumber(value string) float64 {
	number, ok := parseFloat(value, false)
	if !ok {
		return 0
	}
	return number
}

// SafeNumberOrNil 将字符串安全转换为数字，空值或非法值返回 nil。
func SafeNumberOrNil(value string) *float64 {
	number, ok := parseFloat(value, false)
	if !ok {
		return nil
	}
	return &number
}

// ToNumber 将字符串转换为数字，空值、"-" 或非法值返回 nil。
func ToNumber(value string) *float64 {
	number, ok := parseFloat(value, true)
	if !ok {
		return nil
	}
	return &number
}

// ToNumberSafe 将任意值转换为数字，nil 或非法值返回 nil。
func ToNumberSafe(value any) *float64 {
	if value == nil {
		return nil
	}
	return ToNumber(fmt.Sprint(value))
}

func parseFloat(value string, dashAsNil bool) (float64, bool) {
	value = strings.TrimSpace(value)
	if value == "" || dashAsNil && value == "-" {
		return 0, false
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return number, true
}
