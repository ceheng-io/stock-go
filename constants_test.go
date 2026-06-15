package stock

import "testing"

func TestRootReExportsConstants(t *testing.T) {
	if TencentBaseURL != "https://qt.gtimg.cn" {
		t.Fatalf("TencentBaseURL = %q", TencentBaseURL)
	}
	if EMKlineURL != "https://push2his.eastmoney.com/api/qt/stock/kline/get" {
		t.Fatalf("EMKlineURL = %q", EMKlineURL)
	}
	if EMPushToken != "7eea3edcaed734bea9cbfc24409ed989" {
		t.Fatalf("EMPushToken = %q", EMPushToken)
	}
	if DefaultBatchSize != 500 || DefaultConcurrency != 7 {
		t.Fatalf("default batch values mismatch")
	}

	statusCodes := DefaultRetryableStatusCodes()
	if len(statusCodes) != 6 || statusCodes[0] != 408 || statusCodes[5] != 504 {
		t.Fatalf("DefaultRetryableStatusCodes = %#v", statusCodes)
	}
	statusCodes[0] = 0
	if DefaultRetryableStatusCodes()[0] != 408 {
		t.Fatal("DefaultRetryableStatusCodes returned mutable backing slice")
	}

	exchanges := FuturesExchangeMap()
	if exchanges["CFFEX"] != 220 {
		t.Fatalf("FuturesExchangeMap = %#v", exchanges)
	}
	varieties := FuturesVarietyExchange()
	if varieties["rb"] != "SHFE" || varieties["IF"] != "CFFEX" {
		t.Fatalf("FuturesVarietyExchange = %#v", varieties)
	}
	global := GlobalFuturesMarket()
	if global["GC"] != 101 || global["NG"] != 102 {
		t.Fatalf("GlobalFuturesMarket = %#v", global)
	}
	products := CFFEXOptionProductMap()
	if products["io"] != "沪深300" {
		t.Fatalf("CFFEXOptionProductMap = %#v", products)
	}
	commodity := CommodityOptionMap()
	if commodity["au"] != (CommodityOption{Product: "au_o", Exchange: "shfe"}) {
		t.Fatalf("CommodityOptionMap[au] = %#v", commodity["au"])
	}
}

func TestRootReExportsTSCoreConstantNames(t *testing.T) {
	if TENCENT_BASE_URL != TencentBaseURL || EM_KLINE_URL != EMKlineURL || EM_PUSH_TOKEN != EMPushToken {
		t.Fatalf("TS-style URL/token constants are not wired")
	}
	if SINA_OPTION_API_URL != SinaOptionAPIURL || EM_OPTION_CFFEX_URL != EMOptionCFFEXURL {
		t.Fatalf("TS-style option URL constants are not wired")
	}
	if DEFAULT_TIMEOUT != DefaultTimeoutMS || DEFAULT_MAX_RETRIES != DefaultMaxRetries || DEFAULT_BASE_DELAY != DefaultBaseDelayMS {
		t.Fatalf("TS-style default constants are not wired")
	}
	if FUTURES_EXCHANGE_MAP["CFFEX"] != 220 || FUTURES_VARIETY_EXCHANGE["IF"] != "CFFEX" || GLOBAL_FUTURES_MARKET["GC"] != 101 {
		t.Fatalf("TS-style futures maps mismatch")
	}
	if CFFEX_OPTION_PRODUCT_MAP["io"] != "沪深300" || COMMODITY_OPTION_MAP["SR"] != (CommodityOption{Product: "SR_o", Exchange: "czce"}) {
		t.Fatalf("TS-style option maps mismatch")
	}
	if DEFAULT_RETRYABLE_STATUS_CODES[0] != 408 || DEFAULT_RETRYABLE_STATUS_CODES[len(DEFAULT_RETRYABLE_STATUS_CODES)-1] != 504 {
		t.Fatalf("DEFAULT_RETRYABLE_STATUS_CODES = %#v", DEFAULT_RETRYABLE_STATUS_CODES)
	}
}
