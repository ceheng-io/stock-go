package types

// Quote is the common interface implemented by concrete quote payloads.
//
// It mirrors the TypeScript union:
// FullQuote | HKQuote | USQuote | FundQuote.
type Quote interface {
	isQuote()
}

// TencentQuoteItem is a raw Tencent quote assignment split by "~".
type TencentQuoteItem struct {
	Key    string   `json:"key"`
	Fields []string `json:"fields"`
}

// SimpleQuote is a compact CN quote returned by Tencent simple quote APIs.
type SimpleQuote struct {
	MarketID      string   `json:"marketId"`
	Name          string   `json:"name"`
	Code          string   `json:"code"`
	Price         float64  `json:"price"`
	Change        float64  `json:"change"`
	ChangePercent float64  `json:"changePercent"`
	Volume        float64  `json:"volume"`
	Amount        float64  `json:"amount"`
	MarketCap     *float64 `json:"marketCap"`
	MarketType    string   `json:"marketType"`
	Market        Market   `json:"market"`
	AssetType     string   `json:"assetType"`
	Source        string   `json:"source"`
}

// PriceLevel is a bid or ask level.
type PriceLevel struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

// FullQuote is a detailed CN quote returned by Tencent quote APIs.
type FullQuote struct {
	MarketID             string       `json:"marketId"`
	Name                 string       `json:"name"`
	Code                 string       `json:"code"`
	Price                float64      `json:"price"`
	PrevClose            float64      `json:"prevClose"`
	Open                 float64      `json:"open"`
	Volume               float64      `json:"volume"`
	OuterVolume          float64      `json:"outerVolume"`
	InnerVolume          float64      `json:"innerVolume"`
	Bid                  []PriceLevel `json:"bid"`
	Ask                  []PriceLevel `json:"ask"`
	Time                 string       `json:"time"`
	Timestamp            *int64       `json:"timestamp"`
	TZ                   string       `json:"tz"`
	Change               float64      `json:"change"`
	ChangePercent        float64      `json:"changePercent"`
	High                 float64      `json:"high"`
	Low                  float64      `json:"low"`
	Volume2              float64      `json:"volume2"`
	Amount               float64      `json:"amount"`
	TurnoverRate         *float64     `json:"turnoverRate"`
	PE                   *float64     `json:"pe"`
	Amplitude            *float64     `json:"amplitude"`
	CirculatingMarketCap *float64     `json:"circulatingMarketCap"`
	TotalMarketCap       *float64     `json:"totalMarketCap"`
	PB                   *float64     `json:"pb"`
	LimitUp              *float64     `json:"limitUp"`
	LimitDown            *float64     `json:"limitDown"`
	VolumeRatio          *float64     `json:"volumeRatio"`
	AvgPrice             *float64     `json:"avgPrice"`
	PEStatic             *float64     `json:"peStatic"`
	PEDynamic            *float64     `json:"peDynamic"`
	High52W              *float64     `json:"high52W"`
	Low52W               *float64     `json:"low52W"`
	CirculatingShares    *float64     `json:"circulatingShares"`
	TotalShares          *float64     `json:"totalShares"`
	Market               Market       `json:"market"`
	AssetType            string       `json:"assetType"`
	Source               string       `json:"source"`
}

func (FullQuote) isQuote() {}

// HKQuote is a Hong Kong stock quote.
type HKQuote struct {
	MarketID             string   `json:"marketId"`
	Name                 string   `json:"name"`
	Code                 string   `json:"code"`
	Price                float64  `json:"price"`
	PrevClose            float64  `json:"prevClose"`
	Open                 float64  `json:"open"`
	Volume               float64  `json:"volume"`
	Time                 string   `json:"time"`
	Timestamp            *int64   `json:"timestamp"`
	TZ                   string   `json:"tz"`
	Change               float64  `json:"change"`
	ChangePercent        float64  `json:"changePercent"`
	High                 float64  `json:"high"`
	Low                  float64  `json:"low"`
	Amount               float64  `json:"amount"`
	LotSize              *float64 `json:"lotSize"`
	CirculatingMarketCap *float64 `json:"circulatingMarketCap"`
	TotalMarketCap       *float64 `json:"totalMarketCap"`
	Currency             string   `json:"currency"`
	Market               Market   `json:"market"`
	AssetType            string   `json:"assetType"`
	Source               string   `json:"source"`
}

func (HKQuote) isQuote() {}

// USQuote is a US stock quote.
type USQuote struct {
	MarketID       string   `json:"marketId"`
	Name           string   `json:"name"`
	Code           string   `json:"code"`
	Price          float64  `json:"price"`
	PrevClose      float64  `json:"prevClose"`
	Open           float64  `json:"open"`
	Volume         float64  `json:"volume"`
	Time           string   `json:"time"`
	Timestamp      *int64   `json:"timestamp"`
	TZ             string   `json:"tz"`
	Change         float64  `json:"change"`
	ChangePercent  float64  `json:"changePercent"`
	High           float64  `json:"high"`
	Low            float64  `json:"low"`
	Amount         float64  `json:"amount"`
	TurnoverRate   *float64 `json:"turnoverRate"`
	PE             *float64 `json:"pe"`
	Amplitude      *float64 `json:"amplitude"`
	TotalMarketCap *float64 `json:"totalMarketCap"`
	PB             *float64 `json:"pb"`
	High52W        *float64 `json:"high52W"`
	Low52W         *float64 `json:"low52W"`
	Market         Market   `json:"market"`
	AssetType      string   `json:"assetType"`
	Source         string   `json:"source"`
}

func (USQuote) isQuote() {}

// FundQuote is a public fund quote.
type FundQuote struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	NAV       float64 `json:"nav"`
	AccNAV    float64 `json:"accNav"`
	Change    float64 `json:"change"`
	NavDate   string  `json:"navDate"`
	Timestamp *int64  `json:"timestamp"`
	TZ        string  `json:"tz"`
	Market    Market  `json:"market"`
	AssetType string  `json:"assetType"`
	Source    string  `json:"source"`
}

func (FundQuote) isQuote() {}

// FundFlow 是腾讯资金流向数据。
type FundFlow struct {
	Code           string  `json:"code"`
	MainInflow     float64 `json:"mainInflow"`
	MainOutflow    float64 `json:"mainOutflow"`
	MainNet        float64 `json:"mainNet"`
	MainNetRatio   float64 `json:"mainNetRatio"`
	RetailInflow   float64 `json:"retailInflow"`
	RetailOutflow  float64 `json:"retailOutflow"`
	RetailNet      float64 `json:"retailNet"`
	RetailNetRatio float64 `json:"retailNetRatio"`
	TotalFlow      float64 `json:"totalFlow"`
	Name           string  `json:"name"`
	Date           string  `json:"date"`
	Timestamp      *int64  `json:"timestamp"`
	TZ             string  `json:"tz"`
}

// PanelLargeOrder 是腾讯盘口大单占比数据。
type PanelLargeOrder struct {
	BuyLargeRatio  float64 `json:"buyLargeRatio"`
	BuySmallRatio  float64 `json:"buySmallRatio"`
	SellLargeRatio float64 `json:"sellLargeRatio"`
	SellSmallRatio float64 `json:"sellSmallRatio"`
}

// SearchResultType is a normalized search asset category.
type SearchResultType string

const (
	SearchStock   SearchResultType = "stock"
	SearchIndex   SearchResultType = "index"
	SearchFund    SearchResultType = "fund"
	SearchBond    SearchResultType = "bond"
	SearchFutures SearchResultType = "futures"
	SearchOption  SearchResultType = "option"
	SearchOther   SearchResultType = "other"
)

// SearchResult is a Tencent Smartbox search result.
type SearchResult struct {
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Market   string           `json:"market"`
	Type     string           `json:"type"`
	Category SearchResultType `json:"category"`
}
