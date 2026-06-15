package constants

import "testing"

func TestCoreURLConstantsMatchTypeScriptSDK(t *testing.T) {
	tests := map[string]string{
		"TencentBaseURL":         TencentBaseURL,
		"TencentMinuteURL":       TencentMinuteURL,
		"AShareListURL":          AShareListURL,
		"USListURL":              USListURL,
		"HKListURL":              HKListURL,
		"FundListURL":            FundListURL,
		"TradeCalendarURL":       TradeCalendarURL,
		"CodeListURL":            CodeListURL,
		"EMKlineURL":             EMKlineURL,
		"EMTrendsURL":            EMTrendsURL,
		"EMHKKlineURL":           EMHKKlineURL,
		"EMUSKlineURL":           EMUSKlineURL,
		"EMHKTrendsURL":          EMHKTrendsURL,
		"EMUSTrendsURL":          EMUSTrendsURL,
		"EMBoardListURL":         EMBoardListURL,
		"EMBoardSpotURL":         EMBoardSpotURL,
		"EMBoardConsURL":         EMBoardConsURL,
		"EMBoardKlineURL":        EMBoardKlineURL,
		"EMBoardTrendsURL":       EMBoardTrendsURL,
		"EMConceptListURL":       EMConceptListURL,
		"EMConceptSpotURL":       EMConceptSpotURL,
		"EMConceptConsURL":       EMConceptConsURL,
		"EMConceptKlineURL":      EMConceptKlineURL,
		"EMConceptTrendsURL":     EMConceptTrendsURL,
		"EMDatacenterURL":        EMDatacenterURL,
		"EMFFlowURL":             EMFFlowURL,
		"EMClistURL":             EMClistURL,
		"EMNorthboundMinuteURL":  EMNorthboundMinuteURL,
		"EMTopicBaseURL":         EMTopicBaseURL,
		"EMFuturesKlineURL":      EMFuturesKlineURL,
		"EMFuturesGlobalSpotURL": EMFuturesGlobalSpotURL,
		"SinaOptionAPIURL":       SinaOptionAPIURL,
		"SinaOptionDaylineURL":   SinaOptionDaylineURL,
		"SinaSSEOptionListURL":   SinaSSEOptionListURL,
		"SinaSSEOptionExpireURL": SinaSSEOptionExpireURL,
		"SinaSSEOptionMinuteURL": SinaSSEOptionMinuteURL,
		"SinaSSEOptionDailyURL":  SinaSSEOptionDailyURL,
		"SinaSSEOption5DayURL":   SinaSSEOption5DayURL,
		"EMOptionCFFEXURL":       EMOptionCFFEXURL,
		"EMOptionLHBURL":         EMOptionLHBURL,
		"THSLimitUpPoolURL":      THSLimitUpPoolURL,
	}

	want := map[string]string{
		"TencentBaseURL":         "https://qt.gtimg.cn",
		"TencentMinuteURL":       "https://web.ifzq.gtimg.cn/appstock/app/minute/query",
		"AShareListURL":          "https://assets.linkdiary.cn/shares/zh_a_list.json",
		"USListURL":              "https://assets.linkdiary.cn/shares/us_list.json",
		"HKListURL":              "https://assets.linkdiary.cn/shares/hk_list.json",
		"FundListURL":            "https://assets.linkdiary.cn/shares/fund_list",
		"TradeCalendarURL":       "https://assets.linkdiary.cn/shares/trade-data-list.txt",
		"CodeListURL":            "https://assets.linkdiary.cn/shares/zh_a_list.json",
		"EMKlineURL":             "https://push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMTrendsURL":            "https://push2his.eastmoney.com/api/qt/stock/trends2/get",
		"EMHKKlineURL":           "https://33.push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMUSKlineURL":           "https://63.push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMHKTrendsURL":          "https://33.push2his.eastmoney.com/api/qt/stock/trends2/get",
		"EMUSTrendsURL":          "https://63.push2his.eastmoney.com/api/qt/stock/trends2/get",
		"EMBoardListURL":         "https://17.push2.eastmoney.com/api/qt/clist/get",
		"EMBoardSpotURL":         "https://91.push2.eastmoney.com/api/qt/stock/get",
		"EMBoardConsURL":         "https://29.push2.eastmoney.com/api/qt/clist/get",
		"EMBoardKlineURL":        "https://7.push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMBoardTrendsURL":       "https://push2his.eastmoney.com/api/qt/stock/trends2/get",
		"EMConceptListURL":       "https://79.push2.eastmoney.com/api/qt/clist/get",
		"EMConceptSpotURL":       "https://91.push2.eastmoney.com/api/qt/stock/get",
		"EMConceptConsURL":       "https://29.push2.eastmoney.com/api/qt/clist/get",
		"EMConceptKlineURL":      "https://91.push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMConceptTrendsURL":     "https://push2his.eastmoney.com/api/qt/stock/trends2/get",
		"EMDatacenterURL":        "https://datacenter-web.eastmoney.com/api/data/v1/get",
		"EMFFlowURL":             "https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get",
		"EMClistURL":             "https://push2.eastmoney.com/api/qt/clist/get",
		"EMNorthboundMinuteURL":  "https://push2.eastmoney.com/api/qt/kamtbs.rtmin/get",
		"EMTopicBaseURL":         "https://push2ex.eastmoney.com",
		"EMFuturesKlineURL":      "https://push2his.eastmoney.com/api/qt/stock/kline/get",
		"EMFuturesGlobalSpotURL": "https://futsseapi.eastmoney.com/list/COMEX,NYMEX,COBOT,SGX,NYBOT,LME,MDEX,TOCOM,IPE",
		"SinaOptionAPIURL":       "https://stock.finance.sina.com.cn/futures/api/openapi.php/OptionService.getOptionData",
		"SinaOptionDaylineURL":   "https://stock.finance.sina.com.cn/futures/api/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline",
		"SinaSSEOptionListURL":   "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionService.getStockName",
		"SinaSSEOptionExpireURL": "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionService.getRemainderDay",
		"SinaSSEOptionMinuteURL": "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionDaylineService.getOptionMinline",
		"SinaSSEOptionDailyURL":  "https://stock.finance.sina.com.cn/futures/api/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo",
		"SinaSSEOption5DayURL":   "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionDaylineService.getFiveDayLine",
		"EMOptionCFFEXURL":       "https://futsseapi.eastmoney.com/list/option/221",
		"EMOptionLHBURL":         "https://datacenter-web.eastmoney.com/api/data/get",
		"THSLimitUpPoolURL":      "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool",
	}

	for name, got := range tests {
		if got != want[name] {
			t.Fatalf("%s = %q, want %q", name, got, want[name])
		}
	}
}

func TestCoreTokenAndDefaultConstantsMatchTypeScriptSDK(t *testing.T) {
	if EMPushToken != "7eea3edcaed734bea9cbfc24409ed989" {
		t.Fatalf("EMPushToken = %q", EMPushToken)
	}
	if EMDataToken != "b2884a393a59ad64002292a3e90d46a5" {
		t.Fatalf("EMDataToken = %q", EMDataToken)
	}
	if EMOptionLHBToken != EMDataToken {
		t.Fatalf("EMOptionLHBToken = %q, want EMDataToken", EMOptionLHBToken)
	}
	if EMFuturesGlobalSpotToken != "58b2fa8f54638b60b87d69b31969089c" {
		t.Fatalf("EMFuturesGlobalSpotToken = %q", EMFuturesGlobalSpotToken)
	}
	if DefaultTimeoutMS != 30000 || DefaultBatchSize != 500 || MaxBatchSize != 500 || DefaultConcurrency != 7 {
		t.Fatalf("default request values mismatch")
	}
	if DefaultMaxRetries != 3 || DefaultBaseDelayMS != 1000 || DefaultMaxDelayMS != 30000 || DefaultBackoffMultiplier != 2 {
		t.Fatalf("default retry values mismatch")
	}

	statusCodes := DefaultRetryableStatusCodes()
	wantCodes := []int{408, 429, 500, 502, 503, 504}
	if len(statusCodes) != len(wantCodes) {
		t.Fatalf("len(DefaultRetryableStatusCodes) = %d, want %d", len(statusCodes), len(wantCodes))
	}
	for i := range wantCodes {
		if statusCodes[i] != wantCodes[i] {
			t.Fatalf("DefaultRetryableStatusCodes[%d] = %d, want %d", i, statusCodes[i], wantCodes[i])
		}
	}
	statusCodes[0] = 0
	if got := DefaultRetryableStatusCodes()[0]; got != 408 {
		t.Fatalf("DefaultRetryableStatusCodes returned mutable backing slice; got %d", got)
	}
}

func TestCoreMapConstantsMatchTypeScriptSDK(t *testing.T) {
	exchanges := FuturesExchangeMap()
	if exchanges["SHFE"] != 113 || exchanges["CFFEX"] != 220 || exchanges["GFEX"] != 225 {
		t.Fatalf("FuturesExchangeMap = %#v", exchanges)
	}
	exchanges["SHFE"] = 0
	if got := FuturesExchangeMap()["SHFE"]; got != 113 {
		t.Fatalf("FuturesExchangeMap returned mutable backing map; got %d", got)
	}

	varieties := FuturesVarietyExchange()
	if varieties["rb"] != "SHFE" || varieties["IF"] != "CFFEX" || varieties["SR"] != "CZCE" || varieties["si"] != "GFEX" {
		t.Fatalf("FuturesVarietyExchange = %#v", varieties)
	}
	varieties["rb"] = "BAD"
	if got := FuturesVarietyExchange()["rb"]; got != "SHFE" {
		t.Fatalf("FuturesVarietyExchange returned mutable backing map; got %q", got)
	}

	global := GlobalFuturesMarket()
	if global["GC"] != 101 || global["NG"] != 102 || global["NQ"] != 103 || global["SB"] != 108 || global["LCPT"] != 109 {
		t.Fatalf("GlobalFuturesMarket = %#v", global)
	}

	products := CFFEXOptionProductMap()
	if products["ho"] != "上证50" || products["io"] != "沪深300" || products["mo"] != "中证1000" {
		t.Fatalf("CFFEXOptionProductMap = %#v", products)
	}

	commodity := CommodityOptionMap()
	if commodity["au"] != (CommodityOption{Product: "au_o", Exchange: "shfe"}) {
		t.Fatalf("CommodityOptionMap[au] = %#v", commodity["au"])
	}
	if commodity["SR"] != (CommodityOption{Product: "SR_o", Exchange: "czce"}) {
		t.Fatalf("CommodityOptionMap[SR] = %#v", commodity["SR"])
	}
	commodity["au"] = CommodityOption{Product: "bad", Exchange: "bad"}
	if got := CommodityOptionMap()["au"]; got != (CommodityOption{Product: "au_o", Exchange: "shfe"}) {
		t.Fatalf("CommodityOptionMap returned mutable backing map; got %#v", got)
	}
}

func TestCoreTSSnakeCaseAliasesMatchTypeScriptSDK(t *testing.T) {
	if TENCENT_BASE_URL != TencentBaseURL || EM_PUSH_TOKEN != EMPushToken || SINA_OPTION_API_URL != SinaOptionAPIURL {
		t.Fatalf("snake-case URL/token aliases are not wired to canonical constants")
	}
	if DEFAULT_TIMEOUT != DefaultTimeoutMS || DEFAULT_BATCH_SIZE != DefaultBatchSize || MAX_BATCH_SIZE != MaxBatchSize {
		t.Fatalf("snake-case default aliases mismatch")
	}

	if FUTURES_EXCHANGE_MAP["SHFE"] != 113 || FUTURES_VARIETY_EXCHANGE["rb"] != "SHFE" || GLOBAL_FUTURES_MARKET["GC"] != 101 {
		t.Fatalf("snake-case futures maps mismatch")
	}
	if CFFEX_OPTION_PRODUCT_MAP["io"] != "沪深300" || COMMODITY_OPTION_MAP["au"] != (CommodityOption{Product: "au_o", Exchange: "shfe"}) {
		t.Fatalf("snake-case option maps mismatch")
	}
	if DEFAULT_RETRYABLE_STATUS_CODES[0] != 408 || DEFAULT_RETRYABLE_STATUS_CODES[len(DEFAULT_RETRYABLE_STATUS_CODES)-1] != 504 {
		t.Fatalf("DEFAULT_RETRYABLE_STATUS_CODES = %#v", DEFAULT_RETRYABLE_STATUS_CODES)
	}
}
