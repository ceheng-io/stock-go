package stock

import "github.com/ceheng.io/stock-go/constants"

type CommodityOption = constants.CommodityOption

const (
	TencentBaseURL     = constants.TencentBaseURL
	TencentMinuteURL   = constants.TencentMinuteURL
	TENCENT_BASE_URL   = constants.TENCENT_BASE_URL
	TENCENT_MINUTE_URL = constants.TENCENT_MINUTE_URL

	AShareListURL      = constants.AShareListURL
	USListURL          = constants.USListURL
	HKListURL          = constants.HKListURL
	FundListURL        = constants.FundListURL
	TradeCalendarURL   = constants.TradeCalendarURL
	CodeListURL        = constants.CodeListURL
	A_SHARE_LIST_URL   = constants.A_SHARE_LIST_URL
	US_LIST_URL        = constants.US_LIST_URL
	HK_LIST_URL        = constants.HK_LIST_URL
	FUND_LIST_URL      = constants.FUND_LIST_URL
	TRADE_CALENDAR_URL = constants.TRADE_CALENDAR_URL
	CODE_LIST_URL      = constants.CODE_LIST_URL

	EMKlineURL               = constants.EMKlineURL
	EMTrendsURL              = constants.EMTrendsURL
	EMHKKlineURL             = constants.EMHKKlineURL
	EMUSKlineURL             = constants.EMUSKlineURL
	EMHKTrendsURL            = constants.EMHKTrendsURL
	EMUSTrendsURL            = constants.EMUSTrendsURL
	EMBoardListURL           = constants.EMBoardListURL
	EMBoardSpotURL           = constants.EMBoardSpotURL
	EMBoardConsURL           = constants.EMBoardConsURL
	EMBoardKlineURL          = constants.EMBoardKlineURL
	EMBoardTrendsURL         = constants.EMBoardTrendsURL
	EMConceptListURL         = constants.EMConceptListURL
	EMConceptSpotURL         = constants.EMConceptSpotURL
	EMConceptConsURL         = constants.EMConceptConsURL
	EMConceptKlineURL        = constants.EMConceptKlineURL
	EMConceptTrendsURL       = constants.EMConceptTrendsURL
	EMDatacenterURL          = constants.EMDatacenterURL
	EMF10BaseURL             = constants.EMF10BaseURL
	EMAnnouncementListURL    = constants.EMAnnouncementListURL
	EMAnnouncementURL        = constants.EMAnnouncementURL
	EMFFlowURL               = constants.EMFFlowURL
	EMClistURL               = constants.EMClistURL
	EMNorthboundMinuteURL    = constants.EMNorthboundMinuteURL
	EMTopicBaseURL           = constants.EMTopicBaseURL
	EMPushToken              = constants.EMPushToken
	EMDataToken              = constants.EMDataToken
	EM_KLINE_URL             = constants.EM_KLINE_URL
	EM_TRENDS_URL            = constants.EM_TRENDS_URL
	EM_HK_KLINE_URL          = constants.EM_HK_KLINE_URL
	EM_US_KLINE_URL          = constants.EM_US_KLINE_URL
	EM_HK_TRENDS_URL         = constants.EM_HK_TRENDS_URL
	EM_US_TRENDS_URL         = constants.EM_US_TRENDS_URL
	EM_BOARD_LIST_URL        = constants.EM_BOARD_LIST_URL
	EM_BOARD_SPOT_URL        = constants.EM_BOARD_SPOT_URL
	EM_BOARD_CONS_URL        = constants.EM_BOARD_CONS_URL
	EM_BOARD_KLINE_URL       = constants.EM_BOARD_KLINE_URL
	EM_BOARD_TRENDS_URL      = constants.EM_BOARD_TRENDS_URL
	EM_CONCEPT_LIST_URL      = constants.EM_CONCEPT_LIST_URL
	EM_CONCEPT_SPOT_URL      = constants.EM_CONCEPT_SPOT_URL
	EM_CONCEPT_CONS_URL      = constants.EM_CONCEPT_CONS_URL
	EM_CONCEPT_KLINE_URL     = constants.EM_CONCEPT_KLINE_URL
	EM_CONCEPT_TRENDS_URL    = constants.EM_CONCEPT_TRENDS_URL
	EM_DATACENTER_URL        = constants.EM_DATACENTER_URL
	EM_F10_BASE_URL          = constants.EM_F10_BASE_URL
	EM_ANNOUNCEMENT_LIST_URL = constants.EM_ANNOUNCEMENT_LIST_URL
	EM_ANNOUNCEMENT_URL      = constants.EM_ANNOUNCEMENT_URL
	EM_FFLOW_URL             = constants.EM_FFLOW_URL
	EM_CLIST_URL             = constants.EM_CLIST_URL
	EM_NORTHBOUND_MINUTE_URL = constants.EM_NORTHBOUND_MINUTE_URL
	EM_TOPIC_BASE_URL        = constants.EM_TOPIC_BASE_URL
	EM_PUSH_TOKEN            = constants.EM_PUSH_TOKEN
	EM_DATA_TOKEN            = constants.EM_DATA_TOKEN

	EMFuturesKlineURL            = constants.EMFuturesKlineURL
	EMFuturesGlobalSpotURL       = constants.EMFuturesGlobalSpotURL
	EMFuturesGlobalSpotToken     = constants.EMFuturesGlobalSpotToken
	EMOptionCFFEXURL             = constants.EMOptionCFFEXURL
	EMOptionLHBURL               = constants.EMOptionLHBURL
	EMOptionLHBToken             = constants.EMOptionLHBToken
	EM_FUTURES_KLINE_URL         = constants.EM_FUTURES_KLINE_URL
	EM_FUTURES_GLOBAL_SPOT_URL   = constants.EM_FUTURES_GLOBAL_SPOT_URL
	EM_FUTURES_GLOBAL_SPOT_TOKEN = constants.EM_FUTURES_GLOBAL_SPOT_TOKEN
	EM_OPTION_CFFEX_URL          = constants.EM_OPTION_CFFEX_URL
	EM_OPTION_LHB_URL            = constants.EM_OPTION_LHB_URL
	EM_OPTION_LHB_TOKEN          = constants.EM_OPTION_LHB_TOKEN

	SinaOptionAPIURL           = constants.SinaOptionAPIURL
	SinaOptionDaylineURL       = constants.SinaOptionDaylineURL
	SinaSSEOptionListURL       = constants.SinaSSEOptionListURL
	SinaSSEOptionExpireURL     = constants.SinaSSEOptionExpireURL
	SinaSSEOptionMinuteURL     = constants.SinaSSEOptionMinuteURL
	SinaSSEOptionDailyURL      = constants.SinaSSEOptionDailyURL
	SinaSSEOption5DayURL       = constants.SinaSSEOption5DayURL
	SINA_OPTION_API_URL        = constants.SINA_OPTION_API_URL
	SINA_OPTION_DAYLINE_URL    = constants.SINA_OPTION_DAYLINE_URL
	SINA_SSE_OPTION_LIST_URL   = constants.SINA_SSE_OPTION_LIST_URL
	SINA_SSE_OPTION_EXPIRE_URL = constants.SINA_SSE_OPTION_EXPIRE_URL
	SINA_SSE_OPTION_MINUTE_URL = constants.SINA_SSE_OPTION_MINUTE_URL
	SINA_SSE_OPTION_DAILY_URL  = constants.SINA_SSE_OPTION_DAILY_URL
	SINA_SSE_OPTION_5DAY_URL   = constants.SINA_SSE_OPTION_5DAY_URL

	THSLimitUpPoolURL     = constants.THSLimitUpPoolURL
	THS_LIMIT_UP_POOL_URL = constants.THS_LIMIT_UP_POOL_URL

	DefaultTimeoutMS           = constants.DefaultTimeoutMS
	DefaultBatchSize           = constants.DefaultBatchSize
	MaxBatchSize               = constants.MaxBatchSize
	DefaultConcurrency         = constants.DefaultConcurrency
	DefaultMaxRetries          = constants.DefaultMaxRetries
	DefaultBaseDelayMS         = constants.DefaultBaseDelayMS
	DefaultMaxDelayMS          = constants.DefaultMaxDelayMS
	DefaultBackoffMultiplier   = constants.DefaultBackoffMultiplier
	DEFAULT_TIMEOUT            = constants.DEFAULT_TIMEOUT
	DEFAULT_BATCH_SIZE         = constants.DEFAULT_BATCH_SIZE
	MAX_BATCH_SIZE             = constants.MAX_BATCH_SIZE
	DEFAULT_CONCURRENCY        = constants.DEFAULT_CONCURRENCY
	DEFAULT_MAX_RETRIES        = constants.DEFAULT_MAX_RETRIES
	DEFAULT_BASE_DELAY         = constants.DEFAULT_BASE_DELAY
	DEFAULT_MAX_DELAY          = constants.DEFAULT_MAX_DELAY
	DEFAULT_BACKOFF_MULTIPLIER = constants.DEFAULT_BACKOFF_MULTIPLIER
)

var (
	DEFAULT_RETRYABLE_STATUS_CODES = constants.DefaultRetryableStatusCodes()
	FUTURES_EXCHANGE_MAP           = constants.FuturesExchangeMap()
	FUTURES_VARIETY_EXCHANGE       = constants.FuturesVarietyExchange()
	GLOBAL_FUTURES_MARKET          = constants.GlobalFuturesMarket()
	CFFEX_OPTION_PRODUCT_MAP       = constants.CFFEXOptionProductMap()
	COMMODITY_OPTION_MAP           = constants.CommodityOptionMap()
)

// DefaultRetryableStatusCodes 返回默认可重试 HTTP 状态码。
func DefaultRetryableStatusCodes() []int {
	return constants.DefaultRetryableStatusCodes()
}

// FuturesExchangeMap 返回国内期货交易所 market code 映射。
func FuturesExchangeMap() map[string]int {
	return constants.FuturesExchangeMap()
}

// FuturesVarietyExchange 返回期货品种代码到交易所的映射。
func FuturesVarietyExchange() map[string]string {
	return constants.FuturesVarietyExchange()
}

// GlobalFuturesMarket 返回全球期货市场代码映射。
func GlobalFuturesMarket() map[string]int {
	return constants.GlobalFuturesMarket()
}

// CFFEXOptionProductMap 返回中金所股指期权品种到 product 名称的映射。
func CFFEXOptionProductMap() map[string]string {
	return constants.CFFEXOptionProductMap()
}

// CommodityOptionMap 返回商品期权品种到新浪 product 与交易所参数的映射。
func CommodityOptionMap() map[string]CommodityOption {
	return constants.CommodityOptionMap()
}
