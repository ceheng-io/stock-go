package stock

import (
	"net/url"
	"strings"

	"github.com/ceheng.io/stock-go/types"
)

type linkMarket string

const (
	linkMarketSH      linkMarket = "sh"
	linkMarketSZ      linkMarket = "sz"
	linkMarketHK      linkMarket = "hk"
	linkMarketUS      linkMarket = "us"
	linkMarketGlobal  linkMarket = "global"
	linkMarketUnknown linkMarket = "unknown"
)

type normalizedSearchTarget struct {
	code   string
	market linkMarket
}

var tencentMarketMap = map[string]linkMarket{
	"0":   linkMarketSZ,
	"1":   linkMarketSH,
	"100": linkMarketGlobal,
	"105": linkMarketUS,
	"116": linkMarketHK,
}

var xueqiuGlobalIndexMap = map[string]string{
	"IXIC":   ".IXIC",
	"NDX":    ".NDX",
	"NDX100": ".NDX",
}

var eastmoneyGlobalIndexMap = map[string]string{
	"IXIC":   "NDX",
	"NDX":    "NDX",
	"NDX100": "NDX",
}

// GenerateSearchExternalLinks 根据搜索结果生成东方财富和雪球外部链接。
func GenerateSearchExternalLinks(result types.SearchResult) []types.ExternalLink {
	target := normalizeSearchTarget(result)
	return []types.ExternalLink{
		{Name: "东方财富", URL: buildEastmoneyURL(target)},
		{Name: "雪球", URL: buildXueqiuURL(target)},
	}
}

func normalizeSearchTarget(result types.SearchResult) normalizedSearchTarget {
	rawCode := strings.TrimSpace(result.Code)
	rawMarket := strings.ToLower(strings.TrimSpace(result.Market))
	market := linkMarket(rawMarket)
	if mapped, ok := tencentMarketMap[result.Market]; ok {
		market = mapped
	}

	switch market {
	case linkMarketSH, linkMarketSZ:
		return normalizedSearchTarget{code: stripPrefixFold(rawCode, "sh", "sz"), market: market}
	case linkMarketHK:
		return normalizedSearchTarget{code: leftPad(stripPrefixFold(rawCode, "hk"), 5, "0"), market: market}
	case linkMarketUS:
		code := stripTencentNumericPrefix(stripPrefixFold(rawCode, "us"))
		return normalizedSearchTarget{code: strings.ToUpper(code), market: market}
	case linkMarketGlobal:
		code := stripTencentNumericPrefix(rawCode)
		return normalizedSearchTarget{code: strings.ToUpper(code), market: market}
	}

	lowerCode := strings.ToLower(rawCode)
	switch {
	case hasMarketCode(lowerCode, "sh", 6):
		return normalizedSearchTarget{code: rawCode[2:], market: linkMarketSH}
	case hasMarketCode(lowerCode, "sz", 6):
		return normalizedSearchTarget{code: rawCode[2:], market: linkMarketSZ}
	case hasMarketCode(lowerCode, "hk", 5):
		return normalizedSearchTarget{code: rawCode[2:], market: linkMarketHK}
	case strings.HasPrefix(lowerCode, "us") && len(rawCode) > 2:
		return normalizedSearchTarget{code: strings.ToUpper(rawCode[2:]), market: linkMarketUS}
	default:
		return normalizedSearchTarget{code: rawCode, market: linkMarketUnknown}
	}
}

func buildEastmoneyURL(target normalizedSearchTarget) string {
	switch target.market {
	case linkMarketSH, linkMarketSZ:
		if isCNIndex(target.code, target.market) {
			return "https://quote.eastmoney.com/zs" + target.code + ".html"
		}
		return "https://quote.eastmoney.com/" + string(target.market) + target.code + ".html"
	case linkMarketHK:
		return "https://quote.eastmoney.com/hk/" + target.code + ".html"
	case linkMarketUS:
		return "https://quote.eastmoney.com/us/" + target.code + ".html"
	case linkMarketGlobal:
		code := eastmoneyGlobalIndexMap[target.code]
		if code == "" {
			code = target.code
		}
		return "https://quote.eastmoney.com/gb/zs" + code + ".html"
	default:
		return "https://so.eastmoney.com/web/s?keyword=" + url.QueryEscape(target.code)
	}
}

func buildXueqiuURL(target normalizedSearchTarget) string {
	switch target.market {
	case linkMarketSH, linkMarketSZ:
		return "https://xueqiu.com/S/" + strings.ToUpper(string(target.market)) + target.code
	case linkMarketHK, linkMarketUS:
		return "https://xueqiu.com/S/" + target.code
	case linkMarketGlobal:
		code := xueqiuGlobalIndexMap[target.code]
		if code == "" {
			code = target.code
		}
		return "https://xueqiu.com/S/" + code
	default:
		return "https://xueqiu.com/k?q=" + url.QueryEscape(target.code)
	}
}

func isCNIndex(code string, market linkMarket) bool {
	return market == linkMarketSH && strings.HasPrefix(code, "000") ||
		market == linkMarketSZ && strings.HasPrefix(code, "399")
}

func stripPrefixFold(value string, prefixes ...string) string {
	lower := strings.ToLower(value)
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return value[len(prefix):]
		}
	}
	return value
}

func stripTencentNumericPrefix(value string) string {
	if index := strings.IndexByte(value, '.'); index == 3 {
		return value[index+1:]
	}
	return value
}

func hasMarketCode(value string, prefix string, digits int) bool {
	if len(value) != len(prefix)+digits || !strings.HasPrefix(value, prefix) {
		return false
	}
	for _, r := range value[len(prefix):] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func leftPad(value string, length int, pad string) string {
	if len(value) >= length || pad == "" {
		return value
	}
	return strings.Repeat(pad, length-len(value)) + value
}
