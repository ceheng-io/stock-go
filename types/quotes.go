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
	Key    string
	Fields []string
}

// SimpleQuote is a compact CN quote returned by Tencent simple quote APIs.
type SimpleQuote struct {
	MarketID      string
	Name          string
	Code          string
	Price         float64
	Change        float64
	ChangePercent float64
	Volume        float64
	Amount        float64
	MarketCap     *float64
	MarketType    string
	Market        Market
	AssetType     string
	Source        string
}

// PriceLevel is a bid or ask level.
type PriceLevel struct {
	Price  float64
	Volume float64
}

// FullQuote is a detailed CN quote returned by Tencent quote APIs.
type FullQuote struct {
	MarketID             string
	Name                 string
	Code                 string
	Price                float64
	PrevClose            float64
	Open                 float64
	Volume               float64
	OuterVolume          float64
	InnerVolume          float64
	Bid                  []PriceLevel
	Ask                  []PriceLevel
	Time                 string
	Timestamp            *int64
	TZ                   string
	Change               float64
	ChangePercent        float64
	High                 float64
	Low                  float64
	Volume2              float64
	Amount               float64
	TurnoverRate         *float64
	PE                   *float64
	Amplitude            *float64
	CirculatingMarketCap *float64
	TotalMarketCap       *float64
	PB                   *float64
	LimitUp              *float64
	LimitDown            *float64
	VolumeRatio          *float64
	AvgPrice             *float64
	PEStatic             *float64
	PEDynamic            *float64
	High52W              *float64
	Low52W               *float64
	CirculatingShares    *float64
	TotalShares          *float64
	Market               Market
	AssetType            string
	Source               string
}

func (FullQuote) isQuote() {}

// HKQuote is a Hong Kong stock quote.
type HKQuote struct {
	MarketID             string
	Name                 string
	Code                 string
	Price                float64
	PrevClose            float64
	Open                 float64
	Volume               float64
	Time                 string
	Timestamp            *int64
	TZ                   string
	Change               float64
	ChangePercent        float64
	High                 float64
	Low                  float64
	Amount               float64
	LotSize              *float64
	CirculatingMarketCap *float64
	TotalMarketCap       *float64
	Currency             string
	Market               Market
	AssetType            string
	Source               string
}

func (HKQuote) isQuote() {}

// USQuote is a US stock quote.
type USQuote struct {
	MarketID       string
	Name           string
	Code           string
	Price          float64
	PrevClose      float64
	Open           float64
	Volume         float64
	Time           string
	Timestamp      *int64
	TZ             string
	Change         float64
	ChangePercent  float64
	High           float64
	Low            float64
	Amount         float64
	TurnoverRate   *float64
	PE             *float64
	Amplitude      *float64
	TotalMarketCap *float64
	PB             *float64
	High52W        *float64
	Low52W         *float64
	Market         Market
	AssetType      string
	Source         string
}

func (USQuote) isQuote() {}

// FundQuote is a public fund quote.
type FundQuote struct {
	Code      string
	Name      string
	NAV       float64
	AccNAV    float64
	Change    float64
	NavDate   string
	Timestamp *int64
	TZ        string
	Market    Market
	AssetType string
	Source    string
}

func (FundQuote) isQuote() {}

// FundFlow 是腾讯资金流向数据。
type FundFlow struct {
	Code           string
	MainInflow     float64
	MainOutflow    float64
	MainNet        float64
	MainNetRatio   float64
	RetailInflow   float64
	RetailOutflow  float64
	RetailNet      float64
	RetailNetRatio float64
	TotalFlow      float64
	Name           string
	Date           string
	Timestamp      *int64
	TZ             string
}

// PanelLargeOrder 是腾讯盘口大单占比数据。
type PanelLargeOrder struct {
	BuyLargeRatio  float64
	BuySmallRatio  float64
	SellLargeRatio float64
	SellSmallRatio float64
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
	Code     string
	Name     string
	Market   string
	Type     string
	Category SearchResultType
}
