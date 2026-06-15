package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/types"
)

const emTopicPushToken = "7eea3edcaed734bea9cbfc24409ed989"

// ZTPoolType is an Eastmoney limit-up pool category.
type ZTPoolType = types.ZTPoolType

const (
	ZTPoolZT        ZTPoolType = types.ZTPoolZT
	ZTPoolYesterday ZTPoolType = types.ZTPoolYesterday
	ZTPoolStrong    ZTPoolType = types.ZTPoolStrong
	ZTPoolSubNew    ZTPoolType = types.ZTPoolSubNew
	ZTPoolBroken    ZTPoolType = types.ZTPoolBroken
	ZTPoolDT        ZTPoolType = types.ZTPoolDT
)

// StockChangeType is an Eastmoney intraday stock change category.
type StockChangeType = types.StockChangeType

const (
	StockChangeRocketLaunch   StockChangeType = types.StockChangeRocketLaunch
	StockChangeQuickRebound   StockChangeType = types.StockChangeQuickRebound
	StockChangeLargeBuy       StockChangeType = types.StockChangeLargeBuy
	StockChangeLimitUpSeal    StockChangeType = types.StockChangeLimitUpSeal
	StockChangeLimitDownOpen  StockChangeType = types.StockChangeLimitDownOpen
	StockChangeBigBuyOrder    StockChangeType = types.StockChangeBigBuyOrder
	StockChangeAuctionUp      StockChangeType = types.StockChangeAuctionUp
	StockChangeHighOpen5D     StockChangeType = types.StockChangeHighOpen5D
	StockChangeGapUp          StockChangeType = types.StockChangeGapUp
	StockChangeHigh60D        StockChangeType = types.StockChangeHigh60D
	StockChangeSurge60D       StockChangeType = types.StockChangeSurge60D
	StockChangeAccelerateDown StockChangeType = types.StockChangeAccelerateDown
	StockChangeHighDive       StockChangeType = types.StockChangeHighDive
	StockChangeLargeSell      StockChangeType = types.StockChangeLargeSell
	StockChangeLimitDownSeal  StockChangeType = types.StockChangeLimitDownSeal
	StockChangeLimitUpOpen    StockChangeType = types.StockChangeLimitUpOpen
	StockChangeBigSellOrder   StockChangeType = types.StockChangeBigSellOrder
	StockChangeAuctionDown    StockChangeType = types.StockChangeAuctionDown
	StockChangeLowOpen5D      StockChangeType = types.StockChangeLowOpen5D
	StockChangeGapDown        StockChangeType = types.StockChangeGapDown
	StockChangeLow60D         StockChangeType = types.StockChangeLow60D
	StockChangeDrop60D        StockChangeType = types.StockChangeDrop60D
)

// TopicDataClient is the minimal client interface required by topic data providers.
type TopicDataClient interface {
	GetJSON(context.Context, string, any) error
}

type ztPoolConfig struct {
	path string
	sort string
}

var ztPoolConfigs = map[ZTPoolType]ztPoolConfig{
	ZTPoolZT:        {path: "/getTopicZTPool", sort: "fbt:asc"},
	ZTPoolYesterday: {path: "/getYesterdayZTPool", sort: "zs:desc"},
	ZTPoolStrong:    {path: "/getTopicQSPool", sort: "zdp:desc"},
	ZTPoolSubNew:    {path: "/getTopicCXPool", sort: "ods:asc"},
	ZTPoolBroken:    {path: "/getTopicZBPool", sort: "fbt:asc"},
	ZTPoolDT:        {path: "/getTopicDTPool", sort: "fund:asc"},
}

var stockChangeTypeToCode = map[StockChangeType]string{
	StockChangeRocketLaunch:   "8201",
	StockChangeQuickRebound:   "8202",
	StockChangeLargeBuy:       "8193",
	StockChangeLimitUpSeal:    "4",
	StockChangeLimitDownOpen:  "32",
	StockChangeBigBuyOrder:    "64",
	StockChangeAuctionUp:      "8207",
	StockChangeHighOpen5D:     "8209",
	StockChangeGapUp:          "8211",
	StockChangeHigh60D:        "8213",
	StockChangeSurge60D:       "8215",
	StockChangeAccelerateDown: "8204",
	StockChangeHighDive:       "8203",
	StockChangeLargeSell:      "8194",
	StockChangeLimitDownSeal:  "8",
	StockChangeLimitUpOpen:    "16",
	StockChangeBigSellOrder:   "128",
	StockChangeAuctionDown:    "8208",
	StockChangeLowOpen5D:      "8210",
	StockChangeGapDown:        "8212",
	StockChangeLow60D:         "8214",
	StockChangeDrop60D:        "8216",
}

var stockChangeCodeToLabel = map[string]string{
	"8201": "火箭发射",
	"8202": "快速反弹",
	"8193": "大笔买入",
	"4":    "封涨停板",
	"32":   "打开跌停板",
	"64":   "有大买盘",
	"8207": "竞价上涨",
	"8209": "高开5日线",
	"8211": "向上缺口",
	"8213": "60日新高",
	"8215": "60日大幅上涨",
	"8204": "加速下跌",
	"8203": "高台跳水",
	"8194": "大笔卖出",
	"8":    "封跌停板",
	"16":   "打开涨停板",
	"128":  "有大卖盘",
	"8208": "竞价下跌",
	"8210": "低开5日线",
	"8212": "向下缺口",
	"8214": "60日新低",
	"8216": "60日大幅下跌",
}

type ztPoolResponse struct {
	Data struct {
		Pool json.RawMessage `json:"pool"`
	} `json:"data"`
}

type stockChangesResponse struct {
	Data struct {
		AllStock json.RawMessage `json:"allstock"`
	} `json:"data"`
}

type boardChangesResponse struct {
	Data struct {
		AllBK json.RawMessage `json:"allbk"`
	} `json:"data"`
}

// GetZTPool fetches limit-up pool rows.
func GetZTPool(ctx context.Context, client TopicDataClient, baseURL string, poolType ZTPoolType, date string) ([]types.ZTPoolItem, error) {
	if poolType == "" {
		poolType = ZTPoolZT
	}
	config, ok := ztPoolConfigs[poolType]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf("invalid zt pool type %q", poolType))
	}
	queryDate := normalizeTopicDate(date)
	if queryDate == "" {
		queryDate = beijingDateString(time.Now())
	}
	params := url.Values{}
	params.Set("ut", emTopicPushToken)
	params.Set("dpt", "wz.ztzt")
	params.Set("Pageindex", "0")
	params.Set("pagesize", "10000")
	params.Set("sort", config.sort)
	params.Set("date", queryDate)

	var payload ztPoolResponse
	if err := client.GetJSON(ctx, topicURL(baseURL, config.path, params), &payload); err != nil {
		return nil, err
	}
	pool, err := decodeBoardDynamicArray(payload.Data.Pool)
	if err != nil {
		return nil, err
	}
	rows := make([]types.ZTPoolItem, 0, len(pool))
	for _, item := range pool {
		rows = append(rows, parseZTPoolItem(item))
	}
	return rows, nil
}

// GetStockChanges fetches intraday stock change rows.
func GetStockChanges(ctx context.Context, client TopicDataClient, baseURL string, changeType StockChangeType) ([]types.StockChangeItem, error) {
	if changeType == "" {
		changeType = StockChangeLargeBuy
	}
	code, ok := stockChangeTypeToCode[changeType]
	if !ok {
		return nil, invalidArgumentError(fmt.Sprintf("invalid stock change type %q", changeType))
	}
	params := url.Values{}
	params.Set("type", code)
	params.Set("pageindex", "0")
	params.Set("pagesize", "5000")
	params.Set("ut", emTopicPushToken)
	params.Set("dpt", "wzchanges")

	var payload stockChangesResponse
	if err := client.GetJSON(ctx, topicURL(baseURL, "/getAllStockChanges", params), &payload); err != nil {
		return nil, err
	}
	items, err := decodeBoardDynamicArray(payload.Data.AllStock)
	if err != nil {
		return nil, err
	}
	rows := make([]types.StockChangeItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseStockChangeItem(item, changeType, code))
	}
	return rows, nil
}

// GetBoardChanges fetches board change rows.
func GetBoardChanges(ctx context.Context, client TopicDataClient, baseURL string) ([]types.BoardChangeItem, error) {
	params := url.Values{}
	params.Set("ut", emTopicPushToken)
	params.Set("dpt", "wzchanges")
	params.Set("pageindex", "0")
	params.Set("pagesize", "5000")

	var payload boardChangesResponse
	if err := client.GetJSON(ctx, topicURL(baseURL, "/getAllBKChanges", params), &payload); err != nil {
		return nil, err
	}
	items, err := decodeBoardDynamicArray(payload.Data.AllBK)
	if err != nil {
		return nil, err
	}
	rows := make([]types.BoardChangeItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseBoardChangeItem(item))
	}
	return rows, nil
}

func topicURL(baseURL string, path string, params url.Values) string {
	return strings.TrimRight(baseURL, "/") + path + "?" + params.Encode()
}

func normalizeTopicDate(date string) string {
	if date == "" {
		return ""
	}
	return strings.ReplaceAll(date, "-", "")
}

func beijingDateString(now time.Time) string {
	beijing := now.UTC().Add(8 * time.Hour)
	return beijing.Format("20060102")
}

func parseZTPoolItem(item boardDynamicItem) types.ZTPoolItem {
	return types.ZTPoolItem{
		Code:                 stringFromKeys(item, "c", "m"),
		Name:                 stringValue(item["n"]),
		Price:                scaledNumber(item["p"], 1000),
		ChangePercent:        toNumberFromAny(item["zdp"]),
		LimitPrice:           scaledNumber(item["tp"], 1000),
		Amount:               numberFromKeys(item, "amount", "zb"),
		FloatMarketValue:     toNumberFromAny(item["ltsz"]),
		TotalMarketValue:     toNumberFromAny(item["tshare"]),
		TurnoverRate:         toNumberFromAny(item["hs"]),
		ContinuousBoardCount: toNumberFromAny(item["lbc"]),
		FirstBoardTime:       nullableTopicTimeFromKey(item, "fbt"),
		LastBoardTime:        nullableTopicTimeFromKey(item, "lbt"),
		BoardAmount:          toNumberFromAny(item["fund"]),
		SealAmount:           toNumberFromAny(item["fund"]),
		FailedCount:          toNumberFromAny(item["zbc"]),
		Industry:             stringValue(item["hybk"]),
		ZTStatistics:         ztStatisticsString(item["zttj"]),
		Amplitude:            toNumberFromAny(item["zf"]),
		Speed:                toNumberFromAny(item["zs"]),
	}
}

func parseStockChangeItem(item boardDynamicItem, changeType StockChangeType, defaultCode string) types.StockChangeItem {
	changeCode := stringValue(item["t"])
	if changeCode == "" {
		changeCode = defaultCode
	}
	return types.StockChangeItem{
		Time:            topicTimeString(item["tm"]),
		Code:            stringValue(item["c"]),
		Name:            stringValue(item["n"]),
		ChangeType:      changeType,
		ChangeTypeLabel: stockChangeCodeToLabel[changeCode],
		Info:            stringValue(item["i"]),
	}
}

func parseBoardChangeItem(item boardDynamicItem) types.BoardChangeItem {
	ms, _ := item["ms"].(map[string]any)
	direction := ""
	if numberValue(toNumberFromAny(ms["m"])) == 0 {
		direction = "大笔买入"
	} else if numberValue(toNumberFromAny(ms["m"])) == 1 {
		direction = "大笔卖出"
	}
	return types.BoardChangeItem{
		Name:                   stringValue(item["bkn"]),
		ChangePercent:          toNumberFromAny(item["bkz"]),
		MainNetInflow:          toNumberFromAny(item["bkj"]),
		TotalChangeCount:       toNumberFromAny(item["bkc"]),
		TopStockCode:           stringValue(ms["c"]),
		TopStockName:           stringValue(ms["n"]),
		TopStockDirection:      direction,
		ChangeTypeDistribution: topicDistribution(firstValue(item, "bkdf", "bkdfdis")),
	}
}

func scaledNumber(value any, divisor float64) *float64 {
	number := toNumberFromAny(value)
	if number == nil {
		return nil
	}
	scaled := *number / divisor
	return &scaled
}

func topicTimeString(value any) string {
	if value == nil {
		return ""
	}
	text := stringValue(value)
	if dot := strings.IndexByte(text, '.'); dot >= 0 {
		text = text[:dot]
	}
	if len(text) < 6 {
		text = strings.Repeat("0", 6-len(text)) + text
	}
	return text[:2] + ":" + text[2:4] + ":" + text[4:6]
}

func nullableTopicTime(value any) *string {
	text := topicTimeString(value)
	if text == "" {
		return nil
	}
	return &text
}

func nullableTopicTimeFromKey(item boardDynamicItem, key string) *string {
	value, ok := item[key]
	if !ok || value == nil {
		return nil
	}
	text := topicTimeString(value)
	return &text
}

func ztStatisticsString(value any) string {
	obj, ok := value.(map[string]any)
	if !ok {
		return ""
	}
	days := compactNumberString(obj["days"])
	count := compactNumberString(obj["ct"])
	return days + "/" + count
}

func compactNumberString(value any) string {
	number := toNumberFromAny(value)
	if number != nil {
		return fmt.Sprintf("%g", *number)
	}
	return stringValue(value)
}

func topicDistribution(value any) map[string]float64 {
	distribution := map[string]float64{}
	obj, ok := value.(map[string]any)
	if !ok {
		return distribution
	}
	for key, raw := range obj {
		if number := toNumberFromAny(raw); number != nil {
			distribution[key] = *number
		} else {
			distribution[key] = 0
		}
	}
	return distribution
}
