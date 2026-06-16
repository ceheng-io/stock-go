package constants

// 腾讯财经 API。
const (
	TencentBaseURL   = "https://qt.gtimg.cn"
	TencentMinuteURL = "https://web.ifzq.gtimg.cn/appstock/app/minute/query"

	TENCENT_BASE_URL   = TencentBaseURL
	TENCENT_MINUTE_URL = TencentMinuteURL
)

// 股票代码列表和交易日历地址。
const (
	AShareListURL    = "https://assets.linkdiary.cn/shares/zh_a_list.json"
	USListURL        = "https://assets.linkdiary.cn/shares/us_list.json"
	HKListURL        = "https://assets.linkdiary.cn/shares/hk_list.json"
	FundListURL      = "https://assets.linkdiary.cn/shares/fund_list"
	TradeCalendarURL = "https://assets.linkdiary.cn/shares/trade-data-list.txt"
	CodeListURL      = AShareListURL

	A_SHARE_LIST_URL   = AShareListURL
	US_LIST_URL        = USListURL
	HK_LIST_URL        = HKListURL
	FUND_LIST_URL      = FundListURL
	TRADE_CALENDAR_URL = TradeCalendarURL
	CODE_LIST_URL      = CodeListURL
)

// 东方财富 API。
const (
	EMKlineURL            = "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	EMTrendsURL           = "https://push2his.eastmoney.com/api/qt/stock/trends2/get"
	EMHKKlineURL          = "https://33.push2his.eastmoney.com/api/qt/stock/kline/get"
	EMUSKlineURL          = "https://63.push2his.eastmoney.com/api/qt/stock/kline/get"
	EMHKTrendsURL         = "https://33.push2his.eastmoney.com/api/qt/stock/trends2/get"
	EMUSTrendsURL         = "https://63.push2his.eastmoney.com/api/qt/stock/trends2/get"
	EMBoardListURL        = "https://17.push2.eastmoney.com/api/qt/clist/get"
	EMBoardSpotURL        = "https://91.push2.eastmoney.com/api/qt/stock/get"
	EMBoardConsURL        = "https://29.push2.eastmoney.com/api/qt/clist/get"
	EMBoardKlineURL       = "https://7.push2his.eastmoney.com/api/qt/stock/kline/get"
	EMBoardTrendsURL      = "https://push2his.eastmoney.com/api/qt/stock/trends2/get"
	EMConceptListURL      = "https://79.push2.eastmoney.com/api/qt/clist/get"
	EMConceptSpotURL      = "https://91.push2.eastmoney.com/api/qt/stock/get"
	EMConceptConsURL      = "https://29.push2.eastmoney.com/api/qt/clist/get"
	EMConceptKlineURL     = "https://91.push2his.eastmoney.com/api/qt/stock/kline/get"
	EMConceptTrendsURL    = "https://push2his.eastmoney.com/api/qt/stock/trends2/get"
	EMDatacenterURL       = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	EMF10BaseURL          = "https://emweb.securities.eastmoney.com"
	EMAnnouncementListURL = "https://np-anotice-stock.eastmoney.com/api/security/ann"
	EMAnnouncementURL     = "https://np-cnotice-stock.eastmoney.com/api/content/ann"
	EMFFlowURL            = "https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get"
	EMClistURL            = "https://push2.eastmoney.com/api/qt/clist/get"
	EMNorthboundMinuteURL = "https://push2.eastmoney.com/api/qt/kamtbs.rtmin/get"
	EMTopicBaseURL        = "https://push2ex.eastmoney.com"
	EMPushToken           = "7eea3edcaed734bea9cbfc24409ed989"
	EMDataToken           = "b2884a393a59ad64002292a3e90d46a5"

	EM_KLINE_URL             = EMKlineURL
	EM_TRENDS_URL            = EMTrendsURL
	EM_HK_KLINE_URL          = EMHKKlineURL
	EM_US_KLINE_URL          = EMUSKlineURL
	EM_HK_TRENDS_URL         = EMHKTrendsURL
	EM_US_TRENDS_URL         = EMUSTrendsURL
	EM_BOARD_LIST_URL        = EMBoardListURL
	EM_BOARD_SPOT_URL        = EMBoardSpotURL
	EM_BOARD_CONS_URL        = EMBoardConsURL
	EM_BOARD_KLINE_URL       = EMBoardKlineURL
	EM_BOARD_TRENDS_URL      = EMBoardTrendsURL
	EM_CONCEPT_LIST_URL      = EMConceptListURL
	EM_CONCEPT_SPOT_URL      = EMConceptSpotURL
	EM_CONCEPT_CONS_URL      = EMConceptConsURL
	EM_CONCEPT_KLINE_URL     = EMConceptKlineURL
	EM_CONCEPT_TRENDS_URL    = EMConceptTrendsURL
	EM_DATACENTER_URL        = EMDatacenterURL
	EM_F10_BASE_URL          = EMF10BaseURL
	EM_ANNOUNCEMENT_LIST_URL = EMAnnouncementListURL
	EM_ANNOUNCEMENT_URL      = EMAnnouncementURL
	EM_FFLOW_URL             = EMFFlowURL
	EM_CLIST_URL             = EMClistURL
	EM_NORTHBOUND_MINUTE_URL = EMNorthboundMinuteURL
	EM_TOPIC_BASE_URL        = EMTopicBaseURL
	EM_PUSH_TOKEN            = EMPushToken
	EM_DATA_TOKEN            = EMDataToken
)

// 东方财富期货和期权 API。
const (
	EMFuturesKlineURL        = "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	EMFuturesGlobalSpotURL   = "https://futsseapi.eastmoney.com/list/COMEX,NYMEX,COBOT,SGX,NYBOT,LME,MDEX,TOCOM,IPE"
	EMFuturesGlobalSpotToken = "58b2fa8f54638b60b87d69b31969089c"
	EMOptionCFFEXURL         = "https://futsseapi.eastmoney.com/list/option/221"
	EMOptionLHBURL           = "https://datacenter-web.eastmoney.com/api/data/get"
	EMOptionLHBToken         = EMDataToken

	EM_FUTURES_KLINE_URL         = EMFuturesKlineURL
	EM_FUTURES_GLOBAL_SPOT_URL   = EMFuturesGlobalSpotURL
	EM_FUTURES_GLOBAL_SPOT_TOKEN = EMFuturesGlobalSpotToken
	EM_OPTION_CFFEX_URL          = EMOptionCFFEXURL
	EM_OPTION_LHB_URL            = EMOptionLHBURL
	EM_OPTION_LHB_TOKEN          = EMOptionLHBToken
)

// 新浪期权 API。
const (
	SinaOptionAPIURL       = "https://stock.finance.sina.com.cn/futures/api/openapi.php/OptionService.getOptionData"
	SinaOptionDaylineURL   = "https://stock.finance.sina.com.cn/futures/api/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline"
	SinaSSEOptionListURL   = "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionService.getStockName"
	SinaSSEOptionExpireURL = "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionService.getRemainderDay"
	SinaSSEOptionMinuteURL = "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionDaylineService.getOptionMinline"
	SinaSSEOptionDailyURL  = "https://stock.finance.sina.com.cn/futures/api/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo"
	SinaSSEOption5DayURL   = "https://stock.finance.sina.com.cn/futures/api/openapi.php/StockOptionDaylineService.getFiveDayLine"

	SINA_OPTION_API_URL        = SinaOptionAPIURL
	SINA_OPTION_DAYLINE_URL    = SinaOptionDaylineURL
	SINA_SSE_OPTION_LIST_URL   = SinaSSEOptionListURL
	SINA_SSE_OPTION_EXPIRE_URL = SinaSSEOptionExpireURL
	SINA_SSE_OPTION_MINUTE_URL = SinaSSEOptionMinuteURL
	SINA_SSE_OPTION_DAILY_URL  = SinaSSEOptionDailyURL
	SINA_SSE_OPTION_5DAY_URL   = SinaSSEOption5DayURL
)

// 同花顺数据 API。
const (
	THSLimitUpPoolURL = "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool"

	THS_LIMIT_UP_POOL_URL = THSLimitUpPoolURL
)

// 默认请求和重试配置，单位与 TS SDK 保持一致。
const (
	DefaultTimeoutMS         = 30000
	DefaultBatchSize         = 500
	MaxBatchSize             = 500
	DefaultConcurrency       = 7
	DefaultMaxRetries        = 3
	DefaultBaseDelayMS       = 1000
	DefaultMaxDelayMS        = 30000
	DefaultBackoffMultiplier = 2

	DEFAULT_TIMEOUT            = DefaultTimeoutMS
	DEFAULT_BATCH_SIZE         = DefaultBatchSize
	MAX_BATCH_SIZE             = MaxBatchSize
	DEFAULT_CONCURRENCY        = DefaultConcurrency
	DEFAULT_MAX_RETRIES        = DefaultMaxRetries
	DEFAULT_BASE_DELAY         = DefaultBaseDelayMS
	DEFAULT_MAX_DELAY          = DefaultMaxDelayMS
	DEFAULT_BACKOFF_MULTIPLIER = DefaultBackoffMultiplier
)

var defaultRetryableStatusCodes = []int{408, 429, 500, 502, 503, 504}

// DEFAULT_RETRYABLE_STATUS_CODES keeps the TypeScript SDK constant name available.
var DEFAULT_RETRYABLE_STATUS_CODES = DefaultRetryableStatusCodes()

// DefaultRetryableStatusCodes 返回默认可重试 HTTP 状态码。
func DefaultRetryableStatusCodes() []int {
	return append([]int(nil), defaultRetryableStatusCodes...)
}

var futuresExchangeMap = map[string]int{
	"SHFE":  113,
	"DCE":   114,
	"CZCE":  115,
	"INE":   142,
	"CFFEX": 220,
	"GFEX":  225,
}

// FUTURES_EXCHANGE_MAP keeps the TypeScript SDK constant name available.
var FUTURES_EXCHANGE_MAP = FuturesExchangeMap()

// FuturesExchangeMap 返回国内期货交易所 market code 映射。
func FuturesExchangeMap() map[string]int {
	return cloneMap(futuresExchangeMap)
}

var futuresVarietyExchange = map[string]string{
	"cu": "SHFE", "al": "SHFE", "zn": "SHFE", "pb": "SHFE", "au": "SHFE", "ag": "SHFE",
	"rb": "SHFE", "wr": "SHFE", "fu": "SHFE", "ru": "SHFE", "bu": "SHFE", "hc": "SHFE",
	"ni": "SHFE", "sn": "SHFE", "sp": "SHFE", "ss": "SHFE", "ao": "SHFE", "br": "SHFE",
	"c": "DCE", "a": "DCE", "b": "DCE", "m": "DCE", "y": "DCE", "p": "DCE",
	"l": "DCE", "v": "DCE", "j": "DCE", "jm": "DCE", "i": "DCE", "jd": "DCE",
	"pp": "DCE", "cs": "DCE", "eg": "DCE", "eb": "DCE", "pg": "DCE", "lh": "DCE",
	"WH": "CZCE", "CF": "CZCE", "SR": "CZCE", "TA": "CZCE", "OI": "CZCE", "MA": "CZCE",
	"FG": "CZCE", "RM": "CZCE", "SF": "CZCE", "SM": "CZCE", "ZC": "CZCE", "AP": "CZCE",
	"CJ": "CZCE", "UR": "CZCE", "SA": "CZCE", "PF": "CZCE", "PK": "CZCE", "PX": "CZCE",
	"SH": "CZCE",
	"sc": "INE", "nr": "INE", "lu": "INE", "bc": "INE", "ec": "INE",
	"IF": "CFFEX", "IC": "CFFEX", "IH": "CFFEX", "IM": "CFFEX",
	"TS": "CFFEX", "TF": "CFFEX", "T": "CFFEX", "TL": "CFFEX",
	"si": "GFEX", "lc": "GFEX", "ps": "GFEX", "pt": "GFEX", "pd": "GFEX",
}

// FUTURES_VARIETY_EXCHANGE keeps the TypeScript SDK constant name available.
var FUTURES_VARIETY_EXCHANGE = FuturesVarietyExchange()

// FuturesVarietyExchange 返回期货品种代码到交易所的映射。
func FuturesVarietyExchange() map[string]string {
	return cloneMap(futuresVarietyExchange)
}

var globalFuturesMarket = map[string]int{
	"HG": 101, "GC": 101, "SI": 101, "QI": 101, "QO": 101, "MGC": 101,
	"CL": 102, "NG": 102, "RB": 102, "HO": 102, "PA": 102, "PL": 102,
	"ZW": 103, "ZM": 103, "ZS": 103, "ZC": 103, "ZL": 103, "ZR": 103,
	"YM": 103, "NQ": 103, "ES": 103,
	"SB": 108, "CT": 108,
	"LCPT": 109, "LZNT": 109, "LALT": 109,
}

// GLOBAL_FUTURES_MARKET keeps the TypeScript SDK constant name available.
var GLOBAL_FUTURES_MARKET = GlobalFuturesMarket()

// GlobalFuturesMarket 返回全球期货市场代码映射。
func GlobalFuturesMarket() map[string]int {
	return cloneMap(globalFuturesMarket)
}

var cffexOptionProductMap = map[string]string{
	"ho": "上证50",
	"io": "沪深300",
	"mo": "中证1000",
}

// CFFEX_OPTION_PRODUCT_MAP keeps the TypeScript SDK constant name available.
var CFFEX_OPTION_PRODUCT_MAP = CFFEXOptionProductMap()

// CFFEXOptionProductMap 返回中金所股指期权品种到 product 名称的映射。
func CFFEXOptionProductMap() map[string]string {
	return cloneMap(cffexOptionProductMap)
}

// CommodityOption 描述商品期权的新浪 product 与交易所参数。
type CommodityOption struct {
	Product  string
	Exchange string
}

var commodityOptionMap = map[string]CommodityOption{
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

// COMMODITY_OPTION_MAP keeps the TypeScript SDK constant name available.
var COMMODITY_OPTION_MAP = CommodityOptionMap()

// CommodityOptionMap 返回商品期权品种到新浪 product 与交易所参数的映射。
func CommodityOptionMap() map[string]CommodityOption {
	return cloneMap(commodityOptionMap)
}

func cloneMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
