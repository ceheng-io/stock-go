package core

import "github.com/ceheng-io/stock-go/parser"

// TencentQuoteItem is a parsed Tencent quote assignment.
type TencentQuoteItem = parser.TencentQuoteItem

// ParseTencentQuoteResponse parses v_xxx="~"-delimited Tencent quote text.
func ParseTencentQuoteResponse(text string) []TencentQuoteItem {
	return parser.ParseResponse(text)
}
