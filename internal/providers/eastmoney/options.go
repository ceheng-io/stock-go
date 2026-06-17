package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/ceheng-io/stock-go/types"
)

// CFFEXOptionQuotesOptions 配置中金所期权实时行情查询。
type CFFEXOptionQuotesOptions struct {
	PageSize int
}

// OptionClient 是期权 provider 所需的最小请求客户端接口。
type OptionClient interface {
	GetJSON(context.Context, string, any) error
}

type cffexOptionQuotesResponse struct {
	List json.RawMessage `json:"list"`
}

type optionLHBResponse struct {
	Result struct {
		Data json.RawMessage `json:"data"`
	} `json:"result"`
}

// GetCFFEXOptionQuotes 获取中金所全部期权实时行情列表。
func GetCFFEXOptionQuotes(ctx context.Context, client OptionClient, endpoint string, options CFFEXOptionQuotesOptions) ([]types.CFFEXOptionQuote, error) {
	pageSize := options.PageSize
	if pageSize <= 0 {
		pageSize = 20000
	}
	params := url.Values{}
	params.Set("orderBy", "zdf")
	params.Set("sort", "desc")
	params.Set("pageSize", strconv.Itoa(pageSize))
	params.Set("pageIndex", "0")
	params.Set("token", emFuturesGlobalSpotToken)
	params.Set("field", "dm,sc,name,p,zsjd,zde,zdf,f152,vol,cje,ccl,xqj,syr,rz,zjsj,o")

	var payload cffexOptionQuotesResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	items, err := decodeBoardDynamicArray(payload.List)
	if err != nil {
		return nil, err
	}
	rows := make([]types.CFFEXOptionQuote, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseCFFEXOptionQuote(item))
	}
	return rows, nil
}

// GetOptionLHB 获取期权龙虎榜数据。
func GetOptionLHB(ctx context.Context, client OptionClient, endpoint string, symbol string, date string) ([]types.OptionLHBItem, error) {
	params := url.Values{}
	params.Set("type", "RPT_IF_BILLBOARD_TD")
	params.Set("sty", "ALL")
	params.Set("p", "1")
	params.Set("ps", "200")
	params.Set("source", "IFBILLBOARD")
	params.Set("client", "WEB")
	params.Set("ut", emDataToken)
	params.Set("filter", `(SECURITY_CODE="`+symbol+`")(TRADE_DATE='`+date+`')`)

	var payload optionLHBResponse
	if err := client.GetJSON(ctx, endpoint+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	items, err := decodeBoardDynamicArray(payload.Result.Data)
	if err != nil {
		return nil, err
	}
	rows := make([]types.OptionLHBItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, parseOptionLHBItem(item))
	}
	return rows, nil
}

func parseCFFEXOptionQuote(item boardDynamicItem) types.CFFEXOptionQuote {
	return types.CFFEXOptionQuote{
		Code:          stringValue(item["dm"]),
		Name:          stringValue(item["name"]),
		Price:         toNumberFromAny(item["p"]),
		Change:        toNumberFromAny(item["zde"]),
		ChangePercent: toNumberFromAny(item["zdf"]),
		Volume:        toNumberFromAny(item["vol"]),
		Amount:        toNumberFromAny(item["cje"]),
		OpenInterest:  toNumberFromAny(item["ccl"]),
		StrikePrice:   toNumberFromAny(item["xqj"]),
		RemainDays:    toNumberFromAny(item["syr"]),
		DailyChange:   toNumberFromAny(item["rz"]),
		PrevSettle:    toNumberFromAny(item["zjsj"]),
		Open:          toNumberFromAny(item["o"]),
	}
}

func parseOptionLHBItem(item boardDynamicItem) types.OptionLHBItem {
	return types.OptionLHBItem{
		TradeType:          stringValue(item["TRADE_TYPE"]),
		Date:               parseDatacenterDate(item["TRADE_DATE"]),
		Symbol:             stringValue(item["SECURITY_CODE"]),
		TargetName:         stringValue(item["TARGET_NAME"]),
		Rank:               int(numberValue(toNumberFromAny(item["MEMBER_RANK"]))),
		MemberName:         stringValue(item["MEMBER_NAME_ABBR"]),
		SellVolume:         toNumberFromAny(item["SELL_VOLUME"]),
		SellVolumeChange:   toNumberFromAny(item["SELL_VOLUME_CHANGE"]),
		NetSellVolume:      toNumberFromAny(item["NET_SELL_VOLUME"]),
		SellVolumeRatio:    toNumberFromAny(item["SELL_VOLUME_RATIO"]),
		BuyVolume:          toNumberFromAny(item["BUY_VOLUME"]),
		BuyVolumeChange:    toNumberFromAny(item["BUY_VOLUME_CHANGE"]),
		NetBuyVolume:       toNumberFromAny(item["NET_BUY_VOLUME"]),
		BuyVolumeRatio:     toNumberFromAny(item["BUY_VOLUME_RATIO"]),
		SellPosition:       toNumberFromAny(item["SELL_POSITION"]),
		SellPositionChange: toNumberFromAny(item["SELL_POSITION_CHANGE"]),
		NetSellPosition:    toNumberFromAny(item["NET_SELL_POSITION"]),
		SellPositionRatio:  toNumberFromAny(item["SELL_POSITION_RATIO"]),
		BuyPosition:        toNumberFromAny(item["BUY_POSITION"]),
		BuyPositionChange:  toNumberFromAny(item["BUY_POSITION_CHANGE"]),
		NetBuyPosition:     toNumberFromAny(item["NET_BUY_POSITION"]),
		BuyPositionRatio:   toNumberFromAny(item["BUY_POSITION_RATIO"]),
	}
}
