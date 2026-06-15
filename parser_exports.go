package stock

import "github.com/ceheng.io/stock-go/parser"

type TencentQuoteItem = parser.TencentQuoteItem

// DecodeGBK 将 GBK 字节解码为 UTF-8 字符串。
func DecodeGBK(data []byte) (string, error) {
	return parser.DecodeGBK(data)
}

// ParseResponse 解析腾讯财经 v_xxx="~" 响应文本。
func ParseResponse(text string) []TencentQuoteItem {
	return parser.ParseResponse(text)
}

// SafeNumber 将字符串安全转换为数字，空值或非法值返回 0。
func SafeNumber(value string) float64 {
	return parser.SafeNumber(value)
}

// SafeNumberOrNil 将字符串安全转换为数字，空值或非法值返回 nil。
func SafeNumberOrNil(value string) *float64 {
	return parser.SafeNumberOrNil(value)
}

// SafeNumberOrNull 将字符串安全转换为数字，空值或非法值返回 nil。
func SafeNumberOrNull(value string) *float64 {
	return SafeNumberOrNil(value)
}

// ToNumber 将字符串转换为数字，空值、"-" 或非法值返回 nil。
func ToNumber(value string) *float64 {
	return parser.ToNumber(value)
}

// ToNumberSafe 将任意值转换为数字，nil 或非法值返回 nil。
func ToNumberSafe(value any) *float64 {
	return parser.ToNumberSafe(value)
}
