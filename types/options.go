package types

// IndexOptionProduct 是中金所股指期权产品。
type IndexOptionProduct string

const (
	IndexOptionProductHO IndexOptionProduct = "ho"
	IndexOptionProductIO IndexOptionProduct = "io"
	IndexOptionProductMO IndexOptionProduct = "mo"
)

// ETFOptionCate 是新浪 ETF 期权分类。
type ETFOptionCate string

const (
	ETFOptionCate50ETF           ETFOptionCate = "50ETF"
	ETFOptionCate300ETF          ETFOptionCate = "300ETF"
	ETFOptionCate500ETF          ETFOptionCate = "500ETF"
	ETFOptionCateKechuang50      ETFOptionCate = "科创50"
	ETFOptionCateKechuangBoard50 ETFOptionCate = "科创板50"
)

// ETFOptionMonth 是新浪 ETF 期权可用月份信息。
type ETFOptionMonth struct {
	Months   []string
	StockID  string
	CateID   string
	CateList []string
}

// ETFOptionExpireDay 是新浪 ETF 期权到期日信息。
type ETFOptionExpireDay struct {
	ExpireDay     string
	RemainderDays int
	StockID       string
	Name          string
}

// OptionKline 是期权日 K 线。
type OptionKline struct {
	Date   string
	Open   *float64
	High   *float64
	Low    *float64
	Close  *float64
	Volume *float64
}

// OptionMinute 是期权分钟行情。
type OptionMinute struct {
	Time         string
	Date         string
	Price        *float64
	Volume       *float64
	OpenInterest *float64
	AvgPrice     *float64
}

// OptionTQuote 是期权 T 型报价项。
type OptionTQuote struct {
	Symbol       string
	BuyVolume    *float64
	BuyPrice     *float64
	Price        *float64
	AskPrice     *float64
	AskVolume    *float64
	OpenInterest *float64
	Change       *float64
	StrikePrice  *float64
}

// OptionTQuoteResult 是期权 T 型报价结果。
type OptionTQuoteResult struct {
	Calls []OptionTQuote
	Puts  []OptionTQuote
}

// CFFEXOptionQuote 是中金所期权实时行情。
type CFFEXOptionQuote struct {
	Code          string
	Name          string
	Price         *float64
	Change        *float64
	ChangePercent *float64
	Volume        *float64
	Amount        *float64
	OpenInterest  *float64
	StrikePrice   *float64
	RemainDays    *float64
	DailyChange   *float64
	PrevSettle    *float64
	Open          *float64
}

// OptionLHBItem 是期权龙虎榜数据行。
type OptionLHBItem struct {
	TradeType          string
	Date               string
	Symbol             string
	TargetName         string
	Rank               int
	MemberName         string
	SellVolume         *float64
	SellVolumeChange   *float64
	NetSellVolume      *float64
	SellVolumeRatio    *float64
	BuyVolume          *float64
	BuyVolumeChange    *float64
	NetBuyVolume       *float64
	BuyVolumeRatio     *float64
	SellPosition       *float64
	SellPositionChange *float64
	NetSellPosition    *float64
	SellPositionRatio  *float64
	BuyPosition        *float64
	BuyPositionChange  *float64
	NetBuyPosition     *float64
	BuyPositionRatio   *float64
}
