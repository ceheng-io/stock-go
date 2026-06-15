package sina

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

type commodityOptionMapping struct {
	Product  string
	Exchange string
}

var commodityOptionMap = map[string]commodityOptionMapping{
	"au": {Product: "au_o", Exchange: "shfe"},
	"ag": {Product: "ag_o", Exchange: "shfe"},
	"cu": {Product: "cu_o", Exchange: "shfe"},
	"al": {Product: "al_o", Exchange: "shfe"},
	"zn": {Product: "zn_o", Exchange: "shfe"},
	"ru": {Product: "ru_o", Exchange: "shfe"},
	"sc": {Product: "sc_o", Exchange: "ine"},
	"m":  {Product: "m_o", Exchange: "dce"},
	"c":  {Product: "c_o", Exchange: "dce"},
	"i":  {Product: "i_o", Exchange: "dce"},
	"p":  {Product: "p_o", Exchange: "dce"},
	"pp": {Product: "pp_o", Exchange: "dce"},
	"l":  {Product: "l_o", Exchange: "dce"},
	"v":  {Product: "v_o", Exchange: "dce"},
	"pg": {Product: "pg_o", Exchange: "dce"},
	"y":  {Product: "y_o", Exchange: "dce"},
	"a":  {Product: "a_o", Exchange: "dce"},
	"b":  {Product: "b_o", Exchange: "dce"},
	"eg": {Product: "eg_o", Exchange: "dce"},
	"eb": {Product: "eb_o", Exchange: "dce"},
	"SR": {Product: "SR_o", Exchange: "czce"},
	"CF": {Product: "CF_o", Exchange: "czce"},
	"TA": {Product: "TA_o", Exchange: "czce"},
	"MA": {Product: "MA_o", Exchange: "czce"},
	"RM": {Product: "RM_o", Exchange: "czce"},
	"OI": {Product: "OI_o", Exchange: "czce"},
	"PK": {Product: "PK_o", Exchange: "czce"},
	"PF": {Product: "PF_o", Exchange: "czce"},
	"SA": {Product: "SA_o", Exchange: "czce"},
	"UR": {Product: "UR_o", Exchange: "czce"},
}

// GetCommodityOptionSpot 获取新浪商品期权 T 型报价。
func GetCommodityOptionSpot(ctx context.Context, client JSONPClient, endpoint string, variety string, contract string) (types.OptionTQuoteResult, error) {
	mapping, ok := commodityOptionMap[variety]
	if !ok {
		message := fmt.Sprintf("unknown commodity option variety %q: available %s", variety, strings.Join(commodityOptionVarieties(), ", "))
		return types.OptionTQuoteResult{}, core.NewCodedError("INVALID_ARGUMENT", message, nil)
	}

	params := url.Values{}
	params.Set("type", "futures")
	params.Set("product", mapping.Product)
	params.Set("exchange", mapping.Exchange)
	params.Set("pinzhong", contract)

	var payload sinaOptionSpotResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return types.OptionTQuoteResult{}, err
	}
	return types.OptionTQuoteResult{
		Calls: parseOptionCallQuotes(payload.Result.Data.Up),
		Puts:  parseOptionPutQuotes(payload.Result.Data.Down),
	}, nil
}

// GetCommodityOptionKline 获取新浪商品期权合约日 K 线。
func GetCommodityOptionKline(ctx context.Context, client JSONPClient, endpoint string, symbol string) ([]types.OptionKline, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	return getSinaOptionKlinesJSONP(ctx, client, endpoint, params)
}

func commodityOptionVarieties() []string {
	keys := make([]string, 0, len(commodityOptionMap))
	for key := range commodityOptionMap {
		keys = append(keys, key)
	}
	return keys
}
